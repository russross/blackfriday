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
	HtmlBlock
	CodeBlock
	Softbreak
	Hardbreak
	Code
	HtmlSpan
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
	HtmlBlock:      "HtmlBlock",
	CodeBlock:      "CodeBlock",
	Softbreak:      "Softbreak",
	Hardbreak:      "Hardbreak",
	Code:           "Code",
	HtmlSpan:       "HtmlSpan",
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
	FenceLength uint32
	FenceOffset uint32
}

type Node struct {
	Type       NodeType
	Parent     *Node
	FirstChild *Node
	LastChild  *Node
	Prev       *Node // prev sibling
	Next       *Node // next sibling

	content []byte
	open    bool

	Level   uint32 // If Type == Header, this holds the heading level number
	Literal []byte

	ListData             // If Type == List, this holds list info
	CodeBlockData        // If Type == CodeBlock, this holds its properties
	LinkData             // If Type == Link, this holds link info
	HeaderID      string // If Type == Header, this might hold header ID, if present
	IsTitleblock  bool
	IsHeader      bool           // If Type == TableCell, this tells if it's under the header row
	Align         CellAlignFlags // If Type == TableCell, this holds the value for align attribute
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

func (nw *NodeWalker) resumeAt(node *Node, entering bool) {
	nw.current = node
	nw.entering = entering
}

func ForEachNode(root *Node, f func(node *Node, entering bool)) {
	walker := NewNodeWalker(root)
	node, entering := walker.next()
	for node != nil {
		f(node, entering)
		node, entering = walker.next()
	}
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
