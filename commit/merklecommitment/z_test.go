package merklecommitment

import (
	"encoding/json"
	"testing"
	"github.com/QinYuuuu/avid-d/network"
	"github.com/QinYuuuu/avid-d/hasher"
)

// Path: commit/z.go
func TestZGenerate(t *testing.T) {
	// generate
	dataList := [][]byte{[]byte("alice"), []byte("bob"), []byte("cindy"), []byte("david"), []byte("elisa")}
	m, _ := NewMerkleTree(dataList, hasher.SHA256Hasher)
	//comm := Commit(m)
	i := 1
	np := &nodeP{
		RootHash: m.Root().Hash(),
		Content:  dataList[i], //dataList[1]
		Index:    i,
	}
	proof, _ := CreateWitness(m, np.Index)
	//println("hash", proof.HashF, "path", proof.Left)
	np.Proof = proof
	//println("here is the data")
	//fmt.Println(string(data))
	// generate key
	privateKey := GenerateKey()
	publicKey := &privateKey.PublicKey
	z := Zgenerate(np, publicKey)
	//fmt.Println("Z:", z)
	// verify
	t.Run("Verify the z", func(t *testing.T) {
		plaintext := decrypt(z, privateKey)
		//fmt.Println("Plaintext:", plaintext)
		var np nodeP
		json.Unmarshal(plaintext, &np)
		// verify
		//fmt.Println("RootHash", np.RootHash, comm)
		result, _ := Verify(np.RootHash, np.Proof, np.Content, hasher.SHA256Hasher)
		if !result {
			t.Errorf("verify failed")
		}
	})

}
