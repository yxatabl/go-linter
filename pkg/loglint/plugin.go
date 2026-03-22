package loglint

import "golang.org/x/tools/go/analysis"

func NewPlugin() *analysis.Analyzer {
	return Analyzer
}
