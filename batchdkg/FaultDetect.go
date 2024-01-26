package batchdkg

import (
	//"github.com/Nik-U/pbc"
	//"log"
	"math/big"
	//"TMAABE/kzg"
	//"strconv"
)

// As the Pj
func (dkg *BatchDKG)FaultDetectPhase1(){
	r_sum := make([]*big.Int, dkg.secretNum)
	for k:=0;k<dkg.secretNum;k++{
		r_sum[k] = new(big.Int).SetInt64(0)
		for i:=0;i<dkg.nodeNum;i++{
			r_sum[k].Add(r_sum[k], dkg.rReceived[i][k])
		}
	}
	dkg.rsum = r_sum
	for i:=0;i<dkg.nodeNum;i++{
		//log.Println(dkg.sharesReceived[i])
		//log.Println(r_sum)
		dkg.aijsum[i] = step2_calculates(dkg.sharesReceived[i], r_sum, dkg.param.n)
	}
}

// As Verifier Pv
func (dkg *BatchDKG)FaultDetectPhase2(i int) (bool){
	ai_sum := make([]*big.Int, dkg.nodeNum)
	x := make([]*big.Int, dkg.nodeNum)
	for j:=0;j<dkg.nodeNum;j++{
		x[j] = new(big.Int).SetInt64(int64(j+1))
		ai_sum[j] = dkg.aijsumReceived[j][i]
	}
	p, _ := LagrangeInterpolation(x, ai_sum, dkg.param.n)
	for j:=dkg.degree+1;j<dkg.nodeNum;j++{
		tmp := polynomialEval(p, new(big.Int).SetInt64(int64(j+1)), dkg.param.n)
		if tmp.Cmp(ai_sum[j]) == 0{
			continue
		}else{
			return false
		}
	}
	return true
}

// P_j calculate {a1j^sum, ...,  anj^sum}
func step2_calculates(f1 []*big.Int, ri []*big.Int, n *big.Int) (*big.Int){
	aij_sum := new(big.Int).SetInt64(0)
	for k := 0; k < len(ri); k++ {
		tmp := new(big.Int).Mul(ri[k], f1[k])
		aij_sum.Add(aij_sum, tmp)
	}
	//fmt.Println(aij_sum)
	aij_sum.Add(aij_sum, f1[len(ri)])
	return aij_sum.Mod(aij_sum, n)
}