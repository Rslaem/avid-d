package merklecommitment

import (
	//"tmaabe/hasher"
	"errors"
	"bytes"
)

// true:hash is the left node
// false:hash is the right ndoe
type Witness struct{
	hash [][]byte
	left []bool
}

func (w *Witness) SetHash(hash [][]byte) {
	w.hash = hash
}
func (w *Witness) SetPos(left []bool) {
	w.left = left
}
func (w *Witness) Hash() [][]byte{
	return w.hash
}
func (w *Witness) Pos() []bool{
	return w.left
}

// Commit is a function to get the hash of the root of the merkle tree.
//
// Returns:
// - The hash of the root of the tree as node
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	fmt.Println(tree.GetCommitment())
func Commit(m *MerkleTree) []byte {
	//m := NewMerkleTree(data, hasher.SHA256Hasher)
	return m.root.Hash()
}

func CreateWitness(m *MerkleTree, index int) (*Witness, error){
	witnesshash := make([][]byte, 0)
	witnesspos := make([]bool, 0)
	n := m.leafs[index]
	for n.Parent!= nil{
		if n.Parent.Left == n{
			witnesshash = append(witnesshash, n.Parent.Right.hash)
			witnesspos = append(witnesspos, false)
		}else{
			witnesshash = append(witnesshash, n.Parent.Left.hash)
			witnesspos = append(witnesspos, true)
		}
		n = n.Parent
	}
	witness := &Witness{
		hash: witnesshash,
		left: witnesspos,
	}
	return witness, nil
}

func Verify(comm []byte, w *Witness, content []byte, hasher func([]byte) []byte) (bool, error){
	var contenthash []byte
	if len(w.hash) != len(w.left){
		return false, errors.New("error: witness has wrong length")
	}
	for i:=0;i<len(w.hash);i++{
		contenthash = hasher(content)
		if w.left[i]{
			content = append(w.hash[i], contenthash...)
		} else {
			content = append(contenthash, w.hash[i]...)
		}
	}
	return bytes.Equal(comm, hasher(content)), nil
}