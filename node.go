package blackfriday

import (
	"bytes"
	"fmt"
)

type NodeType int

const (
	Document NodeType = iota
	BlockQuote
	List
	Item
	Paragraph
	Header
	HorizontalRule
	Emph
	Strong
	Del
	Link
	Image
	Text
	HTMLBlock
	CodeBlock
	Softbreak
	Hardbreak
	Code
	HTMLSpan
	Table
	TableCell
	TableHead
	TableBody
	TableRow
)

var nodeTypeNames = []string{
	Document:       "Document",
	BlockQuote:     "BlockQuote",
	List:           "List",
	Item:           "Item",
	Paragraph:      "Paragraph",
	Header:         "Header",
	HorizontalRule: "HorizontalRule",
	Emph:           "Emph",
	Strong:         "Strong",
	Del:            "Del",
	Link:           "Link",
	Image:          "Image",
	Text:           "Text",
	HTMLBlock:      "HTMLBlock",
	CodeBlock:      "CodeBlock",
	Softbreak:      "Softbreak",
	Hardbreak:      "Hardbreak",
	Code:           "Code",
	HTMLSpan:       "HTMLSpan",
	Table:          "Table",
	TableCell:      "TableCell",
	TableHead:      "TableHead",
	TableBody:      "TableBody",
	TableRow:       "TableRow",
}

func (t NodeType) String() string {
	return nodeTypeNames[t]
}

type ListData struct {
	ListFlags  ListType
	Tight      bool   // Skip <p>s around list item data if true
	BulletChar byte   // '*', '+' or '-' in bullet lists
	Delimiter  byte   // '.' or ')' after the number in ordered lists
	RefLink    []byte // If not nil, turns this list item into a footnote item and triggers different rendering
}

type LinkData struct {
	Destination []byte
	Title       []byte
	NoteID      int
}

type CodeBlockData struct {
	IsFenced    bool   // Specifies whether it's a fenced code block or an indented one
	Info        []byte // This holds the info string
	FenceChar   byte
	FenceLength int
	FenceOffset int
}

type TableCellData struct {
	IsHeader bool           // This tells if it's under the header row
	Align    CellAlignFlags // This holds the value for align attribute
}

type HeaderData struct {
	Level        int    // This holds the heading level number
	HeaderID     string // This might hold header ID, if present
	IsTitleblock bool   // Specifies whether it's a title block
}

// Node is a single element in the abstract syntax tree of the parsed document.
// It holds connections to the structurally neighboring nodes and, for certain
// types of nodes, additional information that might be needed when rendering.
type Node struct {
	Type       NodeType // Determines the type of the node
	Parent     *Node    // Points to the parent
	FirstChild *Node    // Points to the first child, if any
	LastChild  *Node    // Points to the last child, if any
	Prev       *Node    // Previous sibling; nil if it's the first child
	Next       *Node    // Next sibling; nil if it's the last child

	Literal []byte // Text contents of the leaf nodes

	HeaderData    // Populated if Type == Header
	ListData      // Populated if Type == List
	CodeBlockData // Populated if Type == CodeBlock
	LinkData      // Populated if Type == Link
	TableCellData // Populated if Type == TableCell

	content []byte // Markdown content of the block nodes
	open    bool   // Specifies an open block node that has not been finished to process yet
}

func NewNode(typ NodeType) *Node {
	return &Node{
		Type: typ,
		open: true,
	}
}

func (n *Node) unlink() {
	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else if n.Parent != nil {
		n.Parent.FirstChild = n.Next
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	} else if n.Parent != nil {
		n.Parent.LastChild = n.Prev
	}
	n.Parent = nil
	n.Next = nil
	n.Prev = nil
}

func (n *Node) appendChild(child *Node) {
	child.unlink()
	child.Parent = n
	if n.LastChild != nil {
		n.LastChild.Next = child
		child.Prev = n.LastChild
		n.LastChild = child
	} else {
		n.FirstChild = child
		n.LastChild = child
	}
}

func (n *Node) insertBefore(sibling *Node) {
	sibling.unlink()
	sibling.Prev = n.Prev
	if sibling.Prev != nil {
		sibling.Prev.Next = sibling
	}
	sibling.Next = n
	n.Prev = sibling
	sibling.Parent = n.Parent
	if sibling.Prev == nil {
		sibling.Parent.FirstChild = sibling
	}
}

func (n *Node) isContainer() bool {
	switch n.Type {
	case Document:
		fallthrough
	case BlockQuote:
		fallthrough
	case List:
		fallthrough
	case Item:
		fallthrough
	case Paragraph:
		fallthrough
	case Header:
		fallthrough
	case Emph:
		fallthrough
	case Strong:
		fallthrough
	case Del:
		fallthrough
	case Link:
		fallthrough
	case Image:
		fallthrough
	case Table:
		fallthrough
	case TableHead:
		fallthrough
	case TableBody:
		fallthrough
	case TableRow:
		fallthrough
	case TableCell:
		return true
	default:
		return false
	}
	return false
}

func (n *Node) canContain(t NodeType) bool {
	if n.Type == List {
		return t == Item
	}
	if n.Type == Document || n.Type == BlockQuote || n.Type == Item {
		return t != Item
	}
	if n.Type == Table {
		return t == TableHead || t == TableBody
	}
	if n.Type == TableHead || n.Type == TableBody {
		return t == TableRow
	}
	if n.Type == TableRow {
		return t == TableCell
	}
	return false
}

// WalkStatus allows NodeVisitor to have some control over the tree traversal.
// It is returned from NodeVisitor and different values allow Node.Walk to
// decide which node to go to next.
type WalkStatus int

const (
	GoToNext     WalkStatus = iota // The default traversal of every node.
	SkipChildren                   // Skips all children of current node.
	Terminate                      // Terminates the traversal.
)

// NodeVisitor is a callback to be called when traversing the syntax tree.
// Called twice for every node: once with entering=true when the branch is
// first visited, then with entering=false after all the children are done.
type NodeVisitor func(node *Node, entering bool) WalkStatus

func (root *Node) Walk(visitor NodeVisitor) {
	walker := NewNodeWalker(root)
	node, entering := walker.next()
	for node != nil {
		status := visitor(node, entering)
		switch status {
		case GoToNext:
			node, entering = walker.next()
		case SkipChildren:
			node, entering = walker.resumeAt(node, false)
		case Terminate:
			return
		}
	}
}

type NodeWalker struct {
	current  *Node
	root     *Node
	entering bool
}

func NewNodeWalker(root *Node) *NodeWalker {
	return &NodeWalker{
		current:  root,
		root:     nil,
		entering: true,
	}
}

func (nw *NodeWalker) next() (*Node, bool) {
	if nw.current == nil {
		return nil, false
	}
	if nw.root == nil {
		nw.root = nw.current
		return nw.current, nw.entering
	}
	if nw.entering && nw.current.isContainer() {
		if nw.current.FirstChild != nil {
			nw.current = nw.current.FirstChild
			nw.entering = true
		} else {
			nw.entering = false
		}
	} else if nw.current.Next == nil {
		nw.current = nw.current.Parent
		nw.entering = false
	} else {
		nw.current = nw.current.Next
		nw.entering = true
	}
	if nw.current == nw.root {
		return nil, false
	}
	return nw.current, nw.entering
}

func (nw *NodeWalker) resumeAt(node *Node, entering bool) (*Node, bool) {
	nw.current = node
	nw.entering = entering
	return nw.next()
}

func (ast *Node) String() string {
	return dumpString(ast)
}

func dump(ast *Node) {
	fmt.Println(dumpString(ast))
}

func dump_r(ast *Node, depth int) string {
	if ast == nil {
		return ""
	}
	indent := bytes.Repeat([]byte("\t"), depth)
	content := ast.Literal
	if content == nil {
		content = ast.content
	}
	result := fmt.Sprintf("%s%s(%q)\n", indent, ast.Type, content)
	for n := ast.FirstChild; n != nil; n = n.Next {
		result += dump_r(n, depth+1)
	}
	return result
}

func dumpString(ast *Node) string {
	return dump_r(ast, 0)
}
