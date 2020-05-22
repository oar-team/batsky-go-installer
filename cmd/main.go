package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	callsToReplace := []string{"Now", "Sleep", "NewTimer", "After", "AfterFunc", "NewTicker", "Tick"}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [files]\n", os.Args[0])
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		// For when we run the program with go run
		if arg == "--" {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, arg, nil, 0)
		if err != nil {
			panic(err)
		}

		var nbImports int
		var timeAlias string
		timeIsImported := false
		var importIndex int
		ast.Inspect(f, func(n ast.Node) bool {
			switch n.(type) {
			case *ast.ImportSpec:
				nbImports++
				if n.(*ast.ImportSpec).Path.Value == "\"time\"" {
					timeIsImported = true
					importIndex = fset.Position(n.(*ast.ImportSpec).Path.ValuePos).Line
					if n.(*ast.ImportSpec).Name != nil {
						timeAlias = n.(*ast.ImportSpec).Name.Name
					}
				}
			}
			return true
		})

		if timeIsImported {
			if nbImports > 1 {
				insertStringToFile(arg, "batsky-time \"github.com/oar-team/batsky-go/time\"", importIndex)
			} else if nbImports == 1 {
				//TODO
			}
		}
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

func file2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

/**
 * Insert sting to n-th line of file.
 * If you want to insert a line, append newline '\n' to the end of the string.
 */
func insertStringToFile(path, str string, index int) error {
	lines, err := file2lines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}
