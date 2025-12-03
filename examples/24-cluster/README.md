# Cluster Example

Demonstrates the `Cluster` operation for grouping similar items semantically.

## What is Cluster?

The `Cluster` operation groups items based on semantic similarity:
- Groups related content without predefined categories
- Auto-generates cluster names and descriptions
- Identifies outliers that don't fit any group
- Uses LLM understanding for semantic grouping

## Key Features

- Automatic cluster count detection
- Multiple naming strategies
- Outlier identification
- Similarity threshold control

## Usage

```go
opts := ops.NewClusterOptions().
    WithNumClusters(3).
    WithNamingStrategy("descriptive").
    WithIncludeOutliers(true)

result, err := ops.Cluster(items, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithNumClusters()` | Number of clusters (0 = auto) | 0 |
| `WithNamingStrategy()` | How to name clusters (descriptive, keyword, summary) | "descriptive" |
| `WithSimilarityThreshold()` | Minimum similarity for grouping (0-1) | 0.5 |
| `WithIncludeOutliers()` | Report items that don't fit | false |
| `WithBalanceClusters()` | Try to balance cluster sizes | false |

## Running

```bash
cd examples/24-cluster
go run main.go
```

## Example Output

```
Created 3 clusters:

  Cluster: Machine Learning & AI
  Description: Articles about ML, neural networks, and NLP
  Items (4):
    - Machine learning algorithms improve with more training data
    - Deep learning uses neural networks with multiple layers
    - Natural language processing enables text understanding
    - Transformer models have revolutionized NLP tasks

  Cluster: Web Development Frameworks
  Description: Frontend JavaScript frameworks and tools
  Items (4):
    - JavaScript frameworks like React are popular for web development
    - Vue.js provides reactive data binding for web applications
    - TypeScript adds static typing to JavaScript
    - Angular is a comprehensive framework for enterprise apps
```

## Use Cases

- **Support tickets**: Group by issue type
- **Documents**: Organize by topic
- **Products**: Categorize inventory
- **Feedback**: Identify common themes
