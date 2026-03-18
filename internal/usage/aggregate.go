package usage

import "sort"

// Summary holds aggregated token counts and costs grouped by day and model.
type Summary struct {
	TotalTokens   int64
	TotalCost     float64
	TokensByDay   map[string]int64
	TokensByModel map[string]int64
	CostByDay     map[string]float64
	CostByModel   map[string]float64
}

// DayTotal holds aggregate token and cost data for a single day.
type DayTotal struct {
	Day         string
	TotalTokens int64
	TotalCost   float64
}

// ModelTotal holds aggregate token and cost data for a single model.
type ModelTotal struct {
	Model       string
	TotalTokens int64
	TotalCost   float64
}

// Aggregate computes a Summary from the given records, grouping
// totals by day and by model.
func Aggregate(records []Record) Summary {
	summary := Summary{
		TokensByDay:   make(map[string]int64),
		TokensByModel: make(map[string]int64),
		CostByDay:     make(map[string]float64),
		CostByModel:   make(map[string]float64),
	}

	for _, record := range records {
		summary.TotalTokens += record.TotalTokens
		summary.TotalCost += record.Cost
		summary.TokensByDay[record.Day] += record.TotalTokens
		summary.TokensByModel[record.Model] += record.TotalTokens
		summary.CostByDay[record.Day] += record.Cost
		summary.CostByModel[record.Model] += record.Cost
	}

	return summary
}

// SortedDayTotals returns day totals sorted by date ascending.
func (s Summary) SortedDayTotals() []DayTotal {
	days := make([]DayTotal, 0, len(s.TokensByDay))
	for day, tokens := range s.TokensByDay {
		days = append(days, DayTotal{
			Day:         day,
			TotalTokens: tokens,
			TotalCost:   s.CostByDay[day],
		})
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].Day < days[j].Day
	})

	return days
}

// SortedModelTotals returns model totals sorted by tokens descending,
// with ties broken by model name ascending.
func (s Summary) SortedModelTotals() []ModelTotal {
	models := make([]ModelTotal, 0, len(s.TokensByModel))
	for model, tokens := range s.TokensByModel {
		models = append(models, ModelTotal{
			Model:       model,
			TotalTokens: tokens,
			TotalCost:   s.CostByModel[model],
		})
	}

	sort.Slice(models, func(i, j int) bool {
		if models[i].TotalTokens == models[j].TotalTokens {
			return models[i].Model < models[j].Model
		}
		return models[i].TotalTokens > models[j].TotalTokens
	})

	return models
}
