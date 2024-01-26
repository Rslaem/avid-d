package tmaabe

import (
	"fmt"
	"github.com/Nik-U/pbc"
	"math/big"
	//"crypto/sha256"
)

type GlobalParameters struct {
	pairing	*pbc.Pairing
	g		*pbc.Element
	gn		*pbc.Element
	p       *big.Int
	n0      *big.Int
	n1      *big.Int
	n2      *big.Int
	l       int
}

func (gp *GlobalParameters) GetPairing() (*pbc.Pairing){
	return gp.pairing
}

func (gp *GlobalParameters) GetGenerateG() (*pbc.Element){
	return gp.g
}

func (gp *GlobalParameters) GetGenerateH() (*pbc.Element){
	return gp.gn
}

func NewGlobalParametersFromString(pairingString string, gByte []byte, hByte []byte) (*GlobalParameters){
	gp := new(GlobalParameters)
	paramsObj, err := pbc.NewParamsFromString(pairingString)
	if err != nil {
		fmt.Println("[Error]Something wrong with method pbc.NewParamsFromString()")
	}
	gp.pairing = pbc.NewPairing(paramsObj)

	pStr := "3602291881362578269408900972923883981249023743695260275790375337088899103553606296737776176077021631283499575659377869566382616215275016597346338059"
	p, _ := new(big.Int).SetString(pStr, 10)
	n0, _ := new(big.Int).SetString("1363895147340162124487750544377566700025348452567", 10)
	n1, _ := new(big.Int).SetString("1257354545315887944833595666025792933231792977521", 10)
	n2, _ := new(big.Int).SetString("1296657106138026641358592699056954007605324218609", 10)
	n := new(big.Int)
	n.Mul(n0, n1)
	n.Mul(n, n2)
	l := 1620

	gp.p = p
    gp.n0 = n0
	gp.n1 = n1
	gp.n2 = n2
	gp.l = l
    gp.g = gp.pairing.NewG1().SetBytes(gByte)
    gp.gn = gp.pairing.NewG1().SetBytes(hByte)
    return gp
}
