package tmaabe

import (
	//"fmt"
	"github.com/Nik-U/pbc"
	//"math/big"
)

type AttributeKeyShare struct {
	key1		*pbc.Element	//nodeID j -> K_{1,k,i,j}
	key2		*pbc.Element	//nodeID j -> K_{1,k,i,j}
	nodeID		int
	committeeID	int			//committeeID i
	attribute	string			//attribute k
}

type AttributeKey struct {
	key1		*pbc.Element		//K_{1,k,i}
	key2		*pbc.Element		//K_{2,k,i}
}
