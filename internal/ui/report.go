package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rennanbadaro/cursor-usage-analyzer/internal/usage"
)

const dividerWidth = 40

type reportRow struct {
	label  string
	tokens string
	share  string
	cost   string
}

// Render formats the given Summary into a styled terminal report
// with daily usage, model breakdown, and overall totals.
func Render(summary usage.Summary) string {
	daily := summary.SortedDayTotals()
	models := summary.SortedModelTotals()

	divider := MutedStyle.Render(strings.Repeat("─", dividerWidth))

	dailyRows := make([]reportRow, 0, len(daily))
	for _, d := range daily {
		dailyRows = append(dailyRows, reportRow{
			label:  d.Day,
			tokens: formatIntWithCommas(d.TotalTokens) + " tokens",
			cost:   fmt.Sprintf("$%.2f", d.TotalCost),
		})
	}

	modelRows := make([]reportRow, 0, len(models))
	for _, m := range models {
		tokenShare := 0.0
		if summary.TotalTokens > 0 {
			tokenShare = (float64(m.TotalTokens) / float64(summary.TotalTokens)) * 100
		}
		modelRows = append(modelRows, reportRow{
			label:  m.Model,
			tokens: formatIntWithCommas(m.TotalTokens) + " tokens",
			share:  fmt.Sprintf("%.2f%%", tokenShare),
			cost:   fmt.Sprintf("$%.2f", m.TotalCost),
		})
	}

	overallRows := []reportRow{
		{label: "Tokens", tokens: formatIntWithCommas(summary.TotalTokens)},
		{label: "Cost", tokens: fmt.Sprintf("$%.2f", summary.TotalCost)},
	}

	var content strings.Builder
	content.WriteString(TitleStyle.Render("Usage Summary"))
	content.WriteString("\n\n")
	content.WriteString(TitleStyle.Render("Daily Usage"))
	content.WriteString("\n")
	content.WriteString(renderRows(dailyRows))
	content.WriteString("\n\n")
	content.WriteString(divider)
	content.WriteString("\n\n")
	content.WriteString(TitleStyle.Render("Model Breakdown"))
	content.WriteString("\n")
	content.WriteString(renderRows(modelRows))
	content.WriteString("\n\n")
	content.WriteString(divider)
	content.WriteString("\n\n")
	content.WriteString(TitleStyle.Render("Overall Summary"))
	content.WriteString("\n")
	content.WriteString(renderRows(overallRows))

	return "\n" + PanelStyle.Render(content.String()) + "\n"
}

func renderRows(rows []reportRow) string {
	labelWidth := 0
	tokenWidth := 0
	shareWidth := 0
	hasShare := false
	hasCost := false

	for _, row := range rows {
		if lipgloss.Width(row.label) > labelWidth {
			labelWidth = lipgloss.Width(row.label)
		}
		if lipgloss.Width(row.tokens) > tokenWidth {
			tokenWidth = lipgloss.Width(row.tokens)
		}
		if lipgloss.Width(row.share) > shareWidth {
			shareWidth = lipgloss.Width(row.share)
		}
		if row.share != "" {
			hasShare = true
		}
		if row.cost != "" {
			hasCost = true
		}
	}

	rendered := make([]string, 0, len(rows))
	for _, row := range rows {
		line := LabelStyle.Render(padRight(row.label, labelWidth)) + "  " + ValueStyle.Render(padRight(row.tokens, tokenWidth))
		if hasShare {
			shareValue := row.share
			if shareValue == "" {
				shareValue = "-"
			}
			line += "  " + MutedStyle.Render("|") + "  " + ValueStyle.Render(padRight(shareValue, shareWidth))
		}
		if hasCost {
			costValue := row.cost
			if costValue == "" {
				costValue = "-"
			}
			line += "  " + MutedStyle.Render("|") + "  " + ValueStyle.Render(costValue)
		}
		rendered = append(rendered, line)
	}

	return strings.Join(rendered, "\n")
}

func padRight(value string, width int) string {
	current := lipgloss.Width(value)
	if current >= width {
		return value
	}
	return value + strings.Repeat(" ", width-current)
}

func formatIntWithCommas(value int64) string {
	if value == 0 {
		return "0"
	}

	negative := value < 0
	if negative {
		value = -value
	}

	digits := fmt.Sprintf("%d", value)
	n := len(digits)

	var builder strings.Builder
	if negative {
		builder.WriteByte('-')
	}

	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			builder.WriteByte(',')
		}
		builder.WriteByte(digits[i])
	}

	return builder.String()
}
