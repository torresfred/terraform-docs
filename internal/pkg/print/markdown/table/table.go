package table

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/segmentio/terraform-docs/internal/pkg/doc"
	"github.com/segmentio/terraform-docs/internal/pkg/print"
	"github.com/segmentio/terraform-docs/internal/pkg/print/markdown"
	"github.com/segmentio/terraform-docs/internal/pkg/settings"
)

// Print prints a document as Markdown tables.
func Print(document *doc.Doc, settings settings.Settings) (string, error) {
	var buffer bytes.Buffer

	if document.HasComment() {
		printComment(&buffer, document.Comment, settings)
	}

	if document.HasInputs() {
		if settings.Has(print.WithSortByName) {
			if settings.Has(print.WithSortInputsByRequired) {
				doc.SortInputsByRequired(document.Inputs)
			} else {
				doc.SortInputsByName(document.Inputs)
			}
		}

		printInputs(&buffer, document.Inputs, settings)
	}

	if document.HasOutputs() {
		if settings.Has(print.WithSortByName) {
			doc.SortOutputsByName(document.Outputs)
		}

		if document.HasInputs() {
			buffer.WriteString("\n")
		}

		printOutputs(&buffer, document.Outputs, settings)
	}

	return markdown.Sanitize(buffer.String()), nil
}

func getInputDefaultValue(input *doc.Input, settings settings.Settings) string {
	var result = "n/a"

	if input.HasDefault() {
		result = fmt.Sprintf("`%s`", print.GetPrintableValue(input.Default, settings, false))
	}

	return result
}

func printComment(buffer *bytes.Buffer, comment string, settings settings.Settings) {
	buffer.WriteString(fmt.Sprintf("%s\n", comment))
}

func printInputs(buffer *bytes.Buffer, inputs []doc.Input, settings settings.Settings) {
	buffer.WriteString("## Inputs\n\n")
	buffer.WriteString("| Name | Description | Type | Default |")

	if settings.Has(print.WithRequired) {
		buffer.WriteString(" Required |\n")
	} else {
		buffer.WriteString("\n")
	}

	buffer.WriteString("|------|-------------|:----:|:-----:|")

	if settings.Has(print.WithRequired) {
		buffer.WriteString(":-----:|\n")
	} else {
		buffer.WriteString("\n")
	}

	for _, input := range inputs {
		buffer.WriteString(
			fmt.Sprintf("| %s | %s | %s | %s |",
				strings.Replace(input.Name, "_", "\\_", -1),
				markdown.ConvertMultiLineText(input.Description),
				input.Type,
				getInputDefaultValue(&input, settings)))

		if settings.Has(print.WithRequired) {
			buffer.WriteString(fmt.Sprintf(" %v |\n", printIsInputRequired(&input)))
		} else {
			buffer.WriteString("\n")
		}
	}
}

func printIsInputRequired(input *doc.Input) string {
	if input.IsRequired() {
		return "yes"
	}

	return "no"
}

func printOutputs(buffer *bytes.Buffer, outputs []doc.Output, settings settings.Settings) {
	buffer.WriteString("## Outputs\n\n")
	buffer.WriteString("| Name | Description |\n")
	buffer.WriteString("|------|-------------|\n")

	for _, output := range outputs {
		buffer.WriteString(
			fmt.Sprintf("| %s | %s |\n",
				strings.Replace(output.Name, "_", "\\_", -1),
				markdown.ConvertMultiLineText(output.Description)))
	}
}
