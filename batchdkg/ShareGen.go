package batchdkg

import (
	//"github.com/Nik-U/pbc"
	//"fmt"
	//"crypto/rsa"
	"math/big"
	//"sync"
	"strconv"
)

func (dkg *BatchDKG) ShareGen()([]bool){
	pairing := dkg.param.pairing
	g := dkg.param.g
	universal_fij := step1_function_generate(dkg.secretNum + 1, dkg.degree, dkg.param.n)
	//x := make([]*big.Int, dkg.secretNum + dkg.nodeNum)
	for i:=0;i<dkg.secretNum + 1;i++{
		dkg.secrets[i] = universal_fij[i][0]
	}
	
	for i:=0;i<dkg.nodeNum;i++{
		dkg.shares[i] = step1_calculate(universal_fij, i+1, dkg.param.n)
		for j:=0;j<dkg.secretNum+1;j++{
			dkg.comms[i][j] = pairing.NewG1().PowBig(g, dkg.shares[i][j]) 
		} 
		// p, _ := LagrangeInterpolation(x, dkg.shares[i], dkg.param.n)
		// fmt.Println(p)
		// dkg.kzgC[i] = kzg.Commit(dkg.kzgPK, p)
	}

	// generate the polynomial
	//x := make([]*big.Int, dkg.degree+1)
	//ai0 := make([]*big.Int, dkg.secretNum)
	/*
	for i:=0;i<dkg.degree+1;i++{
		x[i] = new(big.Int).SetInt64(int64(i+1))
	}
	y := make([][]*big.Int, dkg.secretNum + 1)
	for i:=0;i<dkg.secretNum + 1;i++{
		y[i] = make([]*big.Int, dkg.nodeNum)
		for j:=0;j<dkg.nodeNum;j++{
			y[i][j] = dkg.shares[j][i]
		}
	}
	*/
	// 生成向量组ri
	for i := 0; i < dkg.secretNum; i++ {
		//dkg.r[i] = new(big.Int).SetInt64(0)
		dkg.r[i], _ = RandBigInt(dkg.param.n) 
	}

	// verify the shares
	flag := make([]bool, dkg.secretNum + 1)
	for i:=0;i<dkg.secretNum + 1;i++{
		flag[i] = true
		
	}
	return flag
}

// P_i generate k+n polynomials t degree
func step1_function_generate(n int, t int, R *big.Int)([][]*big.Int) {
	universal_fij := make([][]*big.Int, n)
	for j := 0; j < n; j++ {
		universal_fij[j] = make([]*big.Int, t+1)
		for k := 0; k < t+1; k++ {
			universal_fij[j][k], _ = RandBigInt(R)
		}
	}
	return universal_fij
}

// P_i calculate k+n shares to P_j
func step1_calculate(universal_fij [][]*big.Int, j int, R *big.Int) ([]*big.Int){
	l := len(universal_fij)
	f1 := make([]*big.Int, l)
	bigj, _ := new(big.Int).SetString(strconv.Itoa(j), 10)
	for k := 0; k < l; k++ {
		f1[k] = polynomialEval(universal_fij[k], bigj, R)
	}
	return f1
}
