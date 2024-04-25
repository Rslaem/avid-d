package tmaabe

import (
	"fmt"
	"github.com/Nik-U/pbc"
	//"TMAABE/batchdkg"
	"math/big"
	"sync"
	//"crypto/sha256"
)

func GlobalSetup() (*GlobalParameters){
	gp := &GlobalParameters{}
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
    	params := fmt.Sprintf("type a1\np %s\nn %s\nn0 %s\nn1 %s\nn2 %s\nl %d", pStr, n.String(), n0.String(), n1.String(), n2.String(), l)
    	// elliptic curve
	paramsObj, err := pbc.NewParamsFromString(params)
	if err != nil {
		fmt.Println("[Error]Something wrong with method pbc.NewParamsFromString()")
	}
	gp.pairing = pbc.NewPairing(paramsObj)
    	gp.g = gp.pairing.NewG1().Rand()
    	gp.gn = gp.pairing.NewG1().Rand()
    	return gp
}



func CommitteeSetup(cID int, gp *GlobalParameters, nodeNum int, attributes []string, t int)( []*Node){
	//local setup, nodeNum nodes
	nodes := make([]*Node, nodeNum)
	pairing := gp.pairing
	g := gp.g

	indexlist := []int{}
	var wg sync.WaitGroup
	//secret1list := [][]*pbc.Element	//secret1 
	//secret2list := make(map[string][][]*pbc.Element)	//secret2
	// node init 
	
	for i := 0; i < nodeNum; i++ {
		wg.Add(1)
		go func(i int){
			//fmt.Println("node", i)
			nodes[i] = NewNode(gp, cID, i, attributes)
			indexlist = append(indexlist, i+1)
			/* 
		
			generate secret and distribute
		
			*/
			nodes[i].sk1 = pairing.NewZr().Rand()
			nodes[i].pk1 = pairing.NewG1().PowZn(g, nodes[i].sk1)
			for _, att := range attributes{
				nodes[i].sk2[att] = pairing.NewZr().Rand()
				nodes[i].pk2[att] = pairing.NewG1().PowZn(g, nodes[i].sk2[att])
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	/*
	f1, _ := Step1(gp, nodeNum, len(attributes), t)
	//fmt.Println(f1)
	
	Step2(pairing, nodeNum, len(attributes), t, f1)*/
	fmt.Println(indexlist)
	// aggregate PK
	PK1 := pairing.NewG1().Set1()
	PK2 := make(map[string]*pbc.Element)
	for _, att := range attributes{
		PK2[att] = pairing.NewG1().Set1()
	}
	for i := 0; i < nodeNum; i++ {
		lagrange := pairing.NewZr()
		lagrange.SetBig(GenerateLagrangeCoefficient(indexlist, i+1, gp.n0))
		PK1.Mul(PK1, pairing.NewG1().PowZn(nodes[i].pk1, lagrange))
		for _, att := range attributes{
			PK2[att].Mul(PK2[att], pairing.NewG1().PowZn(nodes[i].pk2[att], lagrange))
		}
		
	}
	for i := 0; i < nodeNum; i++ {
		wg.Add(1)
		go func(i int){
			nodes[i].PK1 = PK1
			nodes[i].PK2 = PK2
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nodes
}

