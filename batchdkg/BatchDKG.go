package batchdkg

import (
	"math/big"
	"sync"
	//"log"
	"github.com/Nik-U/pbc"
)

type BatchDKG struct {
	// Basic DKG parameters
	// [+] DKG parameters
	param *Parameters
	// [+] number of Secrets
	secretNum int
	// [+] Number of Nodes
	nodeNum int
	// [+] Polynomial Degree
	degree int

	// utilities
	// [+] KZG commitment's PK
	//kzgPK	*kzg.PK

	//Step1 Share Generation
	// [+] Secrets to be shared
	secrets []*big.Int
	// [+] Shares generated in Step1
	shares [][]*big.Int
	// [+] Commitment
	comms [][]*pbc.Element
	// [+] random vector r
	r []*big.Int

	// Step2 Fault Detection
	// [+] the set of qualified nodes' label D
	index []int
	// [+] shares receive from node Pi
	sharesReceived [][]*big.Int
	// [+] commitment receive from node Pi
	commsReceived [][][]*pbc.Element
	// [+] aij_sum[v]
	aijsum         []*big.Int
	aijsumReceived [][]*big.Int
	rReceived      [][]*big.Int
	rsum           []*big.Int
	// Step3 Public Key Construction
	// [+] the Reconstruct PK
	pk         []*pbc.Element
	pkReceived [][]*pbc.Element

	mutex *sync.RWMutex

	// Step4 Public Key Shares Construction
	falseList []int
}

func NewDKG(param *Parameters, secretNum int, nodeNum int, degree int, mutex *sync.RWMutex) *BatchDKG {
	dkg := new(BatchDKG)
	dkg.param = param
	dkg.degree = degree
	dkg.secretNum = secretNum
	dkg.nodeNum = nodeNum
	dkg.mutex = mutex
	dkg.secrets = make([]*big.Int, dkg.secretNum+1)
	dkg.shares = make([][]*big.Int, dkg.nodeNum)
	dkg.sharesReceived = make([][]*big.Int, dkg.nodeNum)
	dkg.rReceived = make([][]*big.Int, dkg.nodeNum)
	dkg.rsum = make([]*big.Int, dkg.secretNum)
	dkg.pkReceived = make([][]*pbc.Element, dkg.nodeNum)

	for i := range dkg.shares {
		dkg.shares[i] = make([]*big.Int, dkg.secretNum+1)
		dkg.sharesReceived[i] = make([]*big.Int, dkg.secretNum+1)
		dkg.pkReceived[i] = make([]*pbc.Element, dkg.secretNum+1)
	}
	dkg.commsReceived = make([][][]*pbc.Element, dkg.nodeNum)
	dkg.comms = make([][]*pbc.Element, dkg.nodeNum)
	dkg.aijsum = make([]*big.Int, dkg.nodeNum)
	dkg.aijsumReceived = make([][]*big.Int, dkg.nodeNum)
	for v := 0; v < dkg.nodeNum; v++ {
		dkg.aijsumReceived[v] = make([]*big.Int, dkg.nodeNum)
		dkg.commsReceived[v] = make([][]*pbc.Element, dkg.nodeNum)
		for j := 0; j < dkg.nodeNum; j++ {
			dkg.commsReceived[v][j] = make([]*pbc.Element, dkg.secretNum+1)
		}
		dkg.rReceived[v] = make([]*big.Int, dkg.secretNum)
		dkg.comms[v] = make([]*pbc.Element, dkg.secretNum+1)
	}
	dkg.r = make([]*big.Int, dkg.secretNum)
	dkg.pk = make([]*pbc.Element, dkg.secretNum)
	dkg.index = make([]int, dkg.nodeNum)
	for i := 0; i < dkg.nodeNum; i++ {
		dkg.index[i] = i + 1
	}
	return dkg
}

func (dkg *BatchDKG) GetParam() *Parameters {
	return dkg.param
}

func (dkg *BatchDKG) SetShares(i int, shares []*big.Int) {
	dkg.sharesReceived[i] = shares
}

func (dkg *BatchDKG) SetComms(i int, commit [][]*pbc.Element) {
	dkg.commsReceived[i] = commit
}

func (dkg *BatchDKG) Receiveaijsum(j int, aijsum []*big.Int) {
	for i := 0; i < dkg.nodeNum; i++ {
		dkg.aijsumReceived[i][j] = aijsum[i]
	}
}

func (dkg *BatchDKG) ReceiveR(v int, r []*big.Int) {
	dkg.rReceived[v] = r
}
func (dkg *BatchDKG) ReceivePK(v int, r []*pbc.Element) {
	dkg.pkReceived[v] = r
}
func (dkg *BatchDKG) Getr() []*big.Int {
	return dkg.r
}
func (dkg *BatchDKG) GetrReceivedFrom(i int) []*big.Int {
	return dkg.rReceived[i]
}
func (dkg *BatchDKG) GetSharesReceivedFrom(i int) []*big.Int {
	return dkg.sharesReceived[i]
}
func (dkg *BatchDKG) GetCommitReceivedFrom(i int) [][]*pbc.Element {
	return dkg.commsReceived[i]
}
func (dkg *BatchDKG) GetaijsumReceived(i int) []*big.Int {
	return dkg.aijsumReceived[i]
}
func (dkg *BatchDKG) GetpkSharesReceived(i int) []*pbc.Element {
	return dkg.pkReceived[i]
}
func (dkg *BatchDKG) Getaijsum() []*big.Int {
	return dkg.aijsum
}

func (dkg *BatchDKG) GetShares() [][]*big.Int {
	return dkg.shares
}

func (dkg *BatchDKG) GetComms() [][]*pbc.Element {
	return dkg.comms
}
