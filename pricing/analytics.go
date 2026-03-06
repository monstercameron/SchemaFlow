package pricing

import "time"

// CostSummary aggregates request, token, and cost data for a time range.
type CostSummary struct {
	RequestCount            int
	TotalCost               float64
	AverageCostPerRequest   float64
	TotalPromptTokens       int
	TotalCompletionTokens   int
	TotalCachedTokens       int
	TotalReasoningTokens    int
	TotalTokens             int
	AveragePromptTokens     float64
	AverageCompletionTokens float64
	AverageCachedTokens     float64
	AverageReasoningTokens  float64
	AverageTokensPerRequest float64
	AveragePromptCost       float64
	AverageCompletionCost   float64
	AverageCachedCost       float64
	AverageReasoningCost    float64
}

// GetRequestCosts returns the request-level cost history for a time range.
func GetRequestCosts(since time.Time, filters map[string]string) []CostRecord {
	costMutex.RLock()
	defer costMutex.RUnlock()

	records := make([]CostRecord, 0, len(costHistory))
	for _, record := range costHistory {
		if record.Timestamp.Before(since) || !matchesFilters(record, filters) {
			continue
		}
		records = append(records, cloneCostRecord(record))
	}

	return records
}

// GetRequestCost returns a single request-level cost record by request ID.
func GetRequestCost(requestID string) (CostRecord, bool) {
	costMutex.RLock()
	defer costMutex.RUnlock()

	for _, record := range costHistory {
		if record.RequestID == requestID {
			return cloneCostRecord(record), true
		}
	}

	return CostRecord{}, false
}

// GetCostSummary returns aggregate totals and averages for a time range.
func GetCostSummary(since time.Time, filters map[string]string) CostSummary {
	costMutex.RLock()
	defer costMutex.RUnlock()

	var summary CostSummary

	for _, record := range costHistory {
		if record.Timestamp.Before(since) || !matchesFilters(record, filters) {
			continue
		}

		summary.RequestCount++
		summary.TotalCost += record.Cost.TotalCost
		summary.TotalPromptTokens += record.TokenUsage.PromptTokens
		summary.TotalCompletionTokens += record.TokenUsage.CompletionTokens
		summary.TotalCachedTokens += record.TokenUsage.CachedTokens
		summary.TotalReasoningTokens += record.TokenUsage.ReasoningTokens
		summary.TotalTokens += record.TokenUsage.TotalTokens
		summary.AveragePromptCost += record.Cost.PromptCost
		summary.AverageCompletionCost += record.Cost.CompletionCost
		summary.AverageCachedCost += record.Cost.CachedCost
		summary.AverageReasoningCost += record.Cost.ReasoningCost
	}

	if summary.RequestCount == 0 {
		return summary
	}

	requests := float64(summary.RequestCount)
	summary.AverageCostPerRequest = summary.TotalCost / requests
	summary.AveragePromptTokens = float64(summary.TotalPromptTokens) / requests
	summary.AverageCompletionTokens = float64(summary.TotalCompletionTokens) / requests
	summary.AverageCachedTokens = float64(summary.TotalCachedTokens) / requests
	summary.AverageReasoningTokens = float64(summary.TotalReasoningTokens) / requests
	summary.AverageTokensPerRequest = float64(summary.TotalTokens) / requests
	summary.AveragePromptCost /= requests
	summary.AverageCompletionCost /= requests
	summary.AverageCachedCost /= requests
	summary.AverageReasoningCost /= requests

	return summary
}

// ResetCostTracking clears all accumulated cost state.
func ResetCostTracking() {
	costMutex.Lock()
	defer costMutex.Unlock()

	totalCosts = nil
	costHistory = nil
	budgetLimits = nil
	budgetCallback = nil
}

func cloneCostRecord(record CostRecord) CostRecord {
	cloned := record
	if len(record.Tags) > 0 {
		cloned.Tags = make(map[string]string, len(record.Tags))
		for key, value := range record.Tags {
			cloned.Tags[key] = value
		}
	}
	return cloned
}
