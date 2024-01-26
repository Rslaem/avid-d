package batchdkg

import (
	"github.com/Nik-U/pbc"
	"math/big"
)

type Parameters struct{
	pairing	*pbc.Pairing
	g		*pbc.Element
	n		*big.Int
}

// New a batchDKG's parameter
func NewParam(pairing *pbc.Pairing, g *pbc.Element, n *big.Int)(*Parameters){
	return &Parameters{
		pairing:	pairing,
		g:			g,
		n:			n,
	}
}

func (param *Parameters) SetPairing(pairing *pbc.Pairing){
	param.pairing = pairing
}

func (param *Parameters) GetPairing()(pairing *pbc.Pairing){
	return param.pairing
}

func (param *Parameters) GetGenerateG()(*pbc.Element){
	return param.g
}

func (param *Parameters) SetN(n *big.Int){
	param.n = n
}

func (param *Parameters) GetN()(n *big.Int){
	return param.n
}

func (param *Parameters) SetGenerateG(g *pbc.Element){
	param.g = g
}
