package usage

import (
	"math"
	"strings"
	"testing"
)

const testCSV = `Date,Model,Total Tokens,Cost
2026-02-18T10:00:00.000Z,claude-4.6-opus-high-thinking,500000,5.00
2026-02-18T11:00:00.000Z,gpt-5.3-codex,300000,2.50
2026-02-19T09:00:00.000Z,claude-4.6-opus-high-thinking,700000,8.00
2026-02-19T14:00:00.000Z,gpt-5.3-codex,200000,1.50
`

func TestReadCSVAndAggregate(t *testing.T) {
	records, err := ReadCSV(strings.NewReader(testCSV))
	if err != nil {
		t.Fatalf("ReadCSV() error = %v", err)
	}
	if len(records) != 4 {
		t.Fatalf("ReadCSV() returned %d records, want 4", len(records))
	}

	summary := Aggregate(records)

	if len(summary.TokensByDay) != 2 {
		t.Fatalf("expected 2 days, got %d", len(summary.TokensByDay))
	}
	if summary.TotalTokens != 1700000 {
		t.Fatalf("expected 1700000 total tokens, got %d", summary.TotalTokens)
	}
	if !almostEqual(summary.TotalCost, 17.0) {
		t.Fatalf("expected 17.0 total cost, got %f", summary.TotalCost)
	}

	var sumByDay int64
	for _, total := range summary.TokensByDay {
		sumByDay += total
	}
	if sumByDay != summary.TotalTokens {
		t.Fatalf("sum by day (%d) != total tokens (%d)", sumByDay, summary.TotalTokens)
	}

	var sumByModel int64
	for _, total := range summary.TokensByModel {
		sumByModel += total
	}
	if sumByModel != summary.TotalTokens {
		t.Fatalf("sum by model (%d) != total tokens (%d)", sumByModel, summary.TotalTokens)
	}
}

func TestReadCSVMissingRequiredColumn(t *testing.T) {
	input := "Date,Model\n2026-02-19T19:35:35.883Z,gpt-5.3-codex\n"

	_, err := ReadCSV(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expected error for missing required column")
	}
	if !strings.Contains(err.Error(), "missing required column") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadCSVInvalidTotalTokens(t *testing.T) {
	input := "Date,Model,Total Tokens,Cost\n2026-02-19T19:35:35.883Z,gpt-5.3-codex,not-a-number,0.10\n"

	_, err := ReadCSV(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expected error for invalid total tokens")
	}
	if !strings.Contains(err.Error(), "invalid total tokens") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadCSVEmptyTotalTokensWithZeroCost(t *testing.T) {
	const header = "Date,User,Kind,Model,Total Tokens,Cost"
	rows := []struct {
		name string
		row  string
	}{
		{
			name: "aborted not charged",
			row:  `2026-05-14T21:47:30.360Z,user@example.com,"Aborted, Not Charged",claude-opus-4-7-thinking-high,,Free`,
		},
		{
			name: "errored no charge",
			row:  `2026-05-13T21:01:39.686Z,user@example.com,"Errored, No Charge",gpt-5.5-high,,Free`,
		},
	}

	for _, tc := range rows {
		t.Run(tc.name, func(t *testing.T) {
			input := header + "\n" + tc.row + "\n"
			records, err := ReadCSV(strings.NewReader(input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(records) != 1 {
				t.Fatalf("expected 1 record, got %d", len(records))
			}
			if records[0].TotalTokens != 0 {
				t.Fatalf("expected 0 tokens, got %d", records[0].TotalTokens)
			}
			if records[0].Cost != 0 {
				t.Fatalf("expected 0 cost, got %f", records[0].Cost)
			}
		})
	}
}

func TestReadCSVEmptyTotalTokensWithNonZeroCost(t *testing.T) {
	input := "Date,Model,Total Tokens,Cost\n2026-02-19T19:35:35.883Z,gpt-5.3-codex,,1.00\n"

	_, err := ReadCSV(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expected error for empty total tokens with non-zero cost")
	}
	if !strings.Contains(err.Error(), "total tokens missing but cost is non-zero") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadCSVInvalidCost(t *testing.T) {
	input := "Date,Model,Total Tokens,Cost\n2026-02-19T19:35:35.883Z,gpt-5.3-codex,12345,not-a-number\n"

	_, err := ReadCSV(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expected error for invalid cost")
	}
	if !strings.Contains(err.Error(), "invalid cost") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadCSVDashCostAsZero(t *testing.T) {
	input := "Date,Model,Total Tokens,Cost\n2026-03-10T19:22:52.869Z,claude-4.6-opus-max-thinking,3117043,-\n"

	records, err := ReadCSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Cost != 0 {
		t.Fatalf("expected dash cost to parse as 0, got %f", records[0].Cost)
	}
}

func TestReadCSVFreeCostAsZero(t *testing.T) {
	input := "Date,Model,Total Tokens,Cost\n2026-03-11T18:52:15.043Z,gpt-5.3-codex-high,948732,Free\n"

	records, err := ReadCSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Cost != 0 {
		t.Fatalf("expected free cost to parse as 0, got %f", records[0].Cost)
	}
}

func TestReadCSVDollarPrefixCost(t *testing.T) {
	input := "Date,Model,Total Tokens,Cost\n2026-03-10T19:22:52.869Z,gpt-5.3-codex,500000,$1.23\n"

	records, err := ReadCSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if !almostEqual(records[0].Cost, 1.23) {
		t.Fatalf("expected cost 1.23, got %f", records[0].Cost)
	}
}

func TestSortedDayTotals(t *testing.T) {
	records, err := ReadCSV(strings.NewReader(testCSV))
	if err != nil {
		t.Fatalf("ReadCSV() error = %v", err)
	}

	summary := Aggregate(records)
	days := summary.SortedDayTotals()

	if len(days) != 2 {
		t.Fatalf("expected 2 days, got %d", len(days))
	}
	if days[0].Day >= days[1].Day {
		t.Fatalf("days not sorted ascending: %s >= %s", days[0].Day, days[1].Day)
	}
}

func TestSortedModelTotals(t *testing.T) {
	records, err := ReadCSV(strings.NewReader(testCSV))
	if err != nil {
		t.Fatalf("ReadCSV() error = %v", err)
	}

	summary := Aggregate(records)
	models := summary.SortedModelTotals()

	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
	if models[0].TotalTokens < models[1].TotalTokens {
		t.Fatalf("models not sorted by tokens descending: %d < %d", models[0].TotalTokens, models[1].TotalTokens)
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
