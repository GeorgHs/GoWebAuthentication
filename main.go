package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	msg := "This is totally fun get hands-on and learning it from the ground up."
	password := "Ilovedogs"
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatalln("Couldn't encrypt password", err)
	}
	bs = bs[:16]

	wtr := &bytes.Buffer{}
	encWriter, err := encryptWriter(wtr, bs)

	_, err = io.WriteString(encWriter, msg)
	if err != nil {
		log.Fatalln(err)
	}

	encrypted := wtr.String()
	fmt.Println("before base64", encrypted)

	rslt, err := enDecode(bs, msg)
	if err != nil {
		log.Fatal()
	}
	fmt.Println("before base64", string(rslt))

	rslt2, err := enDecode(bs, string(rslt))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(rslt2))
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
	_, err = io.ReadFull(rand.Reader, iv) //will put randomly put 16 bytes in

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

func encryptWriter(wtr io.Writer, key []byte) (io.Writer, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Couldn't newCipher %w", err)
	}

	//Initialization vector
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv) //will put randomly put 16 bytes in

	s := cipher.NewCTR(b, iv)

	return cipher.StreamWriter{
		S: s,
		W: wtr,
	}, nil
}
