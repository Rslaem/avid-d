package tmaabe

import (
	//"fmt"
	"github.com/Nik-U/pbc"
	//"math/big"
	"crypto/sha256"
)

type Receiver struct{
	gp		*GlobalParameters
	gid		string
	attributeKey	map[string]*AttributeKey
	attributes	[]string
	qualified	map[int][]int	//qualifeid nodes B_{i,j}
	tau		map[string]int	// att -> committeeID
}

func NewReceiver(gp *GlobalParameters, gid string, attributes []string) (*Receiver){
	r := new(Receiver)
	r.gp = gp
	r.gid = gid
	r.attributes = attributes
	r.attributeKey = make(map[string]*AttributeKey)
	r.qualified = make(map[int][]int)
	return r
}


func (r *Receiver) KeyVer(aks *AttributeKeyShare, pk1 *pbc.Element, pk2 *pbc.Element) (bool){
	pairing := r.gp.pairing
	g := r.gp.g
	gn := r.gp.gn
	k1 := aks.key1
	k2 := aks.key2

	hgid := pairing.NewG1().SetFromStringHash(r.gid, sha256.New())
	
	left := pairing.NewGT()
	right := pairing.NewGT()
	right1 := pairing.NewGT()
	right2 := pairing.NewGT()
	// verify eq1
	left.Pair(g, k1)
	right1.Pair(gn, pk2)
	right2.Pair(hgid, pk2)
	right.Mul(right1, right2)
	
	if !left.Equals(right){
		return false
	}
	// verify eq2
	left.Pair(g, k2)
	right1.Pair(hgid, pk1)
	right.Mul(right1, right2)
	if !left.Equals(right){
		return false
	}
	return true
}

func (r *Receiver) KeyAggregate (keymap map[int]*AttributeKeyShare, nodelist []int, cID int, attribute string) {
	pairing := r.gp.pairing
	ak := new(AttributeKey)
	ak.key1 = pairing.NewG1().Set1()
	ak.key2 = pairing.NewG1().Set1()
	
	indexlist := make([]int, len(nodelist))
	
	for i, _:= range nodelist{
		indexlist[i] = nodelist[i]+1
	}
	//fmt.Println(indexlist)
	
	for _, j := range indexlist {
		lagrange := pairing.NewZr()
		lagrange.SetBig(GenerateLagrangeCoefficient(indexlist, j, r.gp.n0))
		//fmt.Println(lagrange)
		
		ak.key1 = pairing.NewG1().Mul(ak.key1, pairing.NewG1().PowZn(keymap[j-1].key1, lagrange))
		ak.key2 = pairing.NewG1().Mul(ak.key2, pairing.NewG1().PowZn(keymap[j-1].key2, lagrange))
	}
	r.attributeKey[attribute] = ak
}

func (r *Receiver) Decrypt(c *Ciphertext) (*Message){ 
	pairing := r.gp.pairing
	gn := r.gp.gn
	hgid := pairing.NewG1().SetFromStringHash(r.gid, sha256.New())
	tmp := pairing.NewGT().Set1()
	
	ac := c.accessStructure
	
	toUse := []int{}
	for i := 0; i< len(ac.rho); i++{
		//fmt.Println(i, att)
		att := ac.rho[i]
		if Contains(r.attributes, att){
			toUse = append(toUse, i)
		}
	}
	
	//fmt.Println("toUse", toUse)
	// submatrix of A corresponding to user attribute
	subA := [][]int{}
	for _, index := range toUse {
		subA = append(subA, ac.A[index])
	}
	//fmt.Println("subA\n",subA)
	// matrix transpose
	subA_T := make([][]*pbc.Element, len(subA[0]))
	for i,_:=range subA_T {
		subA_T[i]=make([]*pbc.Element, len(subA))
	}
    	for i:=0;i<len(subA);i++ {
    		for j:=0;j<len(subA[0]);j++ {
    			subA_T[j][i] = pairing.NewZr()
			if subA[i][j] == 1{
				subA_T[j][i].Set1()
			} else if subA[i][j] == 0 {
				subA_T[j][i].Set0()
			} else if subA[i][j] == -1 {
				subA_T[j][i].Neg(pairing.NewZr().Set1())
			}
		}
	}

	b := []*pbc.Element {pairing.NewZr().Set1()}
	for i := 1; i < len(toUse); i++{
		b = append(b, pairing.NewZr().Set0())
	}
	cx := GaussianElimination(subA_T, b, pairing)
	
	for i, index := range toUse{
		att := ac.rho[index]
		tmp1 := pairing.NewGT().Div(pairing.NewGT().Pair(c.c3[att], pairing.NewG1().Mul(gn, hgid)), pairing.NewGT().Pair(c.c1[att], r.attributeKey[att].key1))
		tmp2 := pairing.NewGT().Div(pairing.NewGT().Pair(c.c4[att], hgid), pairing.NewGT().Pair(c.c2[att], r.attributeKey[att].key2))
		tmp.Mul(tmp, pairing.NewGT().PowZn(pairing.NewGT().Mul(tmp1, tmp2), cx[i]))
	}
	m := new(Message)
	m.mElement = pairing.NewGT().Div(c.c0, tmp)
	return m
}
