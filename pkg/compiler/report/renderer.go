package report

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
	"io"
	"strings"
)

type renderingOutput struct {
	buffer strings.Builder
	report Report
	diagnosticStats diagnosticStats
	lineMap *linemap.LineMap
}

func NewRenderingOutput(report Report, lineMap *linemap.LineMap) Output {
	output := &renderingOutput{
		report: report,
		lineMap: lineMap,
	}
	output.calculateDiagnosticStats()
	return output
}

type diagnosticStats struct {
	warningCount int
	infoCount int
	errorCount int
}

func (output *renderingOutput) Print(writer io.Writer) error {
	output.render()
	rawOutput := []byte(output.buffer.String())
	_, err := writer.Write(rawOutput)
	return err
}

func (output *renderingOutput) render() {
	for _, diagnostic := range output.report.Diagnostics {
		output.renderDiagnostic(diagnostic)
		output.buffer.WriteString("\n")
	}
	output.printCompletionMessage()
}

func (output *renderingOutput) renderDiagnostic(diagnostic Diagnostic) {
	errorColor := color.New(color.FgRed)
	rendering := newDiagnosticRendering(diagnostic, errorColor, output.lineMap)
	output.buffer.WriteString(rendering.print())
}

var successColor = color.New(color.FgGreen).Add(color.Bold)
var failureColor = color.New(color.FgRed).Add(color.Bold)
var warningColor = color.New(color.FgBlue).Add(color.Bold)

func (output *renderingOutput) printCompletionMessage() {
	duration := output.formatDuration()
	suffix, color := output.createCompletionSuffix()
	message := color.Sprintf("completed compilation %s, took %s\n", suffix, duration)
	output.buffer.WriteString(message)
}

func (output *renderingOutput) createCompletionSuffix() (string, *color.Color) {
	if output.report.Success {
		if output.diagnosticStats.warningCount > 0 {
			warnings := fmt.Sprintf("successfully (with %d warnings)", output.diagnosticStats.warningCount)
			return warnings, warningColor
		}
		return "successfully", successColor
	}
	return output.createFailedCompletionSuffix(), failureColor
}

func (output *renderingOutput) createFailedCompletionSuffix() string {
	if output.diagnosticStats.errorCount == 1 {
		return "(with 1 error)"
	}
	if output.diagnosticStats.errorCount > 0 {
		return fmt.Sprintf("(with %d errors)", output.diagnosticStats.errorCount)
	}
	return "(with errors)"
}

func (output *renderingOutput) formatDuration() string {
	duration := output.report.Time.CalculateDuration()
	return fmt.Sprintf("%.02f seconds", duration.Seconds())
}

func (output *renderingOutput) calculateDiagnosticStats() {
	output.diagnosticStats = diagnosticStats{}
	for _, entry := range output.report.Diagnostics {
		switch entry.Kind {
		case DiagnosticError:
			output.diagnosticStats.errorCount++
			break
		case DiagnosticInfo:
			output.diagnosticStats.infoCount++
			break
		case DiagnosticWarning:
			output.diagnosticStats.warningCount++
			break
		}
	}
}