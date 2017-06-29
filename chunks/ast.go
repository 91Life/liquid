package chunks

import (
	"io"

	"github.com/osteele/liquid/expressions"
)

// ASTNode is a node of an AST.
type ASTNode interface {
	// Render evaluates an AST node and writes the result to an io.Writer.
	Render(io.Writer, Context) error
}

// ASTRaw holds the text between the start and end of a raw tag.
type ASTRaw struct {
	slices []string
}

// ASTSeq is a sequence of nodes.
type ASTSeq struct {
	Children []ASTNode
}

// ASTFunctional renders itself via a render function that is created during parsing.
type ASTFunctional struct {
	Chunk
	render func(io.Writer, Context) error
}

// ASTText is a text chunk, that is rendered verbatim.
type ASTText struct {
	Chunk
}

// ASTObject is an {{ object }} object.
type ASTObject struct {
	Chunk
	expr expressions.Expression
}

// ASTControlTag is a control tag.
type ASTControlTag struct {
	Chunk
	renderer func(io.Writer, Context) error
	cd       *controlTagDefinition
	Body     []ASTNode
	Branches []*ASTControlTag
}
