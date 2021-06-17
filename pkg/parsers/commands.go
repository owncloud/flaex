package parsers

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ParsedCommand represents a single Configuration option
type ParsedCommand struct {
	Name   string
	Usage  string
	Flags  []string
	Type   string
	FnName string
}

// ParsedCommands represents a set of Configuration options
type ParsedCommands []ParsedCommand

type commandVisitor struct {
	currentFn      string
	parsedCommands []ParsedCommand
}

// ParseCommandFile parses a single go file for Configuration options
func ParseCommandFile(path string) (ParsedCommands, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, 0)
	if err != nil {
		return nil, err
	}

	fv := &commandVisitor{}
	ast.Walk(fv, f)

	return fv.parsedCommands, nil
}

// ParseCommandDir parses a folder for go files with Configuration options
func ParseCommandDir(path string) (ParsedCommands, error) {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, path, nil, 0)
	if err != nil {
		return nil, err
	}

	fv := &commandVisitor{}
	for _, pkg := range pkgs {
		ast.Walk(fv, pkg)
	}

	return fv.parsedCommands, nil
}

func (v *commandVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if fn, ok := node.(*ast.FuncDecl); ok {
		v.currentFn = fn.Name.Name
	}

	if comp, ok := node.(*ast.CompositeLit); ok {
		if selExpr, ok := comp.Type.(*ast.SelectorExpr); ok && isCommandType(selExpr) {
			c := ParsedCommand{}
			c.Type = selExpr.Sel.Name
			for _, v := range comp.Elts {
				if kvExpr, ok := v.(*ast.KeyValueExpr); ok {
					key := exprToStr(kvExpr.Key)
					val := exprToStr(kvExpr.Value)
					switch key {
					case "Name":
						c.Name = val
					case "Usage":
						c.Usage = val
					case "Flags":
						c.Flags = getFlagset(kvExpr.Value)
					}
				}
			}
			c.FnName = v.currentFn
			v.parsedCommands = append(v.parsedCommands, c)
		}
	}

	return v
}

func isCommandType(sel *ast.SelectorExpr) bool {
	switch sel.Sel.Name {
	case "App", "Command":
		return true
	}

	return false
}

func getFlagset(flagset ast.Expr) (flagsets []string) {

	if fc, ok := flagset.(*ast.CallExpr); ok {
		if fi, ok := fc.Fun.(*ast.Ident); ok {
			if fi.Name == "append" {
				for _, arg := range fc.Args {
					fs := getFlagset(arg)
					flagsets = append(flagsets, fs...)
				}
			}
		} else if se, ok := fc.Fun.(*ast.SelectorExpr); ok {
			if fi, ok := se.X.(*ast.Ident); ok {
				if fi.Name == "flagset" {
					flagsets = append(flagsets, se.Sel.Name)
				}
			}
		}
	}
	return flagsets
}
