package loglint

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "Checks for log formatting issues",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var logMethods = map[string]bool{
	"Info":  true,
	"Error": true,
	"Warn":  true,
	"Debug": true,
	"Trace": true,
	"Fatal": true,
	"Panic": true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkLogCall(pass, call)
	})

	return nil, nil
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	var receiverName string
	switch x := sel.X.(type) {
	case *ast.Ident:
		receiverName = x.Name
	case *ast.SelectorExpr:
		return // для случая типа zap.L().Info(...) — пока игнорируем
	default:
		return
	}

	if receiverName != "log" && receiverName != "slog" {
		return
	}

	methodName := sel.Sel.Name
	if !logMethods[methodName] {
		return
	}

	if len(call.Args) == 0 {
		return
	}

	firstArg := call.Args[0]
	checkMessage(pass, firstArg, call.Pos())
}

func checkMessage(pass *analysis.Pass, arg ast.Expr, pos token.Pos) {
	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	if len(msg) == 0 {
		return
	}

	firstRune := []rune(msg)[0]
	if !unicode.IsLower(firstRune) {
		pass.Report(analysis.Diagnostic{
			Pos:      pos,
			Message:  "log message should start with a lowercase letter",
			Category: "loglint",
		})
	}
}
