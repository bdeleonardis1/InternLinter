package linter

import (
	"go/parser"
	"go/ast"
	"go/token"

	"github.com/waigani/diffparser"
	
	"strings"
)

// FindFlaws returns all of the todo comments an dprint statements introduced in the diff
func FindFlaws(diff *diffparser.Diff) (todos map[int]string, prints map[int]string, err error) {
	todos = make(map[int]string)
	prints = make(map[int]string)
	for _, file := range diff.Files {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, file.NewName, nil, parser.ParseComments)
		if err != nil {
			return nil, nil, err
		}

		parsedComments := make(map[int]string)
		for _, comment := range parsedFile.Comments {
			if strings.Contains(comment.Text(), "TODO") {
				parsedComments[fset.Position(comment.Pos()).Line] = strings.TrimSpace(comment.Text())
			}
		}

		parsedPrintlns := make(map[int]string)
		ast.Inspect(parsedFile, func(x ast.Node) bool {
			c, ok := x.(*ast.CallExpr)
			if !ok {
				return true
			}
			s, ok := c.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			i, ok := s.X.(*ast.Ident)
			if !ok {
				return true
			}
			if i.Name == "fmt" && s.Sel.String() == "Println" {
				parsedPrintlns[fset.Position(x.Pos()).Line] = "fmt.Println"
			}
			return false
		})

		for _, hunk := range file.Hunks {
			for _, line := range hunk.NewRange.Lines {
				if line.Mode == diffparser.ADDED {
					lineContent := strings.TrimSpace(line.Content)

					commentText, ok := parsedComments[line.Number]
					if ok && len(lineContent) > 4 && lineContent[0:2] == "//" && commentText == lineContent[3:] {
						todos[line.Number] = lineContent
					}

					_, ok = parsedPrintlns[line.Number]
					if ok && len(lineContent) > 11 && lineContent[0:11] == "fmt.Println" {
						prints[line.Number] = lineContent
					}
				}
			}
		}
	}

	return todos, prints, nil
}