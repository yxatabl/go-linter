package loglint

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"unicode"

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

var logMethods = map[string]bool{
	"Info":    true,
	"Error":   true,
	"Warn":    true,
	"Debug":   true,
	"Trace":   true,
	"Fatal":   true,
	"Panic":   true,
	"Print":   true,
	"Println": true,
	"Printf":  true,
}

var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bpassword\b`),
	regexp.MustCompile(`(?i)\bpasswd\b`),
	regexp.MustCompile(`(?i)\btoken\b`),
	regexp.MustCompile(`(?i)api[_-]?key`),
	regexp.MustCompile(`(?i)\bsecret\b`),
	regexp.MustCompile(`(?i)\bcredential\b`),
}

var (
	nonASCIIPattern = regexp.MustCompile(`[^\x00-\x7F]`)
)

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

func startsWithLowercase(s string) bool {
	if len(s) == 0 {
		return true
	}

	firstRune := []rune(s)[0]
	if !unicode.IsLetter(firstRune) {
		return true
	}

	return unicode.IsLower(firstRune)
}

func isEnglishOnly(s string) bool {
	if len(s) == 0 {
		return true
	}

	return !nonASCIIPattern.MatchString(s)
}

func hasSpecialCharsOrEmojis(s string) bool {
	for _, r := range s {
		if r > 0xFFFF {
			return true
		}
	}

	if strings.Contains(s, "!!!") ||
		strings.Contains(s, "???") ||
		strings.Contains(s, "...") {
		return true
	}

	return false
}

func containsSensitiveData(s string) bool {
	sLower := strings.ToLower(s)
	for _, pattern := range defaultSensitivePatterns {
		if pattern.MatchString(sLower) {
			return true
		}
	}
	return false
}

type LogCallInfo struct {
	Receiver string
	Method   string
}
