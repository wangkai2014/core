package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"hash"
	"io"
)

type hmacData struct {
	B []byte
	M []byte
}

func init() {
	gob.Register(hmacData{})
}

type hmacStreamWriterCloser struct {
	w    io.Writer
	ww   io.WriteCloser
	buf  *bytes.Buffer
	hash hash.Hash
}

func (h *hmacStreamWriterCloser) Write(p []byte) (int, error) {
	return h.buf.Write(p)
}

func (h *hmacStreamWriterCloser) Close() error {
	defer h.buf.Reset()
	h.hash.Write(h.buf.Bytes())

	data := hmacData{
		B: h.buf.Bytes(),
		M: h.hash.Sum(nil),
	}

	err := gob.NewEncoder(h.w).Encode(data)
	if err != nil {
		return err
	}

	return h.ww.Close()
}

type Crypto struct {
	c *Context
}

func (c *Context) Crypto() Crypto {
	return Crypto{c}
}

func (c Crypto) AesOfbWriter(w io.Writer, blockKey []byte) (io.Writer, error) {
	block, err := aes.NewCipher(blockKey)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte

	stream := cipher.NewOFB(block, iv[:])

	return &cipher.StreamWriter{S: stream, W: w}, nil
}

func (c Crypto) AesOfbReader(r io.Reader, blockKey []byte) (io.Reader, error) {
	block, err := aes.NewCipher(blockKey)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte

	stream := cipher.NewOFB(block, iv[:])

	return &cipher.StreamReader{S: stream, R: r}, nil
}

func (c Crypto) HmacWriterCloser(w io.Writer, hashKey, blockKey []byte) (io.WriteCloser, error) {
	enc := base64.NewEncoder(base64.URLEncoding, w)
	w = enc
	ww, err := c.AesOfbWriter(w, blockKey)
	if err == nil {
		w = ww
	}

	fn := c.c.App.HashFunc
	if fn == nil {
		fn = func() hash.Hash {
			return sha256.New()
		}
	}

	return &hmacStreamWriterCloser{w, enc, &bytes.Buffer{}, hmac.New(fn, hashKey)}, nil
}

func (c Crypto) HmacReader(r io.Reader, hashKey, blockKey []byte) (io.Reader, error) {
	r = base64.NewDecoder(base64.URLEncoding, r)
	rr, err := c.AesOfbReader(r, blockKey)
	if err == nil {
		r = rr
	}

	fn := c.c.App.HashFunc
	if fn == nil {
		fn = func() hash.Hash {
			return sha256.New()
		}
	}

	data := hmacData{}

	err = gob.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}

	hash := hmac.New(fn, hashKey)
	hash.Write(data.B)

	if !hmac.Equal(data.M, hash.Sum(nil)) {
		return nil, ErrorStr("Data has been tempered with!")
	}

	return bytes.NewReader(data.B), nil
}
