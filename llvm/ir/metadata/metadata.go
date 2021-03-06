// === [ Metadata ] ============================================================
//
// References:
//    http://llvm.org/docs/LangRef.html#metadata

// Package metadata provides access to LLVM IR metadata.
package metadata

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/enc"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// A Node represents an LLVM IR metadata node.
//
// Node may have one of the following underlying types.
//
//    *metadata.Metadata   (https://godoc.org/github.com/geode-lang/geode/llvm/ir/metadata#Metadata)
//    *metadata.String     (https://godoc.org/github.com/geode-lang/geode/llvm/ir/metadata#String)
//    *metadata.Value      (https://godoc.org/github.com/geode-lang/geode/llvm/ir/metadata#Value)
//    constant.Constant    (https://godoc.org/github.com/geode-lang/geode/llvm/ir/constant#Constant)
type Node interface {
	value.Value
	// MetadataNode ensures that only metadata nodes can be assigned to the
	// metadata.Node interface.
	MetadataNode()
}

// --- [ metadata ] ------------------------------------------------------------

// Metadata represents a set of LLVM IR metadata nodes.
//
// Metadata may be referenced from instructions (e.g. call), and are thus
// considered LLVM IR values of metadata type.
type Metadata struct {
	// Metadata ID; or empty if metadata literal.
	ID string
	// Metadata nodes.
	Nodes []Node
}

// Type returns the type of the metadata.
func (md *Metadata) Type() types.Type {
	return types.Metadata
}

// Ident returns the identifier associated with the metadata.
func (md *Metadata) Ident() string {
	if len(md.ID) > 0 {
		return enc.Metadata(md.ID)
	}
	return md.Def()
}

// Def returns the LLVM syntax representation of the definition of the metadata.
func (md *Metadata) Def() string {
	buf := &bytes.Buffer{}
	buf.WriteString("!{")
	for i, node := range md.Nodes {
		if i != 0 {
			buf.WriteString(", ")
		}
		if !types.Equal(node.Type(), types.Metadata) {
			fmt.Fprintf(buf, "%s ", node.Type())
		}
		buf.WriteString(node.Ident())
	}
	buf.WriteString("}")
	return buf.String()
}

// MetadataNode ensures that only metadata nodes can be assigned to the
// metadata.Node interface.
func (*Metadata) MetadataNode() {}

// Add a metadata node to the metadata listing
func (md *Metadata) Add(n Node) {
	md.Nodes = append(md.Nodes, n)
}

// --- [ metadata string ] -----------------------------------------------------

// A String represents an LLVM IR metadata string.
type String struct {
	// String value.
	Val string
}

// Ident returns the identifier associated with the metadata.
func (md *String) Ident() string {
	return fmt.Sprintf(`!"%s"`, enc.EscapeString(md.Val))
}

// Type returns the type of the metadata.
func (md *String) Type() types.Type {
	return types.Metadata
}

// MetadataNode ensures that only metadata nodes can be assigned to the
// metadata.Node interface.
func (*String) MetadataNode() {}

// --- [ metadata value ] ------------------------------------------------------

// A Value represents an LLVM IR metadata value.
type Value struct {
	// Value.
	X value.Value
}

// Ident returns the identifier associated with the metadata.
func (md *Value) Ident() string {
	return fmt.Sprintf("%s %s", md.X.Type(), md.X.Ident())
}

// Type returns the type of the metadata.
func (md *Value) Type() types.Type {
	return types.Metadata
}

// MetadataNode ensures that only metadata nodes can be assigned to the
// metadata.Node interface.
func (*Value) MetadataNode() {}

// --- [ named metadata ] ------------------------------------------------------

// Named represents a named collection of metadata, which belongs to a
// module.
type Named struct {
	// Metadata name.
	Name string
	// Associated metadata.
	Metadata []*Metadata
}

// Def returns the LLVM syntax representation of the definition of the named
// metadata.
func (md *Named) Def() string {
	buf := &bytes.Buffer{}
	buf.WriteString("!{")
	for i, metadata := range md.Metadata {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(metadata.Ident())
	}
	buf.WriteString("}")
	return buf.String()
}

// Ident returns the ident of the metadata node
func (md *Named) Ident() string {
	return enc.Metadata(md.Name)
}

// --- [ metadata raw ] -----------------------------------------------------

// NewRaw constructs a raw metadata value
func NewRaw(val string) *Raw {
	r := &Raw{}
	r.Val = val
	return r
}

// A Raw represents raw text to place in the ir
type Raw struct {
	// String value.
	Val string
}

// Ident returns the identifier associated with the metadata.
func (md *Raw) Ident() string {
	return fmt.Sprintf(`!%s`, md.Val)
}

// Type returns the type of the metadata.
func (md *Raw) Type() types.Type {
	return types.Metadata
}

// MetadataNode ensures that only metadata nodes can be assigned to the
// metadata.Node interface.
func (*Raw) MetadataNode() {}
