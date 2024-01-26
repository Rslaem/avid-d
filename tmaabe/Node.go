package tmaabe

import (
	"github.com/Nik-U/pbc"
	//"math/big"
	"crypto/sha256"
)

type Node struct {
	committeeID	int				//committee i
	nodeID		int				//node i,j
	gp		*GlobalParameters		
	sk1		*pbc.Element			//m_(i,j)
	pk1		*pbc.Element			//g^(m_(i,j))
	PK1		*pbc.Element			//g^(m_(i))
	sk2		map[string]*pbc.Element	//attribute k -> alpha_(k,i,j)
	pk2		map[string]*pbc.Element	//attribute k -> g^(alpha_(k,i,j))
	PK2		map[string]*pbc.Element	//attribute k -> g^(alpha_(k,i))
	attributes	[]string
	//ipinfo		*IPinfo
}

func NewNode(gp *GlobalParameters, cID int, nid int, attributes []string) (*Node){
	n := new(Node)
	n.committeeID = cID
	n.nodeID = nid
	n.gp = gp
	n.attributes = attributes
	n.sk2 = make(map[string]*pbc.Element)
	n.pk2 = make(map[string]*pbc.Element)
	n.PK2 = make(map[string]*pbc.Element)
	return n
}

func (n *Node) KeyGen(attribute string, gid string) (*AttributeKeyShare){
	pairing := n.gp.pairing
	gn := n.gp.gn
	hgid := pairing.NewG1().SetFromStringHash(gid, sha256.New())
	sk1 := n.sk1			//m_(i,j)
	sk2 := n.sk2[attribute]	//attribute k -> alpha_(k,i,j)
	
	k1 := pairing.NewG1().PowZn(pairing.NewG1().Mul(gn,hgid), sk2)
	k2 := pairing.NewG1().PowZn(hgid, pairing.NewZr().Add(sk1, sk2))
	
	aks := new(AttributeKeyShare)
	aks.key1 = k1
	aks.key2 = k2
	aks.nodeID = n.nodeID
	aks.committeeID = n.committeeID
	aks.attribute = attribute
	//fmt.Println(aks)
	return aks
}
