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

var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)passwd`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)api[_-]?key`),
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)credential`),
	regexp.MustCompile(`(?i)auth`),
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

	msg := extractMessage(call, logInfo)
	if msg == "" {
		return
	}

	applyRules(pass, call, msg, logInfo)
}

func extractMessage(call *ast.CallExpr, info *LogCallInfo) string {
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
			if left != "" && right != "" {
				return left + right
			}
		}
	case *ast.Ident:
		// переменная-сообщение - пока игнорируем
		return ""
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

func applyRules(pass *analysis.Pass, call *ast.CallExpr, msg string, info *LogCallInfo) {
	if !startswithLowercase(msg) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message should start with a lowercase letter",
			Category: "loglint",
		})
	}

	if !isEnglishOnly(msg) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message must be in English only",
			Category: "loglint",
		})
	}

	if hasSpecialCahrsOrEmojis(msg) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message should not contain special characters or emojis",
			Category: "loglint",
		})
	}

	if containsSensitiveData(msg) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Message:  "log message should not contain sensitive data",
			Category: "loglint",
		})
	}
}

func startswithLowercase(s string) bool {
	if len(s) == 0 {
		return true
	}

	firstRune := []rune(s)[0]
	return unicode.IsLower(firstRune)
}

func isEnglishOnly(s string) bool {
	allowedPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\.,!?\-':;\(\)\[\]{}<>/\\_+=*&^%$#@~` + "`" + `" ]+$`)
	nonEnglishPattern := regexp.MustCompile(`[^\x00-\x7F]`)

	if len(s) == 0 {
		return true
	}

	if nonEnglishPattern.MatchString(s) {
		return false
	}

	return allowedPattern.MatchString(s)
}

func hasSpecialCahrsOrEmojis(s string) bool {
	emojisPattern := regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F1E0}-\x{1F1FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}]`)
	specialCahrs := regexp.MustCompile(`[!?…]{2,}$`)

	if emojisPattern.MatchString(s) {
		return true
	}

	if specialCahrs.MatchString(s) {
		return true
	}

	return false
}

func containsSensitiveData(s string) bool {
	for _, pattern := range defaultSensitivePatterns {
		if pattern.MatchString(s) {
			return true
		}
	}

	return false
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
	case *ast.SelectorExpr:
		return nil // для случая типа zap.L().Info(...) — пока игнорируем
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

type LogCallInfo struct {
	Receiver string
	Method   string
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
