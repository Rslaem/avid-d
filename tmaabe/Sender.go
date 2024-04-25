package tmaabe

import (
	//"fmt"
	"github.com/Nik-U/pbc"
	//"math/big"
)

type Sender struct{
	gp		*GlobalParameters
	PK1		map[int]*pbc.Element			//committeeID i -> PK_(1,i)
	PK2		map[int]map[string]*pbc.Element	//committeeID i -> attribute k -> g^(alpha_(k,i))
	attributes	map[int][]string			//committeeID i -> attribute k
}

func NewSender(gp *GlobalParameters, attributes map[int][]string) (*Sender){
	s := new(Sender)
	s.gp = gp
	s.attributes = attributes
	s.PK1 = make(map[int]*pbc.Element)
	s.PK2 = make(map[int]map[string]*pbc.Element)
	return s
}


func (sender *Sender) Encrypt(m *Message, ac *AccessStructure) (*Ciphertext){
	pairing := sender.gp.pairing
	gn := sender.gp.gn
	g := sender.gp.g
	
	lenth := ac.GetL()
	n := ac.GetN()
	
	s := pairing.NewZr().Rand()
	//fmt.Println(pairing.NewGT().PowZn(pairing.NewGT().Pair(gn, g), s))
	v := []*pbc.Element {s}
	w := []*pbc.Element {pairing.NewZr().Neg(s)}
	for i := 1; i < lenth; i++{
		v = append(v, pairing.NewZr().Rand())
		w = append(w, pairing.NewZr().Rand())
	}
	
	c := NewCiphertext(sender.gp, ac)
	c.c0 = pairing.NewGT().Mul(m.mElement, pairing.NewGT().PowZn(pairing.NewGT().Pair(gn, g), s))
	for i := 0; i < n; i++ {
		lambdax := DotProduct(ac.A[i], v, pairing)

		wx := DotProduct(ac.A[i], w, pairing)
		rx := pairing.NewZr().Rand()
		qx := pairing.NewZr().Rand()
		
		att := ac.rho[i]
		cID := ac.tau[att]
		c.c1[att] = pairing.NewG1().PowZn(g, rx)
		c.c2[att] = pairing.NewG1().PowZn(g, qx)
		
		//Ax index i -> attribute k = rho(i) -> committeeIds tau(k) -> pks
		
		PK1 := sender.PK1[cID]
		PK2 := sender.PK2[cID][att]
		c.c3[att] = pairing.NewG1().Mul(pairing.NewG1().PowZn(g, lambdax), pairing.NewG1().PowZn(PK2, rx))
		c.c4[att] = pairing.NewG1().Mul(pairing.NewG1().PowZn(g, wx), pairing.NewG1().PowZn(pairing.NewG1().Mul(PK2, PK1), qx))
	}
	return c
}
