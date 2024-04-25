package vid 

import (
	"bytes"
	//"encoding/binary"
	"encoding/gob"
	//"fmt"

	escode "TMAABE/erasurecode"
)

// StoredErasureCodeChunk is an erasure chunk that is stored. It could be in the memory
// or on the database.
type StoredErasureCodeChunk struct {
	InMemoryValue escode.ErasureCodeChunk
	DBKey         []byte
	IsStored      bool
}

// BUG(leiy): StoredErasureCodeChunk does not have a "delete" or "update" method.

type InMessageChunk struct {
	db      KVStore // will be set in the outgoing path
	key     []byte  // will be set in the outgoing path
	data    []byte  // will be set in the incoming path
	decoded bool
}

func (s *StoredErasureCodeChunk) LoadPointer(db KVStore) *InMessageChunk {
	if s.inMemory() {
		// TODO
		panic("can't load in-memory chunk into pointer")
	}
	return &InMessageChunk{
		db:  db,
		key: s.DBKey,
	}
}

func (s *StoredErasureCodeChunk) StorePointer(c *InMessageChunk, db KVStore) {
	if c.decoded {
		err := db.Put(s.DBKey, c.data)
		if err != nil {
			panic(err)
		}
	} else {
		d, err := db.Get(c.key)
		if err != nil {
			panic(err)
		}
		err = db.Put(s.DBKey, d)
		if err != nil {
			panic(err)
		}
	}
	s.IsStored = true
}

func (s *StoredErasureCodeChunk) Stored() bool {
	return s.IsStored
}

// inMemory returns if the value is currently on the disk
func (s *StoredErasureCodeChunk) inMemory() bool {
	// if the DBKey is set, it is definitely not in the memory
	if s.DBKey != nil {
		return false
	} else {
		return true
	}
}

func (s *StoredErasureCodeChunk) storeOnDisk(val escode.ErasureCodeChunk, db KVStore) {
	//bm, is := val.(encoding.BinaryMarshaler)
	var b []byte
	/*
		// we can't use marshalbinary here because we don't have the concrete type
		// to decode into when retrieving from the disk
		if is {
			b, err := bm.MarshalBinary()
			if err != nil {
				panic(err)
			}
			s.bmarshaller = true
		} else {
	*/
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(&val)
	if err != nil {
		panic(err)
	}
	b = buf.Bytes()
	//}
	err = db.Put(s.DBKey, b)
	if err != nil {
		panic(err)
	}
}

// Store stores the given ErasureCodeChunk to the StoredErasureCodeChunk.
func (s *StoredErasureCodeChunk) Store(val escode.ErasureCodeChunk, db KVStore) {
	s.IsStored = true
	if s.inMemory() {
		s.InMemoryValue = val
	}
}

func (s *StoredErasureCodeChunk) Load() escode.ErasureCodeChunk {
	return s.InMemoryValue
}

func (s *StoredErasureCodeChunk) Stash(db KVStore, key []byte) {
	s.DBKey = key
	if s.IsStored {
		s.storeOnDisk(s.InMemoryValue, db)
		s.InMemoryValue = nil
	}
}