package tmaabe

import (
	"github.com/Nik-U/pbc"
	//"math/big"
)

type Message struct {
	gp		*GlobalParameters
	mElement	*pbc.Element
	mByte		[]byte
}

func MessageFromByte (gp, buf []byte) (m *Message){
	m.mElement = m.gp.pairing.NewGT().SetBytes(buf)
	m.mByte = buf
	return m
}
