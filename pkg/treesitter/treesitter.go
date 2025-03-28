package treesitter

import (
	"fmt"
	"strings"
	"unsafe"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// TreeSitterManager handles tree-sitter initialization and language loading
type TreeSitterManager struct {
	// Map of language name to language instance
	languages map[string]*tree_sitter.Language

	// Map of language to custom queries
	queries map[string]map[string]*tree_sitter.Query
}

// NewTreeSitterManager creates a new tree-sitter manager
func NewTreeSitterManager() *TreeSitterManager {
	return &TreeSitterManager{
		languages: make(map[string]*tree_sitter.Language),
		queries:   make(map[string]map[string]*tree_sitter.Query),
	}
}

// LoadLanguage loads a tree-sitter language
func (m *TreeSitterManager) LoadLanguage(name string, language *tree_sitter.Language) error {
	if language == nil {
		return fmt.Errorf("language %s is nil", name)
	}

	m.languages[name] = language
	m.queries[name] = make(map[string]*tree_sitter.Query)
	return nil
}

// GetLanguage returns a loaded tree-sitter language
func (m *TreeSitterManager) GetLanguage(name string) (*tree_sitter.Language, error) {
	language, ok := m.languages[name]
	if !ok {
		return nil, fmt.Errorf("language %s not loaded", name)
	}

	return language, nil
}

// RegisterQuery registers a named query for a language
func (m *TreeSitterManager) RegisterQuery(language string, name string, queryText string) error {
	lang, err := m.GetLanguage(language)
	if err != nil {
		return err
	}

	query, err := tree_sitter.NewQuery([]byte(queryText), lang)
	if err != nil {
		return fmt.Errorf("failed to create query %s for language %s: %w", name, language, err)
	}

	m.queries[language][name] = query
	return nil
}

// GetQuery returns a registered query
func (m *TreeSitterManager) GetQuery(language string, name string) (*tree_sitter.Query, error) {
	langQueries, ok := m.queries[language]
	if !ok {
		return nil, fmt.Errorf("no queries registered for language %s", language)
	}

	query, ok := langQueries[name]
	if !ok {
		return nil, fmt.Errorf("query %s not registered for language %s", name, language)
	}

	return query, nil
}

// CreateParser creates a new tree-sitter parser for the specified language
func (m *TreeSitterManager) CreateParser(language string) (*tree_sitter.Parser, error) {
	lang, err := m.GetLanguage(language)
	if err != nil {
		return nil, err
	}

	parser := tree_sitter.NewParser()
	parser.SetLanguage(lang)

	return parser, nil
}

// ExecuteQuery executes a named query on a tree-sitter node
func (m *TreeSitterManager) ExecuteQuery(language string, queryName string, node *tree_sitter.Node) ([]*tree_sitter.QueryMatch, error) {
	query, err := m.GetQuery(language, queryName)
	if err != nil {
		return nil, err
	}

	cursor := tree_sitter.NewQueryCursor()
	cursor.Exec(query, node)

	var matches []*tree_sitter.QueryMatch
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		matches = append(matches, match)
	}

	return matches, nil
}

// NodeToString converts a tree-sitter node to a string for debugging
func (m *TreeSitterManager) NodeToString(node *tree_sitter.Node, source []byte, indent string) string {
	if node == nil {
		return indent + "<nil>"
	}

	result := fmt.Sprintf("%sType: %s\n", indent, node.Type())
	result += fmt.Sprintf("%sRange: (%d,%d) - (%d,%d)\n",
		indent, node.StartPoint().Row, node.StartPoint().Column,
		node.EndPoint().Row, node.EndPoint().Column)

	if source != nil {
		nodeText := source[node.StartByte():node.EndByte()]
		if len(nodeText) > 40 {
			nodeText = append(nodeText[:37], []byte("...")...)
		}
		result += fmt.Sprintf("%sText: %q\n", indent, nodeText)
	}

	result += fmt.Sprintf("%sChildren: %d\n", indent, node.ChildCount())

	if node.ChildCount() > 0 {
		result += fmt.Sprintf("%sNamed Children: %d\n", indent, node.NamedChildCount())

		// Print named children
		for i := uint32(0); i < node.NamedChildCount(); i++ {
			child := node.NamedChild(i)
			if child != nil {
				result += fmt.Sprintf("%sNamed Child %d:\n", indent, i)
				result += m.NodeToString(child, source, indent+"  ")
			}
		}
	}

	return result
}

// Custom tree-sitter node for privacy scanning

// PrivacyNode extends tree-sitter Node with privacy-specific metadata
type PrivacyNode struct {
	Node             *tree_sitter.Node
	IssueType        string
	Sensitivity      int // 0-10 scale
	ReplacementValue string
	Context          map[string]interface{}
}

// NewPrivacyNode creates a new privacy node
func NewPrivacyNode(node *tree_sitter.Node) *PrivacyNode {
	return &PrivacyNode{
		Node:        node,
		Sensitivity: 5, // Default medium sensitivity
		Context:     make(map[string]interface{}),
	}
}

// ID returns a unique identifier for the node
func (n *PrivacyNode) ID() uint32 {
	return uint32(uintptr(unsafe.Pointer(n.Node)))
}

// GetText extracts the text content of the node
func (n *PrivacyNode) GetText(source []byte) string {
	if n.Node == nil {
		return ""
	}

	start := n.Node.StartByte()
	end := n.Node.EndByte()

	if start < 0 || end > uint32(len(source)) || start > end {
		return ""
	}

	return string(source[start:end])
}

// GetPositionText returns a string representation of the node's position
func (n *PrivacyNode) GetPositionText() string {
	if n.Node == nil {
		return "unknown position"
	}

	return fmt.Sprintf("line %d:%d to %d:%d",
		n.Node.StartPoint().Row+1, n.Node.StartPoint().Column+1,
		n.Node.EndPoint().Row+1, n.Node.EndPoint().Column+1)
}

// TreeSitterNodeTraversal provides utility functions to traverse tree-sitter trees
type TreeSitterNodeTraversal struct {
	// Optional: Store traversal state here
}

// FindNodesOfType finds all nodes of a specified type in the tree
func (t *TreeSitterNodeTraversal) FindNodesOfType(rootNode *tree_sitter.Node, nodeType string) []*tree_sitter.Node {
	var results []*tree_sitter.Node
	t.traverseAndFind(rootNode, func(node *tree_sitter.Node) bool {
		return node.Type() == nodeType
	}, &results)
	return results
}

// FindNodesMatching finds all nodes matching a predicate
func (t *TreeSitterNodeTraversal) FindNodesMatching(rootNode *tree_sitter.Node,
	predicate func(*tree_sitter.Node) bool,
) []*tree_sitter.Node {
	var results []*tree_sitter.Node
	t.traverseAndFind(rootNode, predicate, &results)
	return results
}

// traverseAndFind recursively traverses the tree and applies the predicate
func (t *TreeSitterNodeTraversal) traverseAndFind(node *tree_sitter.Node,
	predicate func(*tree_sitter.Node) bool, results *[]*tree_sitter.Node,
) {
	if node == nil {
		return
	}

	// Check if this node matches
	if predicate(node) {
		*results = append(*results, node)
	}

	// Recursively check children
	for i := uint32(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child != nil {
			t.traverseAndFind(child, predicate, results)
		}
	}
}

// FindNodesByText finds nodes containing specific text
func (t *TreeSitterNodeTraversal) FindNodesByText(rootNode *tree_sitter.Node, source []byte,
	textPredicate func(string) bool,
) []*tree_sitter.Node {
	return t.FindNodesMatching(rootNode, func(node *tree_sitter.Node) bool {
		if node == nil {
			return false
		}

		// Get node text
		start := node.StartByte()
		end := node.EndByte()

		if start < 0 || end > uint32(len(source)) || start > end {
			return false
		}

		nodeText := string(source[start:end])
		return textPredicate(nodeText)
	})
}

// FindNodesContainingText finds nodes containing a specific text
func (t *TreeSitterNodeTraversal) FindNodesContainingText(rootNode *tree_sitter.Node, source []byte,
	text string,
) []*tree_sitter.Node {
	return t.FindNodesByText(rootNode, source, func(nodeText string) bool {
		return strings.Contains(nodeText, text)
	})
}
