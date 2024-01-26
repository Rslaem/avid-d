package tmaabe

import(
	"fmt"
)

func TestAC() {
	as1 := new(AccessStructure)
	as1.BuildFromPolicy("A and B and C and D");
	fmt.Println(as1.A)
}