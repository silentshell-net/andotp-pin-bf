package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/crypto/pbkdf2"
)

func equalSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// this might not be needed, di should always be 16 according to andOTP source code
func matchSlice(original, match []byte) bool {
	for di := range original {

		if len(original[di:]) >= len(match) {

			newSlice := original[di : di+len(match)]

			if equalSlice(newSlice, match) {
				fmt.Println(di)
				return true
			}
		}

	}

	return false
}

// Worker func
func Worker(jobs <-chan string, authDec, saltDec []byte) {
	iter := 2827
	pbkdf2Lenght := 256

	for j := range jobs {

		dk := pbkdf2.Key([]byte(j), saltDec, iter, pbkdf2Lenght, sha1.New)

		if matchSlice(dk, authDec) {
			fmt.Println("Found", j)
			os.Exit(0)
		}

	}
}

func main() {
	var (
		auth, salt       string
		saltDec, authDec []byte
		err              error
	)

	flag.StringVar(&auth, "auth", "", "Base64 value of the 'pref_auth_credentials' key")
	flag.StringVar(&salt, "salt", "", "Base64 value of the 'pref_auth_salt' key")
	flag.Parse()

	if auth == "" || salt == "" {
		flag.Usage()
		os.Exit(0)
	}

	// Decode salt
	if saltDec, err = base64.StdEncoding.DecodeString(salt); err != nil {
		fmt.Println("error decoding salt:", err)
		return
	}
	fmt.Printf("Salt: %x\n", saltDec)

	// Decode auth
	if authDec, err = base64.StdEncoding.DecodeString(auth); err != nil {
		fmt.Println("error decoding auth:", err)
		return
	}
	fmt.Printf("Auth: %x\n", authDec)
	fmt.Println()

	// Simple worker pool
	jobs := make(chan string)
	for w := 1; w <= runtime.NumCPU(); w++ {
		go Worker(jobs, authDec, saltDec)
	}

	for i := 0; i < 9999; i++ {
		jobs <- fmt.Sprintf("%0.4v", i)
	}
}
