package jango

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
)

var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func rndBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func rndString(n int) string {
	var bb bytes.Buffer
	bb.Grow(n)
	l := uint32(len(runes))
	for i := 0; i < n; i++ {
		bb.WriteRune(runes[binary.BigEndian.Uint32(rndBytes(4))%l])
	}
	return bb.String()
}
