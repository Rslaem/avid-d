package merklecommitment

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/QinYuuuu/avid-d/hasher"
)

type nodeP struct { //(ri,Mi,Pri),i
	RootHash []byte   `json:"rootHash"`
	Content  []byte   `json:"content"` //not hash but the original data
	Proof    *Witness `json:"proof"`
	Index    int      `json:"index"`
}

func GenerateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096) //larger if encryption is needed
	if err != nil {
		fmt.Println("Failed to generate RSA private key:", err)
		return nil
	}
	return privateKey
}
func Zgenerate(p *nodeP, publicKey *rsa.PublicKey) []byte {
	data, _ := json.Marshal(p)
	//fmt.Println("Data:", data)
	z := encrypt(data, publicKey)
	return z
}
func encrypt(data []byte, publicKey *rsa.PublicKey) []byte {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		fmt.Println("Failed to encrypt JSON data:", err)
		return nil
	}
	return ciphertext
}
func decrypt(data []byte, privateKey *rsa.PrivateKey) []byte {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
	if err != nil {
		fmt.Println("Failed to decrypt ciphertext:", err)
		return nil
	}
	//fmt.Println("Plaintext:", plaintext)
	return plaintext
}
func ZVerify(z []byte, privateKey *rsa.PrivateKey) bool {
	plaintext := decrypt(z, privateKey)
	var np nodeP
	json.Unmarshal(plaintext, &np)
	// verify
	result, _ := Verify(np.RootHash, np.Proof, np.Content, hasher.SHA256Hasher)
	return result
}
