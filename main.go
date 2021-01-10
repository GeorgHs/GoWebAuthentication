package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("sample-file.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	h := sha256.New()

	_, err = io.Copy(h, f)
	if err != nil {
		log.Fatalln("Couldn't io.copy", err)
	}

	// just for fun get the cumulated Sha256
	fmt.Printf("here's the type BEFORE Sum: %T \n", h)
	fmt.Printf("%v \n", h)
	xb := h.Sum(nil)
	fmt.Printf("here's the type after Sum % T\n", xb)
	fmt.Println("%x\n", xb)

	xb = h.Sum(nil)
	fmt.Printf("here's the type AFTER SECOND Sum: %T\n", xb)
	fmt.Println("%x\n", xb)

	xb = h.Sum(xb)
	fmt.Printf("here's the type AFTER THIRD sum and passing in xb: %T\n", xb)
	fmt.Println("%x\n", xb)
}
