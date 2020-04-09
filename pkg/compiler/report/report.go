package report

import (
	times "time"
)

type Report struct {
	Success     bool         `json:"success"`
	Time        Time         `json:"time"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Time struct {
	Begin      int64 `json:"begin"`
	Completion int64 `json:"completion"`
}

func (time Time) CalculateDuration() times.Duration {
	return times.Duration(time.Completion - time.Begin)
}

type DiagnosticKind string

const (
	DiagnosticError   DiagnosticKind = "error"
	DiagnosticInfo                   = "info"
	DiagnosticWarning                = "warning"
)

type Diagnostic struct {
	TextRange TextRange      `json:"textRange"`
	Message   string         `json:"message"`
	Kind      DiagnosticKind `json:"kind"`
}

type TextRange struct {
	Text     string   `json:"text"`
	Range    PositionRange `json:"range"`
	File      string  `json:"file"`
}

type PositionRange struct {
	BeginPosition Position `json:"beginPosition"`
	EndPosition Position `json:"endPosition"`
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
	Offset int `json:"offset"`
}