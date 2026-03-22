package loglint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "Checks for log formatting and content issues",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

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
	logInfo := identifyLogCall(call)
	if logInfo == nil {
		return
	}

	msg := extractMessage(call)
	if msg == "" {
		return
	}

	applyRules(pass, call, msg)
}

func identifyLogCall(call *ast.CallExpr) *LogCallInfo {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	var receiverName string
	switch x := sel.X.(type) {
	case *ast.Ident:
		receiverName = x.Name
	default:
		return nil
	}

	if receiverName != "log" && receiverName != "slog" {
		return nil
	}

	methodName := sel.Sel.Name
	if !logMethods[methodName] {
		return nil
	}

	return &LogCallInfo{
		Receiver: receiverName,
		Method:   methodName,
	}
}

func applyRules(pass *analysis.Pass, call *ast.CallExpr, msg string) {
	hasEmojiOrSpecial := hasSpecialCharsOrEmojis(msg)

	// First rule
	if !startsWithLowercase(msg) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message should start with a lowercase letter",
			Category: "loglint",
		})
	}

	// Second rule
	if !hasEmojiOrSpecial && !isEnglishOnly(msg) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message must be in English only",
			Category: "loglint",
		})
	}

	// Third rule
	if hasEmojiOrSpecial {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message should not contain special characters or emojis",
			Category: "loglint",
		})
	}

	// Fourth rule
	if containsSensitiveData(msg) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message should not contain sensitive data",
			Category: "loglint",
		})
	}
}
