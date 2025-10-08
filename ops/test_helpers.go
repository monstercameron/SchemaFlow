package ops

import (
	"context"

	"github.com/monstercameron/SchemaFlow/core"
)

type Person struct {
	Name string
	Age  int
}

func mockLLMResponse(ctx context.Context, system, user string, opts core.OpOptions) (string, error) {
	return `{"mock": "response"}`, nil
}

func setupMockClient() {
	core.SetLLMCaller(mockLLMResponse)
}
