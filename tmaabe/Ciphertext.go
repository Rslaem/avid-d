package tmaabe

import (
	//"fmt"
	"github.com/Nik-U/pbc"
	//"math/big"
)

type Ciphertext struct{
	c0	*pbc.Element
	c1	map[string]*pbc.Element
	c2	map[string]*pbc.Element
	c3	map[string]*pbc.Element
	c4	map[string]*pbc.Element
	accessStructure	*AccessStructure
	gp	*GlobalParameters
}

func NewCiphertext(gp *GlobalParameters, ac *AccessStructure) (*Ciphertext){
	c := new(Ciphertext)
	c.gp = gp
	c.accessStructure = ac
	c.c1 = make(map[string]*pbc.Element)
	c.c2 = make(map[string]*pbc.Element)
	c.c3 = make(map[string]*pbc.Element)
	c.c4 = make(map[string]*pbc.Element)
	return c
}
