# Predict Example

Demonstrates the `Predict` operation for forecasting and extrapolation.

## What is Predict?

The `Predict` operation forecasts future values based on historical data:
- Identifies trends and patterns
- Generates confidence intervals
- Creates alternative scenarios
- Explains prediction reasoning

## Key Features

- Multiple prediction methods
- Confidence intervals
- Scenario generation
- Factor analysis
- Risk identification

## Usage

```go
opts := ops.NewPredictOptions().
    WithHorizon("next_quarter").
    WithIncludeConfidenceInterval(true).
    WithIncludeScenarios(true)

result, err := ops.Predict[Forecast](historicalData, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithHorizon()` | Prediction timeframe | required |
| `WithMethod()` | Method (trend, pattern, regression, auto) | "auto" |
| `WithConfidenceLevel()` | Confidence for intervals | 0.8 |
| `WithIncludeScenarios()` | Generate alternative scenarios | false |
| `WithNumScenarios()` | Number of scenarios | 3 |
| `WithFactors()` | Factors to consider | [] |
| `WithAssumptions()` | Assumptions to use | [] |
| `WithIncludeReasoning()` | Explain prediction logic | true |

## Running

```bash
cd examples/33-predict
go run main.go
```

## Example Output

```
Q3 2024 Sales Forecast:
  Predicted Revenue: $1,550,000.00
  Predicted Units: 7,800
  Growth Rate: 6.9%
  Confidence: 78%

  80% Confidence Interval:
    Revenue: $1,450,000.00 - $1,650,000.00

  Reasoning: Based on observed 15% YoY growth trend and Q3
  seasonality patterns, with adjustment for recent market...

Alternative Scenarios:

  Optimistic (25% probability)
  Description: Strong market conditions and successful launch
  Conditions: [new product launch, competitor issues]

  Pessimistic (20% probability)
  Description: Market downturn affects demand
  Conditions: [economic slowdown, increased competition]
```

## Use Cases

- **Sales forecasting**: Predict revenue and units
- **Resource planning**: Forecast staffing needs
- **Inventory management**: Predict demand
- **Financial projections**: Estimate budgets
- **Trend analysis**: Identify patterns
