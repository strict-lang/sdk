package report

import (
	"fmt"
	"io"
	"strings"
)

type renderingOutput struct {
	buffer strings.Builder
	report Report
	diagnosticStats diagnosticStats
}

func NewRenderingOutput(report Report) Output {
	output := &renderingOutput{report: report}
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
	output.printCompletionMessage()
}

func (output *renderingOutput) printCompletionMessage() {
	duration := output.formatDuration()
	suffix := output.createCompletionSuffix()
	message := fmt.Sprintf("completed compilation %s, took %s\n", suffix, duration)
	output.buffer.WriteString(message)
}

func (output *renderingOutput) createCompletionSuffix() string {
	if output.report.Success {
		if output.diagnosticStats.warningCount > 0 {
			return fmt.Sprintf("successfully (with %d warnings)", output.diagnosticStats.warningCount)
		}
		return "successfully"
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