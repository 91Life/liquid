package render

import (
	"io"

	"github.com/osteele/liquid/expression"
)

// ASTNode is a node of an AST.
type ASTNode interface {
	// Render evaluates an AST node and writes the result to an io.Writer.
}

// ASTBlock represents a {% tag %}…{% endtag %}.
type ASTBlock struct {
	Chunk
	renderer func(io.Writer, Context) error
	syntax      BlockSyntax
	Body     []ASTNode
	Branches []*ASTBlock
}

// ASTRaw holds the text between the start and end of a raw tag.
type ASTRaw struct {
	slices []string
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
	expr expression.Expression
}

// ASTSeq is a sequence of nodes.
type ASTSeq struct {
	Children []ASTNode
}
