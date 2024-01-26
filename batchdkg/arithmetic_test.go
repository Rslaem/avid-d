package batchdkg

import (
	//"bytes"
	//"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestLagrangeInterpolation(t *testing.T) {
	x0 := big.NewInt(3)
	y0 := big.NewInt(35)
	x1 := big.NewInt(10)
	y1 := big.NewInt(1015)
	x2 := big.NewInt(256)
	y2 := big.NewInt(16777477)
	x3 := big.NewInt(50)
	y3 := big.NewInt(125055)
	
	n0, _ := new(big.Int).SetString("1363895147340162124487750544377566700025348452567", 10)
	n1, _ := new(big.Int).SetString("1257354545315887944833595666025792933231792977521", 10)
	n2, _ := new(big.Int).SetString("1296657106138026641358592699056954007605324218609", 10)
	n := new(big.Int)
	n.Mul(n0, n1)
	n.Mul(n, n2)
	
	xs := []*big.Int{x0, x1, x2, x3}
	ys := []*big.Int{y0, y1, y2, y3}

	p, err := LagrangeInterpolation(xs, ys, n)
	assert.Nil(t, err)
	assert.Equal(t, "x³ + x¹ + 5", PolynomialToString(p))

	assert.Equal(t, y0, polynomialEval(p, x0, n))
	assert.Equal(t, y1, polynomialEval(p, x1, n))
	assert.Equal(t, y2, polynomialEval(p, x2, n))
}
