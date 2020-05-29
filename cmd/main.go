package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

var callsToReplace []string = []string{"Now", "Sleep", "NewTimer", "After", "AfterFunc", "NewTicker", "Tick", "Timer", "Ticker"}
var timeAlias string = "time"
var timePath string = "time"
var batskyAlias string = "batskyTime"
var batskyPath string = "github.com/oar-team/batsky-go/time"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [path1 path2...] [--not [path3 path4...]]\n", os.Args[0])
		os.Exit(1)
	}

	var paths []string
	var ignorePaths []string
	not := false
	for _, str := range os.Args[1:] {
		switch str {
		case "--":
			// this happens when this code is run with go run
		case "--not":
			not = true
		default:
			if not {
				ignorePaths = append(ignorePaths, str)
			} else {
				paths = append(paths, str)
			}
		}
	}

	fmt.Println("This action will replace all calls to \"time\"'s", callsToReplace, "in the given directories and files")
	warning := true
	for warning {
		fmt.Println("This action is irreversible. Do you wish to continue? [y/N/show-files]")

		text := ""
		fmt.Scanln(&text)
		switch text {
		case "N", "":
			return
		case "y":
			warning = false
		case "show-files":
			fmt.Println("These files will be modified :")
			fmt.Println()
			walkDirs(paths, ignorePaths, true)
		}
		fmt.Println()
	}
	walkDirs(paths, ignorePaths, false)
}

/*
Looks for specific calls to "time"'s functions (see the list in the function
warning message in order to change those calls to call
github.com/oar-team/batsky-go/time functions instead.

Leaves the rest of time objects and functions as is.

dryRun : only shows the files it has an effect on without actually doing anything
*/
func walkDirs(paths, ignorePaths []string, dryRun bool) error {
	var err error
	for _, path := range paths {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}
			if !info.IsDir() && !isPathIn(path, ignorePaths) && filepath.Ext(path) == ".go" {
				if dryRun {
					fmt.Println(path)
				} else {
					searchAndReplace(path)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func searchAndReplace(path string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	timeIsImported := false
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

	if !timeIsImported {
		// Do not even bother
		return nil
	}

	importIsNeeded := false
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
		astutil.AddNamedImport(fset, f, batskyAlias, batskyPath)
	}
	if !astutil.UsesImport(newAST.(*ast.File), timePath) {
		switch timeAlias {
		case timePath:
			if !astutil.DeleteImport(fset, newAST.(*ast.File), timePath) {
				return errors.New("Could not remove unused time import")
			}
		default:
			if !astutil.DeleteNamedImport(fset, newAST.(*ast.File), timeAlias, timePath) {
				return errors.New("Could not remove unused time import")
			}
		}
	}

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, newAST)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, buf.Bytes(), 0644)
}

/*
Returns whether s is in the slice
*/
func isIn(s string, slice []string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}

/*
Returns whether the specified path is a child of any of the paths listed in the
input slice
*/
func isPathIn(path string, paths []string) bool {
	for _, e := range paths {
		if strings.Contains(path, e) {
			return true
		}
	}
	return false
}
