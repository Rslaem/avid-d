package batchdkg

import (
	"github.com/Nik-U/pbc"
	"log"
	"math/big"
	//"strconv"
)
/*
func (dkg *BatchDKG) PKRecStep1()([]*pbc.Element){
	pairing := dkg.param.GetPairing()
	g := dkg.param.GetGenerateG()
	tmp := make([]*pbc.Element,dkg.secretNum+dkg.nodeNum)
	for i:=0;i< dkg.secretNum+dkg.nodeNum;i++{
		tmp[i] = pairing.NewG1().PowBig(g, dkg.secrets[i])
	}
	return tmp
}
*/

/*
	receive from Pi, execute by verifier Pv Step3
	input:
	output:
*/
func (dkg *BatchDKG) PKRecVerify(i int)(bool){
	pairing := dkg.param.GetPairing()
	g := dkg.param.GetGenerateG()
	gij := dkg.GetCommitReceivedFrom(i)

	//fmt.Println("gi0", gi0)
	//left := pairing.NewG1().Set1()
	//fmt.Println(dkg.aijsumReceived[i])
	
	D_big := make([]*big.Int, len(dkg.index))
	for j:=0;j<len(dkg.index);j++{
		D_big[j] = big.NewInt(int64(dkg.index[j]))
	}
	gi0 := make([]*pbc.Element, dkg.secretNum+1)
	for k:=0;k<dkg.secretNum+1;k++{
		gi0[k] = pairing.NewG1().Set1()
	}
	
	//fmt.Println("D_big", D_big)
	left := pairing.NewG1().Set1()
	for j:=0;j<len(dkg.index);j++{
		lagCoeffecient := generate0LagrangeCoefficient(D_big, j, dkg.param.n)
		power := new(big.Int).Mul(lagCoeffecient, dkg.aijsumReceived[i][j])
		tmp := pairing.NewG1().PowBig(g, power)
		left.Mul(left, tmp)
		
		for k:=0;k<dkg.secretNum+1;k++{
			tmp2 := pairing.NewG1().PowBig(gij[j][k], lagCoeffecient)
			gi0[k].Mul(gi0[k], tmp2)
			//log.Printf("%v",gi0[k])
		} 
	}
	//log.Printf("%v",gi0)
	//fmt.Println("left", left)
	right := gi0[dkg.secretNum]
	for i:=0;i<dkg.secretNum;i++{
		tmp := pairing.NewG1().PowBig(gi0[i], dkg.rsum[i])
		right.Mul(right, tmp)
	}
	dkg.pkReceived[i] = gi0
	return left.Equals(right)
}

// Step4c
func (dkg *BatchDKG) PKRecStep3(){
	a := dkg.pkReceived
	for i:=0; i<dkg.secretNum; i++{
		tmp := dkg.param.pairing.NewG1().Set1()
		for j:=0;j<len(dkg.index);j++{
			tmp.Mul(tmp, a[j][i])
		}
		//poly, _ := LagrangeInterpolation(D_big, tmp, dkg.param.n)
		dkg.pk[i] = tmp
	}
	log.Printf("the public key:\n %v", dkg.pk)
}
