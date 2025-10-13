package finder

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

type Filter struct {
	// TODO: when dev freeze, pull dump logic from switch to Dumper implements Renderer for better polymorphism.
	Dump   bool
	Result []string
}

type HasLines interface {
	Lines() *text.Segments
}

func Extract(source []byte, node HasLines) string {
	return string(node.Lines().Value(source))
}

func (f *Filter) traverse(source []byte, node ast.Node, depth int) error {
	indent := strings.Repeat("  ", depth)
	if f.Dump {
		fmt.Printf(indent+"node %v %v\n", node.Type(), node.Kind())
	}
	if hn, ok := node.(*ast.Heading); ok {
		if f.Dump {
			fmt.Println(indent + Extract(source, hn))
		}
		f.Result = append(f.Result, Extract(source, hn))
	}
	if pn, ok := node.(*ast.Paragraph); ok {
		if f.Dump {
			fmt.Println(indent + Extract(source, pn))
		}
		f.Result = append(f.Result, Extract(source, pn))
	}
	if _, ok := node.(*ast.Blockquote); ok {
		return nil
	}
	if node.HasChildren() {
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			if err := f.traverse(source, child, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// Render of [Filter] is a dummy as the data would not output to [io.Writer] but inside the receiver struct members,
// the output input is ignored because words are collected rather than source rendered.
func (f *Filter) Render(_ io.Writer, source []byte, n ast.Node) error {
	if err := f.traverse(source, n, 0); err != nil {
		return err
	}
	if f.Dump {
		n.Dump(source, 0)
	}
	return nil
}

func (f *Filter) AddOptions(option ...renderer.Option) {
	panic(fmt.Errorf("not implemented on options %v", option))
}

// FilterText parses passage as markdown source code, understand the document and output lines that are candidates.
// Headings and Paragraphs are typical candidates, while blockquote and code are not.
func FilterText(passage string) []string {
	md := goldmark.New()
	f := &Filter{Dump: false}
	md.SetRenderer(f)
	err := md.Convert([]byte(passage), os.Stdout)
	if err != nil {
		panic(err)
	}
	return f.Result
}
