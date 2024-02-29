package wallet

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"log"
)

type Wallet struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  []byte
}

const (
	version        = byte(0x00)
	checksumLength = 4
)

func (w Wallet) Address() []byte {
	pubKeyHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubKeyHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}

func NewKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panic(err)
	}
	return privateKey, &privateKey.PublicKey
}

func MakeWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	publicKeyBytes := PublicKeyToBytes(publicKey)
	wallet := Wallet{privateKey, publicKeyBytes}
	return &wallet
}

func PublicKeyToBytes(pubKey *rsa.PublicKey) []byte {
    eBytes := make([]byte, 8)
    binary.BigEndian.PutUint64(eBytes, uint64(pubKey.E))

    keyBytes := append(pubKey.N.Bytes(), eBytes...)

    return keyBytes
}

func PublicKeyHash(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)
	return hash[:]
}

func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:checksumLength]
}
