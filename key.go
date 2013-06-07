package core

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
)

// Convert Unsigned 64-bit Int to Bytes.
func uint64ToByte(num uint64) [8]byte {
	var buf [8]byte
	buf[0] = byte(num >> 0)
	buf[1] = byte(num >> 8)
	buf[2] = byte(num >> 16)
	buf[3] = byte(num >> 24)
	buf[4] = byte(num >> 32)
	buf[5] = byte(num >> 40)
	buf[6] = byte(num >> 48)
	buf[7] = byte(num >> 56)
	return buf
}

// AES-256 Friendly, Great for Session ID's
func KeyGen() string {
	const keyLen = 32

	curtime := time.Now()
	second := uint64ToByte(uint64(curtime.Unix()))
	nano := uint64ToByte(uint64(curtime.UnixNano()))

	rand1 := uint64ToByte(uint64(rand.Int63()))
	rand2 := uint64ToByte(uint64(rand.Int63()))

	b := []byte{}

	for key, value := range second {
		b = append(b, value, rand1[key])
	}

	for key, value := range nano {
		b = append(b, value, rand2[key])
	}

	hash := sha256.New()
	defer hash.Reset()

	hash.Write(b)
	b = hash.Sum(nil)

	str := base64.URLEncoding.EncodeToString(b)

	return str[:keyLen]
}
