package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

type ParsedOption struct {
	Name    string
	Env     []string
	Usage   string
	Default string
	FnName  string
}

type ParsedOptions []ParsedOption

type flagSetVisitor struct {
	currentFn     string
	parsedOptions []ParsedOption
}

func ParseFile(path string) (ParsedOptions, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, 0)
	if err != nil {
		return nil, err
	}

	fv := &flagSetVisitor{}
	ast.Walk(fv, f)

	return fv.parsedOptions, nil
}

func (v *flagSetVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if fn, ok := node.(*ast.FuncDecl); ok {
		v.currentFn = fn.Name.Name
	}

	if comp, ok := node.(*ast.CompositeLit); ok {
		if selExpr, ok := comp.Type.(*ast.SelectorExpr); ok && isFlagType(selExpr) {
			o := ParsedOption{}
			for _, v := range comp.Elts {
				if kvExpr, ok := v.(*ast.KeyValueExpr); ok {
					key := exprToStr(kvExpr.Key)
					val := exprToStr(kvExpr.Value)
					switch key {
					case "Name":
						o.Name = val
					case "Usage":
						o.Usage = val
					case "Value":
						o.Default = val
					case "EnvVars":
						if comp, ok := kvExpr.Value.(*ast.CompositeLit); ok {
							for _, envVar := range comp.Elts {
								ev := exprToStr(envVar)
								o.Env = append(o.Env, ev)
							}
						}
					}
				}
			}

			o.FnName = v.currentFn
			v.parsedOptions = append(v.parsedOptions, o)
		}
	}

	return v
}

func isFlagType(sel *ast.SelectorExpr) bool {
	switch sel.Sel.Name {
	case "StringFlag", "BoolFlag", "SliceFlag":
		return true
	}

	return false
}

func exprToStr(x ast.Expr) string {
	exprStr := types.ExprString(x)
	return strings.Trim(exprStr, "\"")
}
