package report

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
	"strings"
)

type diagnosticRendering struct {
	diagnostic   Diagnostic
	color        *color.Color
	buffer       *strings.Builder
	regionInLine input.Region
	line         input.Line
}

func newDiagnosticRendering(
	diagnostic Diagnostic,
	color *color.Color,
	lineMap *linemap.LineMap) diagnosticRendering {

	line := findLineForDiagnostic(diagnostic, lineMap)
	region := createRelativeRegion(diagnostic.TextRange.Range, line)
	return diagnosticRendering{
		diagnostic:   diagnostic,
		color:        color,
		buffer:       &strings.Builder{},
		regionInLine: region,
		line:         line,
	}
}

func findLineForDiagnostic(diagnostic Diagnostic, lineMap *linemap.LineMap) input.Line {
	textRange := diagnostic.TextRange.Range
	lineIndex := textRange.BeginPosition.Line
	return lineMap.LineAtIndex(input.LineIndex(lineIndex - 1))
}

func createRelativeRegion(region PositionRange, line input.Line) input.Region {
	length := findFixedLength(line)
	begin := region.BeginPosition.Column
	if region.EndPosition.Line == region.BeginPosition.Line &&
		region.EndPosition.Offset != region.BeginPosition.Offset {
		end := minimum(region.EndPosition.Column, length+1)
		return input.CreateRegion(input.Offset(begin), input.Offset(end))
	}
	return input.CreateRegion(input.Offset(begin), input.Offset(length)+1)
}

func findFixedLength(line input.Line) int {
	if strings.HasSuffix(line.Text, "\n") {
		return len(line.Text) - 1
	}
	return len(line.Text) - 1
}

func minimum(left, right int) int {
	if left > right {
		return right
	}
	return left
}

func (rendering *diagnosticRendering) print() string {
	rendering.buffer.WriteString(rendering.description())
	rendering.buffer.WriteString(rendering.lineInformation())
	rendering.buffer.WriteString("| ")
	rendering.buffer.WriteString(rendering.highlightRegionInLine())
	if underscore := rendering.underscoreRegion(); isBlank(underscore) {
		rendering.buffer.WriteString("| ")
		rendering.buffer.WriteString(underscore)
	}
	return rendering.buffer.String()
}

func (rendering *diagnosticRendering) description() string {
	name := rendering.color.Sprintf("[%s]", rendering.diagnostic.Kind)
	return fmt.Sprintf("%s %s\n", name, rendering.diagnostic.Message)
}

func (rendering *diagnosticRendering) lineInformation() string {
	path := rendering.diagnostic.TextRange.File
	textRange := rendering.diagnostic.TextRange.Range
	return fmt.Sprintf(
		"in %s %d:%d\n",
		path,
		textRange.BeginPosition.Line,
		textRange.BeginPosition.Column)
}

func isBlank(text string) bool {
	for _, character := range text {
		if character != ' ' {
			return false
		}
	}
	return true
}

func (rendering *diagnosticRendering) highlightRegionInLine() string {
	if rendering.isFullLineInRegion() || len(rendering.line.Text) == 0 {
		colored := rendering.color.Sprint(rendering.line.Text)
		return withNewLine(colored)
	}
	fullText := rendering.line.Text
	region := rendering.regionInLine
	end := minimum(len(fullText), int(region.End()))
	codeBeforeRegion := fullText[:region.Begin()]
	codeInRegion := fullText[region.Begin():end]
	codeAfterRegion := fullText[end:]
	coloredCodeInRegion := rendering.color.Sprint(codeInRegion)
	return withNewLine(codeBeforeRegion + coloredCodeInRegion + codeAfterRegion)
}

func withNewLine(text string) string {
	if strings.HasSuffix(text, "\n") {
		return text
	}
	return text + "\n"
}

func (rendering *diagnosticRendering) underscoreRegion() string {
	offsetToRegion := int(rendering.regionInLine.Begin())
	offsetText := strings.Repeat(" ", offsetToRegion)
	end := minimum(len(rendering.line.Text), int(rendering.regionInLine.End()))
	regionSize := end - int(rendering.regionInLine.Begin())
	underscore := strings.Repeat("^", regionSize)
	coloredUnderscore := rendering.color.Sprint(underscore)
	return offsetText + coloredUnderscore
}

func (rendering *diagnosticRendering) isFullLineInRegion() bool {
	return rendering.regionInLine.Begin() == 0 &&
		rendering.regionInLine.End() == rendering.line.Length
}
