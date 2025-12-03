# Rank Example

Demonstrates the `Rank` operation for ordering items by query relevance.

## What is Rank?

The `Rank` operation orders items based on relevance to a query:
- Uses semantic understanding, not just keyword matching
- Supports field boosting and penalization
- Returns scores with explanations
- Configurable top-K limit

## Key Features

- Semantic relevance scoring
- Field-level boost/penalize
- Score explanations
- Top-K filtering

## Usage

```go
opts := ops.NewRankOptions().
    WithQuery("Go programming patterns").
    WithTopK(5).
    WithBoostFields([]string{"title", "skills"}).
    WithScoreDetails(true)

result, err := ops.Rank(items, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithQuery()` | The search/relevance query | required |
| `WithTopK()` | Maximum results to return | 0 (all) |
| `WithBoostFields()` | Fields to weight higher | [] |
| `WithPenalizeFields()` | Fields to weight lower | [] |
| `WithScoreDetails()` | Include score explanations | false |
| `WithNormalize()` | Normalize scores to 0-1 | true |

## Running

```bash
cd examples/25-rank
go run main.go
```

## Example Output

```
Query: 'Go programming patterns and best practices'
Top 3 results:

  1. Advanced Go Patterns (Score: 0.95)
     Reason: Directly discusses Go design patterns

  2. Go Concurrency Deep Dive (Score: 0.88)
     Reason: Covers Go-specific concurrency patterns

  3. Go Performance Optimization (Score: 0.82)
     Reason: Related to Go programming best practices
```

## Difference from Sort

- **Sort**: Orders by explicit criteria (e.g., "by date descending")
- **Rank**: Orders by semantic relevance to a query

## Use Cases

- **Search results**: Order by query relevance
- **Recommendations**: Rank by user preference match
- **Candidate matching**: Score job applicants
- **Content curation**: Prioritize relevant articles
