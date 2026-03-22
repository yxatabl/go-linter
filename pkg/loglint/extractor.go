package loglint

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func extractMessage(call *ast.CallExpr) string {
	if len(call.Args) == 0 {
		return ""
	}

	firstArg := call.Args[0]
	return extractStringFromExpr(firstArg)
}

func extractStringFromExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return unquoteString(e.Value)
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			left := extractStringFromExpr(e.X)
			right := extractStringFromExpr(e.Y)
			return left + right
		}
	}
	return ""
}

func unquoteString(s string) string {
	if strings.HasPrefix(s, "`") {
		return s[1 : len(s)-1]
	}

	unquoted, err := strconv.Unquote(s)
	if err != nil {
		return ""
	}
	return unquoted
}
