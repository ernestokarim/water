package main

import (
	"fmt"
	"strconv"
)

type Node interface {
	Type() NodeType
}

// ========================================================

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeList NodeType = iota
	NodeCall
	NodeNumber
	NodeString
	NodeDefine
	NodeVar
)

// ========================================================

type ListNode struct {
	NodeType
	Nodes []Node
}

func newList() *ListNode {
	return &ListNode{
		NodeType: NodeList,
	}
}

// ========================================================

type CallNode struct {
	NodeType
	Name string
	Args []Node
}

func newCall(name string) *CallNode {
	return &CallNode{
		NodeType: NodeCall,
		Name:     name,
		Args:     []Node{},
	}
}

// ========================================================

type NumberNode struct {
	NodeType
	Text string

	IsInt, IsUint bool

	Int64  int64
	Uint64 uint64
}

func newNumber(text string) (*NumberNode, error) {
	n := &NumberNode{
		NodeType: NodeNumber,
		Text:     text,
	}

	u, err := strconv.ParseUint(text, 0, 64)
	if err == nil {
		n.IsUint = true
		n.Uint64 = u
	}

	i, err := strconv.ParseInt(text, 0, 64)
	if err == nil {
		n.IsInt = true
		n.Int64 = i

		if i == 0 {
			n.IsUint = true
			n.Uint64 = u
		}
	}

	if !n.IsUint || !n.IsInt {
		return nil, fmt.Errorf("illegal number syntax: %s", text)
	}

	return n, nil
}

// ========================================================

type StringNode struct {
	NodeType
	Text string
}

func newString(text string) (*StringNode, error) {
	n := &StringNode{
		NodeType: NodeString,
	}

	var err error
	n.Text, err = strconv.Unquote(text)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// ========================================================

type VarNode struct {
	NodeType
	Name string
}

func newVar(name string) *VarNode {
	return &VarNode{
		NodeType: NodeVar,
		Name:     name,
	}
}

// ========================================================

type DefineNode struct {
	NodeType
	variable *VarNode
	value    Node
}

func newDefine(variable *VarNode, value Node) *DefineNode {
	return &DefineNode{
		variable: variable,
		value:    value,
	}
}
