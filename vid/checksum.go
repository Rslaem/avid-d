package vid

import "github.com/QinYuuuu/avid-d/commit/merklecommitment"

type Checksum struct {
	Value merklecommitment.Witness
	Root  *merklecommitment.MerkleTree
}

type StoredChecksum struct {
	Value    merklecommitment.Witness
	Root     *merklecommitment.MerkleTree
	IsStored bool
}

func (checksum *Checksum) Size() int {
	size := 0
	for _, c := range checksum.Value.Hash() {
		size += len(c)
	}
	return size
}

func (checksum *StoredChecksum) Store(c Checksum) {
	checksum.Value = c.Value
	checksum.IsStored = true
	checksum.Root = c.Root
}

func (checksum *StoredChecksum) Load() Checksum {
	return Checksum{
		Value: checksum.Value,
	}
}

func (checksum *StoredChecksum) Stored() bool {
	return checksum.IsStored
}

func (checksum *StoredChecksum) Size() int {
	size := 0
	for _, c := range checksum.Value.Hash() {
		size += len(c)
	}
	return size
}
