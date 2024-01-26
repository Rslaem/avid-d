package vid

type Checksum struct {
	Value [][]byte
}

type StoredChecksum struct {
	Value    [][]byte
	IsStored bool
}

func (checksum *Checksum) Size() int {
	size := 0
	for _, c := range checksum.Value {
		size += len(c)
	}
	return size
}

func (checksum *StoredChecksum) Store(c Checksum) {
	checksum.Value = c.Value
	checksum.IsStored = true
}

func (checksum *StoredChecksum) Load() (Checksum){
	return Checksum{
		Value: checksum.Value,
	}
}

func (checksum *StoredChecksum) Stored() bool {
	return checksum.IsStored
}

func (checksum *StoredChecksum) Size() int {
	size := 0
	for _, c := range checksum.Value {
		size += len(c)
	}
	return size
}