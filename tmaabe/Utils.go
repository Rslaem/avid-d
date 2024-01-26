package tmaabe

import (
	//"fmt"
	"github.com/Nik-U/pbc"
	"math/big"
	"crypto/sha256"
)

//lagrange coefficient
func GenerateLagrangeCoefficient (list []int, index int, q *big.Int) (*big.Int){
	acc, _ := new(big.Int).SetString("1", 10)
	iBIg := new(big.Int).SetUint64(uint64(index))
	jBig := new(big.Int)
	diff := new(big.Int)
	diffModInv := new(big.Int)
	for _, j := range list {
		if j == index {
			continue
		} else{ 
			jBig.SetUint64(uint64(j))
			diff.Sub(jBig, iBIg)
			diffModInv.ModInverse(diff, q)
			acc.Mul(acc, jBig)
			acc.Mul(acc, diffModInv)
			acc.Mod(acc, q)
		}
	}
	return acc
}

func DotProduct(v1 []int, v2 []*pbc.Element, pairing *pbc.Pairing) (*pbc.Element){
	e := pairing.NewZr()
	result := pairing.NewZr().Set0()
	for i := 0;i<len(v1);i++{
		if v1[i] == 1{
			e.Set1()
		} else if v1[i] == 0 {
			e.Set0()
		} else if v1[i] == -1 {
			e.Neg(pairing.NewZr().Set1())
		}
		result = pairing.NewZr().Add(result, pairing.NewZr().Mul(e, v2[i]))
	}
	return result
}

// Gaussian Elimination Solve Equations
func GaussianElimination (A [][]*pbc.Element, b []*pbc.Element, pairing *pbc.Pairing)([]*pbc.Element){
	// A{m row, n column}; b = {1,0,...,0} m
	// return x {x1,...,xm } m (num of attribute)
	n := len(b)
	x := make([]*pbc.Element, n)
	for k:=0; k<n; k++{
		for i:=k+1;i<n;i++{
			scale := pairing.NewZr().Div(A[i][k], A[k][k])
            		for j:=k;j<n;j++ {
                		A[i][j] = pairing.NewZr().Sub(A[i][j], pairing.NewZr().Mul(scale, A[k][j]))
                	}
            		b[i] = pairing.NewZr().Sub(b[i], pairing.NewZr().Mul(scale, b[k]))
		}
	}
    	for k:=n-1; k>=0; k--{
    		tmp := pairing.NewZr().Set0()
    		for i:=k+1;i<n;i++{
    			tmp = pairing.NewZr().Add(tmp, pairing.NewZr().Mul(A[k][i],x[i]))
    		}
    		x[k] = pairing.NewZr().Div(pairing.NewZr().Sub(b[k], tmp), A[k][k])
    	}
	return x
}

// if s in slist
func Contains(slist []string, s string) (bool){
	smap := make(map[string]struct{}, len(slist))
	for _, v := range slist {
		smap[v] = struct{}{}
	}
	_, ok := smap[s]
	return ok
}

func hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
