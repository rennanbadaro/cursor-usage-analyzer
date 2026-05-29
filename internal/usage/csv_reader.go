package usage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	columnDate        = "Date"
	columnModel       = "Model"
	columnTotalTokens = "Total Tokens"
	columnCost        = "Cost"

	initialRecordsCap = 64
)

// Record holds a single parsed row from a Cursor usage CSV.
type Record struct {
	Day         string
	Model       string
	TotalTokens int64
	Cost        float64
}

// ReadCSV parses Cursor usage records from the given reader,
// which should provide CSV content with Date, Model, Total Tokens,
// and Cost columns.
func ReadCSV(r io.Reader) ([]Record, error) {
	reader := csv.NewReader(r)

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("csv is empty")
		}
		return nil, fmt.Errorf("read header: %w", err)
	}

	indices, err := requiredColumnIndices(header)
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0, initialRecordsCap)
	rowNumber := 1

	for {
		rowNumber++
		row, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("read row %d: %w", rowNumber, err)
		}

		dateRaw, err := rowValue(row, indices[columnDate], columnDate)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNumber, err)
		}

		day, err := parseDay(dateRaw)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid date %q: %w", rowNumber, dateRaw, err)
		}

		modelRaw, err := rowValue(row, indices[columnModel], columnModel)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNumber, err)
		}

		model := strings.TrimSpace(modelRaw)
		if model == "" {
			return nil, fmt.Errorf("row %d: model is empty", rowNumber)
		}

		costValue, err := rowValue(row, indices[columnCost], columnCost)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNumber, err)
		}

		costRaw := strings.TrimSpace(costValue)
		cost, err := parseCost(costRaw)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid cost %q: %w", rowNumber, costRaw, err)
		}

		totalTokensValue, err := rowValue(row, indices[columnTotalTokens], columnTotalTokens)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNumber, err)
		}

		totalTokens, err := parseTotalTokens(strings.TrimSpace(totalTokensValue), cost)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNumber, err)
		}

		records = append(records, Record{
			Day:         day,
			Model:       model,
			TotalTokens: totalTokens,
			Cost:        cost,
		})
	}

	return records, nil
}

// ReadCSVFile opens the file at path and parses its Cursor usage records.
func ReadCSVFile(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()
	return ReadCSV(file)
}

func requiredColumnIndices(header []string) (map[string]int, error) {
	indexByName := make(map[string]int, len(header))
	for i, name := range header {
		indexByName[strings.TrimSpace(name)] = i
	}

	columns := []string{columnDate, columnModel, columnTotalTokens, columnCost}
	indices := make(map[string]int, len(columns))

	for _, column := range columns {
		index, ok := indexByName[column]
		if !ok {
			return nil, fmt.Errorf("missing required column %q", column)
		}
		indices[column] = index
	}

	return indices, nil
}

func parseDay(raw string) (string, error) {
	parsed, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(raw))
	if err != nil {
		return "", err
	}
	return parsed.Format("2006-01-02"), nil
}

func rowValue(row []string, index int, column string) (string, error) {
	if index < 0 || index >= len(row) {
		return "", fmt.Errorf("missing value for %q", column)
	}
	return row[index], nil
}

func parseCost(raw string) (float64, error) {
	if raw == "-" || strings.EqualFold(raw, "free") {
		return 0, nil
	}
	raw = strings.TrimPrefix(raw, "$")
	return strconv.ParseFloat(raw, 64)
}

func parseTotalTokens(raw string, cost float64) (int64, error) {
	if raw == "" {
		if cost != 0 {
			return 0, fmt.Errorf("total tokens missing but cost is non-zero")
		}
		return 0, nil
	}

	totalTokens, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid total tokens %q: %w", raw, err)
	}
	return totalTokens, nil
}
