package ast

import (
	"bytes"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

//
// ForNode is a for loop structure representation
type ForNode struct {
	NodeType
	TokenReference

	Index int
	Init  Node
	Cond  Node
	Step  Node
	Body  Node
}

func (n ForNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "for %s; %s; %s %s", n.Init, n.Cond, n.Step, n.Body)
	return buff.String()
}

// NameString implements Node.NameString
func (n ForNode) NameString() string { return "ForNode" }

// Codegen implements Node.Codegen for ForNode
func (n ForNode) Codegen(prog *Program) (value.Value, error) {

	// The name of the blocks is prefixed so we can determine which for loop a block is for.
	namePrefix := fmt.Sprintf("F%X_", n.Index)
	parentBlock := prog.Compiler.CurrentBlock()

	prog.ScopeDown(n.Token)
	var err error
	var predicate value.Value
	var condBlk *ir.Block
	var bodyBlk *ir.Block
	var bodyGenBlk *ir.Block
	var endBlk *ir.Block
	parentFunc := parentBlock.Parent

	condBlk = parentFunc.NewBlock(namePrefix + "cond")

	n.Init.Codegen(prog)

	parentBlock.NewBr(condBlk)

	err = prog.Compiler.genInBlock(condBlk, func() error {
		predicate, _ = n.Cond.Codegen(prog)

		c, err := createTypeCast(prog, predicate, types.I1)
		if err != nil {
			return err
		}
		predicate = c
		return nil
	})

	if err != nil {
		return nil, err
	}

	bodyBlk = parentFunc.NewBlock(namePrefix + "body")

	stepBlk := parentFunc.NewBlock(namePrefix + "step")

	err = prog.Compiler.genInBlock(bodyBlk, func() error {
		scp := prog.Scope
		gen, err := n.Body.Codegen(prog)
		if err != nil {
			return err
		}
		prog.Scope = scp
		bodyGenBlk = gen.(*ir.Block)
		if err != nil {
			return err
		}
		BranchIfNoTerminator(bodyGenBlk, stepBlk)
		BranchIfNoTerminator(bodyBlk, stepBlk)
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = prog.Compiler.genInBlock(stepBlk, func() error {
		scp := prog.Scope
		_, err := n.Step.Codegen(prog)
		prog.Scope = scp
		return err
	})

	if err != nil {
		return nil, err
	}

	BranchIfNoTerminator(stepBlk, condBlk)
	endBlk = parentFunc.NewBlock(namePrefix + "end")
	prog.Compiler.PushBlock(endBlk)
	condBlk.NewCondBr(predicate, bodyBlk, endBlk)

	if err := prog.ScopeUp(); err != nil {
		return nil, err
	}
	return endBlk, nil
}
