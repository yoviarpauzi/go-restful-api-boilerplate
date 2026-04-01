package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"aidanwoods.dev/go-paseto"
)

func main() {
	hexKey := "e3fe3a8acf0c5613678076d749290a9a62725d58ee1a2b49974a2bd491cbe7f7"

	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Fatal(err)
	}

	key, err := paseto.V4SymmetricKeyFromBytes(keyBytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(key)
}
