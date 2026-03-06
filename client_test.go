package schemaflow

import (
	"context"
	"testing"
	"time"

	"github.com/monstercameron/schemaflow/internal/llm"
)

type stubProvider struct {
	name string
}

func (provider *stubProvider) Complete(context.Context, llm.CompletionRequest) (llm.CompletionResponse, error) {
	return llm.CompletionResponse{Provider: provider.name}, nil
}

func (provider *stubProvider) Name() string {
	return provider.name
}

func (provider *stubProvider) EstimateCost(llm.CompletionRequest) float64 {
	return 0
}

func (provider *stubProvider) RetryPolicy() (int, time.Duration) {
	return 0, 0
}

func TestWithProviderUsesVendorSpecificEnv(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "deepseek-env-key")

	client := NewClient("")
	client.WithProvider("deepseek")

	if client.provider == nil {
		t.Fatal("expected provider to be configured")
	}
	if client.provider.Name() != "deepseek" {
		t.Fatalf("expected deepseek provider, got %s", client.provider.Name())
	}
}

func TestWithProviderConfigUsesRegisteredFactory(t *testing.T) {
	const providerName = "custom-factory"

	err := RegisterProviderFactory(providerName, func(config ProviderConfig) (Provider, error) {
		return &stubProvider{name: providerName}, nil
	})
	if err != nil {
		t.Fatalf("failed to register provider factory: %v", err)
	}

	client := NewClient("")
	client.WithProviderConfig(providerName, ProviderConfig{})

	if client.provider == nil {
		t.Fatal("expected provider to be configured")
	}
	if client.provider.Name() != providerName {
		t.Fatalf("expected %s provider, got %s", providerName, client.provider.Name())
	}
}

func TestWithProviderInstance(t *testing.T) {
	client := NewClient("")
	client.WithProviderInstance(&stubProvider{name: "instance-provider"})

	if client.provider == nil {
		t.Fatal("expected provider instance to be set")
	}
	if client.provider.Name() != "instance-provider" {
		t.Fatalf("expected instance-provider, got %s", client.provider.Name())
	}
}
