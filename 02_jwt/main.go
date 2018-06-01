package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	bgocrypto "github.com/building-microservices-with-go/crypto"
)

var rsaPrivate *rsa.PrivateKey
var rsaPublic *rsa.PublicKey

func main() {
	command := os.Args[1]
	if command == "generate" {
		fmt.Println(string(GenerateJWT()))
		os.Exit(0)
	}
	token := os.Args[2]
	err := ValidateJWT([]byte(token))
	if err != nil {
		fmt.Printf("error validating token: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("token validated")
}

func init() {
	var err error
	rsaPrivate, err = bgocrypto.UnmarshalRSAPrivateKeyFromFile("./keys/jwt.key")
	if err != nil {
		log.Fatal("Unable to parse private key", err)
	}

	rsaPublic, err = bgocrypto.UnmarshalRSAPublicKeyFromFile("./keys/jwt.key.pub")
	if err != nil {
		log.Fatal("Unable to parse public key", err)
	}
}

// GenerateJWT creates a new JWT and signs it with the private key
func GenerateJWT() []byte {
	claims := jws.Claims{}
	claims.SetExpiration(time.Now().Add(2880 * time.Minute))
	claims.Set("userID", "abcsd232jfjf")
	claims.Set("accessLevel", "user")

	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)

	b, _ := jwt.Serialize(rsaPrivate)

	return b
}

// ValidateJWT validates that the given slice is a valid JWT and the signature matches
// the public key
func ValidateJWT(token []byte) error {
	jwt, err := jws.ParseJWT(token)
	if err != nil {
		return fmt.Errorf("Unable to parse token: %v", err)
	}

	if err = jwt.Validate(rsaPublic, crypto.SigningMethodRS256); err != nil {
		return fmt.Errorf("Unable to validate token: %v", err)
	}

	fmt.Println(jwt.Claims().Get("userID"))

	return nil
}
