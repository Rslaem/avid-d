package tmaabe

import (
	// "github.com/Nik-U/pbc"
	// "math/big"
	//"fmt"
	"strings"
)

type AccessStructure struct {
	rho        map[int]string // row -> att
	tau        map[string]int // att -> committeeID
	A          [][]int
	partsIndex int
	policyTree *TreeNode
}

func (ac *AccessStructure) GetL() int {
	return len(ac.A[0])
}

func (ac *AccessStructure) GetN() int {
	return len(ac.A)
}

func (ac *AccessStructure) GetRow(row int) []int {
	return ac.A[row]
}

type TreeNode struct {
	left  *TreeNode
	right *TreeNode
	name  string
	label string
}

func NewTreeNode() *TreeNode {
	n := new(TreeNode)
	n.left = nil
	n.right = nil
	n.name = ""
	n.label = ""
	return n
}
func (n *TreeNode) SetLeft(left *TreeNode) {
	n.left = left
}
func (n *TreeNode) GetLeft() *TreeNode {
	return n.left
}
func (n *TreeNode) SetRight(right *TreeNode) {
	n.right = right
}
func (n *TreeNode) GetRight() *TreeNode {
	return n.right
}
func (n *TreeNode) SetName(s string) {
	n.name = s
}
func (n *TreeNode) GetName() string {
	return n.name
}
func (n *TreeNode) SetLabel(s string) {
	n.label = s
}
func (n *TreeNode) GetLabel() string {
	return n.label
}
func (aRho *AccessStructure) BuildFromPolicy(policy string) {
	aRho.rho = make(map[int]string)
	aRho.policyTree = NewTreeNode()
	//fmt.Println(policy)
	aRho.generateTree(policy)
	//fmt.Println(aRho.policyTree)
	aRho.generateMatrix()

}

func (aRho *AccessStructure) generateNode(policyParts []string) *TreeNode {
	aRho.partsIndex++
	node := new(TreeNode)
	node.SetName(policyParts[aRho.partsIndex])
	if node.GetName() == "or" || node.GetName() == "and" {
		node.SetLeft(aRho.generateNode(policyParts))
		node.SetRight(aRho.generateNode(policyParts))
	}

	return node
}

func (aRho *AccessStructure) generateTree(policy string) {
	aRho.partsIndex = -1
	var policyParts []string
	if !(strings.HasPrefix(policy, "and") || strings.HasPrefix(policy, "or")) {
		policy = strings.Replace(strings.Replace(policy, " )", ")", -1), "( ", "(", -1)
		policyParts = infixNotationToPolishNotation(strings.Split(policy, " "))
	} else {
		policyParts = strings.Split(policy, " ")
	}
	//fmt.Println(policyParts)
	aRho.policyTree = aRho.generateNode(policyParts)
}

func infixNotationToPolishNotation(policy []string) []string {
	precedence := make(map[string]int)
	precedence["and"] = 2
	precedence["or"] = 1
	precedence["("] = 0

	rpn := NewStack() //rpn stands for Reverse Polish Notation
	operators := NewStack()
	for _, token := range policy {
		//fmt.Println(i)
		if token == "(" {
			operators.Push(token)
		} else if token == ")" {
			for tmp, _ := operators.Top(); tmp != "("; tmp, _ = operators.Top() {
				
				e, _ := operators.Pop()
				rpn.Push(e)
			}
			operators.Pop()
		} else if _, ok := precedence[token]; ok {
			//fmt.Println("operators.Elements")
			for tmp, _ := operators.Top(); (!operators.Empty() && (precedence[token] <= precedence[tmp])); tmp, _ = operators.Top() {
				//fmt.Println("tmp")
				e, _ := operators.Pop()
				rpn.Push(e)
			}
			operators.Push(token)
		} else {
			rpn.Push(token)
		}
	}
	for !operators.Empty() {
		e, _ := operators.Pop()
		rpn.Push(e)
	}
	//fmt.Println(rpn.Elements)
	// reversing the result to obtain Normal Polish Notation
	polishNotation := make([]string, 0)
	for i := range rpn.Elements {
		polishNotation = append(polishNotation, rpn.Elements[len(rpn.Elements)-1-i])
	}
	return polishNotation
}

func (aRho *AccessStructure) toString(builder *strings.Builder, node *TreeNode) {
	if builder.Len() != 0 {
		builder.WriteString(" ")
	}
	if node.GetName() == "and" || node.GetName() == "or" {
		builder.WriteString(node.GetName())
		aRho.toString(builder, node.GetLeft())
		aRho.toString(builder, node.GetRight())
	} else {
		builder.WriteString(node.GetName())
	}
}

func (aRho *AccessStructure) ToString() string {
	builder := new(strings.Builder)
	aRho.toString(builder, aRho.policyTree)
	return builder.String()
}

func (aRho *AccessStructure) generateMatrix() {
	c := computeLabels(aRho.policyTree)
	queue := NewQueue()
	queue.Add(aRho.policyTree)

	for !queue.Empty() {
		node, _ := queue.Poll()

		if node.GetName() == "and" || node.GetName() == "or" {
			queue.Add(node.GetLeft())
			queue.Add(node.GetRight())
		} else {
			aRho.rho[len(aRho.A)] = node.GetName()
			Ax := make([]int, 0)
			for i := 0; i < len(node.GetLabel()); i++ {
				switch node.GetLabel()[i] {
				case '0':
					Ax = append(Ax, 0)
					break 
				case '1':
					Ax = append(Ax, 1)
					break 
				case '*':
					Ax = append(Ax, -1)
					break
				}
			}

			for c > len(Ax) {
				Ax = append(Ax, 0)
			}

			aRho.A = append(aRho.A, Ax)
		}
	}
}

func computeLabels(root *TreeNode) int {
	queue := NewQueue()
	sb := new(strings.Builder)
	c := 1
	root.SetLabel("1")
	queue.Add(root)

	for !queue.Empty() {
		node, _ := queue.Poll()

		if !(node.GetName() == "and" || node.GetName() == "or") {
			continue
		} else if node.GetName() == "or" {
			node.GetLeft().SetLabel(node.GetLabel())
			queue.Add(node.GetLeft())
			node.GetRight().SetLabel(node.GetLabel())
			queue.Add(node.GetRight())
		} else if node.GetName() == "and" {
			sb.Reset()
			sb.WriteString(node.GetLabel())

			for c > sb.Len() {
				sb.WriteString("0")
			}

			sb.WriteString("1")
			node.GetLeft().SetLabel(sb.String())
			queue.Add(node.GetLeft())

			sb.Reset()

			for c > sb.Len() {
				sb.WriteString("0")
			}

			sb.WriteString("*")

			node.GetRight().SetLabel(sb.String())
			queue.Add(node.GetRight())

			c++
		}
	}

	return c
}

func (aRho *AccessStructure) GetMatrixAsString() string {
	sb := new(strings.Builder)
	for x := 0; x < len(aRho.A); x++ {
		Ax := aRho.A[x]
		sb.WriteString(aRho.rho[x] + ": [")
		for _, aAx := range Ax {
			switch aAx {
			case 1:
				sb.WriteString("  1")
				break
			case -1:
				sb.WriteString(" -1")
				break
			case 0:
				sb.WriteString("  0")
				break
			}
		}
		sb.WriteString("]\n")
	}
	result := sb.String()
	return result[:len(result)-1]
}
