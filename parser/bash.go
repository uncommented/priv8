package parser

import (
	"fmt"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
)

// BashParser implements bash script parsing using tree-sitter
type BashParser struct {
	parser   *tree_sitter.Parser
	language *tree_sitter.Language
}

// NewBashParser creates a new bash parser instance
func NewBashParser() (*BashParser, error) {
	// Get the bash language from tree-sitter
	language := tree_sitter.NewLanguage(tree_sitter_bash.Language())

	// Create a new parser
	parser := tree_sitter.NewParser()

	// Set the language for the parser
	parser.SetLanguage(language)

	return &BashParser{
		parser:   parser,
		language: language,
	}, nil
}

// Parse parses bash script content and returns the syntax tree
func (p *BashParser) Parse(content []byte) (*tree_sitter.Tree, error) {
	// Parse the content into a syntax tree
	tree := p.parser.Parse(content, nil)
	if tree == nil {
		return nil, fmt.Errorf("failed to parse bash script")
	}

	return tree, nil
}
