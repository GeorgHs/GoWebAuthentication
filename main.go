package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	msg := "This is totally fun get hands-on and learning it from the ground up."
	encoded := encode(msg)
	fmt.Println("ENCODED MSG", encoded)

	s, err := decode(encoded)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("DECODED MSG", s)
}

func encode(msg string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(msg))
	return encoded
}

func decode(encoded string) (string, error) {
	s, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("couldn't decode string %w", err)
	}
	return string(s), nil
}

func enDecode(key []byte, input string) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Couldn't newCipher %w", err)
	}

	//Initialization vector
	iv := make([]byte, aes.BlockSize)

	s := cipher.NewCTR(b, iv)

	buff := &bytes.Buffer{}
	sw := cipher.StreamWriter{
		S: s,
		W: buff,
	}
	_, err = sw.Write([]byte(input))
	if err != nil {
		return nil, fmt.Errorf("couldn't sw.write to streamwriter %w", err)
	}

	return buff.Bytes(), nil
}
