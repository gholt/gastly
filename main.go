package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(`%s <infile> <outfile> <package> [from=to] ...

Copies a Go source file with rewriting rules.

<infile> will be read.
<outfile> will be written.
<package> will be the new package name.

Each from=to specifies a rewrite rule, replacing any occurrences of "from" with
"to". These will be applied in order, so a leftmost replacement may affect what
the following replacements match.

If "to" begins with "droptype:" it will also drop any type specification that
matches "from" exactly. For example, "NumericType=droptype:int" will replace
"NumericType" with "int" everywhere and it will drop any "type NumericType ..."
specification it finds.
`, os.Args[0])
		os.Exit(1)
	}
	infile := os.Args[1]
	outfile := os.Args[2]
	packageName := os.Args[3]
	r := &rewriter{dropTypes: map[string]struct{}{}}
	for _, arg := range os.Args[4:] {
		s := strings.SplitN(arg, "=", 2)
		if len(s) != 2 {
			fmt.Println("Invalid syntax:", arg)
			os.Exit(1)
		}
		if strings.HasPrefix(s[1], "droptype:") {
			s[1] = s[1][len("droptype:"):]
			r.dropTypes[s[0]] = struct{}{}
		}
		r.translations = append(r.translations, []string{s[0], s[1]})
	}
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, infile, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	astf.Name.Name = packageName
	ast.Walk(r, astf)
	f, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintln(f, "// Automatically generated with:", strings.Join(os.Args, " "))
	fmt.Fprintln(f)
	if err = format.Node(f, fset, astf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type rewriter struct {
	dropTypes    map[string]struct{}
	translations [][]string
}

func (r rewriter) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		var drops []int
		for i, m := range n.Decls {
			if d, ok := m.(*ast.GenDecl); ok {
				for _, k := range d.Specs {
					if ts, ok := k.(*ast.TypeSpec); ok {
						if _, ok := r.dropTypes[ts.Name.Name]; ok {
							drops = append(drops, i)
						}
					}
				}
			}
		}
		for d := len(drops) - 1; d >= 0; d-- {
			if drops[d] == len(n.Decls)-1 {
				n.Decls = n.Decls[:len(n.Decls)-1]
			} else {
				n.Decls = append(n.Decls[:drops[d]], n.Decls[drops[d]+1:]...)
			}
		}
	case *ast.Ident:
		for _, fromto := range r.translations {
			n.Name = strings.Replace(n.Name, fromto[0], fromto[1], -1)
		}
	case *ast.Comment:
		for _, fromto := range r.translations {
			n.Text = strings.Replace(n.Text, fromto[0], fromto[1], -1)
		}
	}
	return r
}
