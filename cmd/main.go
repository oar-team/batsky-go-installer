package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	callsToReplace := []string{"Now", "Sleep", "NewTimer", "After", "AfterFunc", "NewTicker", "Tick"}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [files]\n", os.Args[0])
		os.Exit(1)
	}

	for _, path := range os.Args[1:] {
		// For when we run the program with go run
		if path == "--" {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		timeIsImported := false
		timeAlias := "time"
		imports := astutil.Imports(fset, f)
		for _, specs := range imports {
			for _, spec := range specs {
				if spec.Path.Value == "\"time\"" {
					timeIsImported = true
					if spec.Name != nil {
						timeAlias = spec.Name.Name
					}
				}
			}
		}

		fmt.Println(timeAlias)

		if !timeIsImported {
			// Do not even bother
			return
		}

		importIsNeeded := false
		batskyAlias := "batskyTime"
		newAST := astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
			switch c.Node().(type) {
			case *ast.Ident:
				if c.Node().(*ast.Ident).Name == timeAlias {
					p := c.Parent()
					switch p.(type) {
					case *ast.SelectorExpr:
						if isIn(p.(*ast.SelectorExpr).Sel.Name, callsToReplace) {
							replacementNode := c.Node()
							replacementNode.(*ast.Ident).Name = batskyAlias
							c.Replace(replacementNode)
							importIsNeeded = true
						}
					}
				}
			}
			return true
		})

		if importIsNeeded {
			astutil.AddNamedImport(fset, f, batskyAlias, "github.com/oar-team/batsky-go/time")
		}
		if !astutil.UsesImport(newAST.(*ast.File), "time") {
			switch timeAlias {
			case "time":
				if !astutil.DeleteImport(fset, newAST.(*ast.File), "time") {
					panic("Could not remove unused time import")
				}
			default:
				if !astutil.DeleteNamedImport(fset, newAST.(*ast.File), timeAlias, "time") {
					panic("Could not remove unused time import")
				}
			}
		}

		ast.Print(fset, newAST)

		buf := &bytes.Buffer{}
		err = format.Node(buf, fset, newAST)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(path, buf.Bytes(), 0644)
	}
}

func isIn(s string, slice []string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}
