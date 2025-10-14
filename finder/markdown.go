package finder

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

type Text struct {
	Item string
	Path string
}

type Filter struct {
	// TODO: when dev freeze, pull dump logic from switch to Dumper implements Renderer for better polymorphism.
	Dump   bool
	Result []Text

	hh HeadingHelper
}

type HasLines interface {
	Lines() *text.Segments
}

func Extract(source []byte, node HasLines) string {
	return string(node.Lines().Value(source))
}

// HeadingHelper helps to conclude the path of current content.
// It requires a non-skip headings hierarchy to work.
type HeadingHelper struct {
	memory [6]string // markdown supports H1 to H6 total 6 levels of headings
	recent int       // the most recent updated memory key
}

func (hh *HeadingHelper) Path() string {
	levels := hh.memory[:hh.recent+1] // remove outdated
	return strings.Join(levels, "/")  // not [path.Join] to escape redundant filepath behaviours
}

// Next updates the current heading info.
// PANIC when facing skipping level which shall have been avoided outside.
func (hh *HeadingHelper) Next(headingLevelStartsFromOne int, text string) {
	neo := headingLevelStartsFromOne - 1
	if neo > hh.recent && neo-hh.recent > 1 {
		panic(fmt.Errorf("skipping level %d=>%d %v", hh.recent, neo, hh.memory))
	}
	hh.memory[neo] = text
	hh.recent = neo
}

func (f *Filter) save(indent string, node HasLines, source []byte) {
	if f.Dump {
		fmt.Println(indent + Extract(source, node))
	}
	f.Result = append(f.Result, Text{
		Item: Extract(source, node),
		Path: f.hh.Path(),
	})
}

func (f *Filter) traverse(source []byte, node ast.Node, depth int) error {
	indent := strings.Repeat("  ", depth)
	if f.Dump {
		fmt.Printf(indent+"node %v %v\n", node.Type(), node.Kind())
	}
	switch node := node.(type) {
	case *ast.Heading:
		f.save(indent, node, source)
		f.hh.Next(node.Level, Extract(source, node))
	case *ast.Paragraph:
		f.save(indent, node, source)
	case *east.TableCell:
		f.save(indent, node, source)
	case *ast.Blockquote:
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
func FilterText(passage string) []Text {
	md := goldmark.New(goldmark.WithExtensions(extension.Table))
	f := &Filter{Dump: false}
	md.SetRenderer(f)
	err := md.Convert([]byte(passage), os.Stdout)
	if err != nil {
		panic(err)
	}
	return f.Result
}
