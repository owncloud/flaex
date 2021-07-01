package parsers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

// ParsedOption represents a single Configuration option
type ParsedOption struct {
	Name    string
	Env     []string
	Usage   string
	Default string
	FnName  string
}

// ParsedOptions represents a set of Configuration options
type ParsedOptions []ParsedOption

type flagSetVisitor struct {
	currentFn     string
	parsedOptions []ParsedOption
}

// ParseFlagsetFile parses a single go file for Configuration options
func ParseFlagsetFile(path string) (ParsedOptions, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, 0)
	if err != nil {
		return nil, err
	}

	fv := &flagSetVisitor{}
	ast.Walk(fv, f)

	return fv.parsedOptions, nil
}

// ParseFlagsetDir parses a folder for go files with Configuration options
func ParseFlagsetDir(path string) (ParsedOptions, error) {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, path, nil, 0)
	if err != nil {
		return nil, err
	}

	fv := &flagSetVisitor{}
	for _, pkg := range pkgs {
		ast.Walk(fv, pkg)
	}

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
						o.Default = getValue(kvExpr)
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

func getValue(kvExpr ast.Expr) (value string) {
	if kv, ok := kvExpr.(*ast.KeyValueExpr); ok {
		v := kv.Value
		if fun, ok := v.(*ast.CallExpr); ok {
			if ce, ok := fun.Fun.(*ast.SelectorExpr); ok {
				if strings.HasPrefix(ce.Sel.Name, "OverrideDefault") {
					if len(fun.Args) == 2 {
						if val, ok := fun.Args[1].(*ast.BasicLit); ok {
							return val.Value
						}
					}
				}
			}
		}
		return exprToStr(kv.Value)
	}
	return value
}
