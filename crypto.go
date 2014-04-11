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

type dummyWriterCloser int

func (_ dummyWriterCloser) Write(p []byte) (int, error) {
	return 0, nil
}

func (_ dummyWriterCloser) Close() error {
	return nil
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

func (c Crypto) HmacWriterCloser(w io.Writer, hashKey, blockKey []byte) io.WriteCloser {
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

	return &hmacStreamWriterCloser{w, dummyWriterCloser(0), &bytes.Buffer{}, hmac.New(fn, hashKey)}
}

func (c Crypto) HmacReader(r io.Reader, hashKey, blockKey []byte) (io.Reader, error) {
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
		return nil, ErrorStr(c.c.Lang().Key("errHmacDataIntegrity"))
	}

	return bytes.NewReader(data.B), nil
}

func (c Crypto) Base64HmacWriterCloser(w io.Writer, hashKey, blockKey []byte) io.WriteCloser {
	enc := base64.NewEncoder(base64.URLEncoding, w)
	w = enc
	writer := c.HmacWriterCloser(w, hashKey, blockKey).(*hmacStreamWriterCloser)
	writer.ww = enc
	return writer
}

func (c Crypto) Base64HmacReader(r io.Reader, hashKey, blockKey []byte) (io.Reader, error) {
	return c.HmacReader(base64.NewDecoder(base64.URLEncoding, r), hashKey, blockKey)
}
