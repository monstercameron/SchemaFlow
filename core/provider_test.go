package schemaflow

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestProviders(t *testing.T) {
	t.Run("OpenAIProvider", func(t *testing.T) {
		provider, err := NewOpenAIProvider(ProviderConfig{
			APIKey: "test-key",
		})
		
		if err != nil {
			t.Fatalf("Failed to create OpenAI provider: %v", err)
		}
		
		if provider.Name() != "openai" {
			t.Errorf("Expected provider name 'openai', got %s", provider.Name())
		}
		
		// Test cost estimation
		cost := provider.EstimateCost(CompletionRequest{
			SystemPrompt: "You are a helpful assistant",
			UserPrompt:   "Hello, how are you?",
			MaxTokens:    100,
		})
		
		if cost <= 0 {
			t.Error("Expected positive cost estimate")
		}
	})
	
	t.Run("AnthropicProvider", func(t *testing.T) {
		provider, err := NewAnthropicProvider(ProviderConfig{
			APIKey: "test-key",
		})
		
		if err != nil {
			t.Fatalf("Failed to create Anthropic provider: %v", err)
		}
		
		if provider.Name() != "anthropic" {
			t.Errorf("Expected provider name 'anthropic', got %s", provider.Name())
		}
		
		// Test mock completion
		resp, err := provider.Complete(context.Background(), CompletionRequest{
			SystemPrompt: "Test",
			UserPrompt:   "Hello",
		})
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if resp.Provider != "anthropic" {
			t.Errorf("Expected provider 'anthropic' in response, got %s", resp.Provider)
		}
	})
	
	t.Run("LocalProvider", func(t *testing.T) {
		provider, err := NewLocalProvider(ProviderConfig{})
		
		if err != nil {
			t.Fatalf("Failed to create local provider: %v", err)
		}
		
		if provider.Name() != "local" {
			t.Errorf("Expected provider name 'local', got %s", provider.Name())
		}
		
		// Test cost (should be 0)
		cost := provider.EstimateCost(CompletionRequest{})
		if cost != 0 {
			t.Errorf("Expected 0 cost for local provider, got %f", cost)
		}
		
		// Test mock completion
		resp, err := provider.Complete(context.Background(), CompletionRequest{
			SystemPrompt: "You are an extraction expert",
			UserPrompt:   "Extract data from text",
			ResponseFormat: "json",
		})
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.Contains(resp.Content, "{") {
			t.Error("Expected JSON response from local provider")
		}
	})
	
	t.Run("LocalProviderWithHandler", func(t *testing.T) {
		provider, _ := NewLocalProvider(ProviderConfig{})
		
		// Set custom handler
		provider.WithHandler(func(ctx context.Context, req CompletionRequest) (string, error) {
			return fmt.Sprintf("Custom response to: %s", req.UserPrompt), nil
		})
		
		resp, err := provider.Complete(context.Background(), CompletionRequest{
			UserPrompt: "test prompt",
		})
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if resp.Content != "Custom response to: test prompt" {
			t.Errorf("Expected custom response, got: %s", resp.Content)
		}
	})
}

func TestProviderRegistry(t *testing.T) {
	t.Run("RegisterAndGet", func(t *testing.T) {
		registry := NewProviderRegistry()
		
		// Create and register providers
		openai, _ := NewOpenAIProvider(ProviderConfig{APIKey: "test"})
		local, _ := NewLocalProvider(ProviderConfig{})
		
		err := registry.Register("openai", openai)
		if err != nil {
			t.Fatalf("Failed to register OpenAI provider: %v", err)
		}
		
		err = registry.Register("local", local)
		if err != nil {
			t.Fatalf("Failed to register local provider: %v", err)
		}
		
		// Get provider
		provider, err := registry.Get("openai")
		if err != nil {
			t.Fatalf("Failed to get OpenAI provider: %v", err)
		}
		
		if provider.Name() != "openai" {
			t.Errorf("Expected OpenAI provider, got %s", provider.Name())
		}
		
		// Get default (should be first registered)
		provider, err = registry.Get("")
		if err != nil {
			t.Fatalf("Failed to get default provider: %v", err)
		}
		
		if provider.Name() != "openai" {
			t.Errorf("Expected default to be OpenAI, got %s", provider.Name())
		}
	})
	
	t.Run("SetDefault", func(t *testing.T) {
		registry := NewProviderRegistry()
		
		openai, _ := NewOpenAIProvider(ProviderConfig{APIKey: "test"})
		local, _ := NewLocalProvider(ProviderConfig{})
		
		registry.Register("openai", openai)
		registry.Register("local", local)
		
		// Change default
		err := registry.SetDefault("local")
		if err != nil {
			t.Fatalf("Failed to set default: %v", err)
		}
		
		provider, _ := registry.Get("")
		if provider.Name() != "local" {
			t.Errorf("Expected default to be local, got %s", provider.Name())
		}
	})
	
	t.Run("ListProviders", func(t *testing.T) {
		registry := NewProviderRegistry()
		
		openai, _ := NewOpenAIProvider(ProviderConfig{APIKey: "test"})
		local, _ := NewLocalProvider(ProviderConfig{})
		
		registry.Register("openai", openai)
		registry.Register("local", local)
		
		providers := registry.List()
		if len(providers) != 2 {
			t.Errorf("Expected 2 providers, got %d", len(providers))
		}
		
		// Check both providers are in the list
		hasOpenAI := false
		hasLocal := false
		for _, name := range providers {
			if name == "openai" {
				hasOpenAI = true
			}
			if name == "local" {
				hasLocal = true
			}
		}
		
		if !hasOpenAI || !hasLocal {
			t.Error("Missing expected providers in list")
		}
	})
}

func TestClientWithProviders(t *testing.T) {
	t.Run("DefaultProvider", func(t *testing.T) {
		client := NewClient("test-key")
		
		if client.provider == nil {
			t.Error("Expected provider to be initialized")
		}
		
		if client.providerName != "openai" {
			t.Errorf("Expected default provider to be openai, got %s", client.providerName)
		}
	})
	
	t.Run("WithProvider", func(t *testing.T) {
		client := NewClient("").
			WithProvider("local")
		
		if client.provider == nil {
			t.Error("Expected provider to be set")
		}
		
		if client.providerName != "local" {
			t.Errorf("Expected provider to be local, got %s", client.providerName)
		}
		
		// Test switching providers
		client.WithProvider("anthropic")
		
		if client.providerName != "anthropic" {
			t.Errorf("Expected provider to be anthropic, got %s", client.providerName)
		}
	})
	
	t.Run("WithProviderInstance", func(t *testing.T) {
		client := NewClient("")
		
		customProvider, _ := NewLocalProvider(ProviderConfig{})
		customProvider.WithHandler(func(ctx context.Context, req CompletionRequest) (string, error) {
			return "custom response", nil
		})
		
		client.WithProviderInstance(customProvider)
		
		if client.provider != customProvider {
			t.Error("Expected custom provider to be set")
		}
		
		if client.providerName != "local" {
			t.Errorf("Expected provider name to be local, got %s", client.providerName)
		}
	})
}

func TestProviderIntegration(t *testing.T) {
	// Save original state
	oldClient := defaultClient
	oldCallLLM := callLLM
	defer func() {
		defaultClient = oldClient
		callLLM = oldCallLLM
	}()
	
	t.Run("ExtractWithProvider", func(t *testing.T) {
		t.Skip("Skipped due to import cycle - ops imports core, so core tests can't import ops")
		
		// Create client with local provider
		client := NewClient("").WithProvider("local")
		defaultClient = client
		
		// Reset callLLM to use default implementation
		callLLM = defaultCallLLM
		
		type TestData struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		
		// This should use the provider-based implementation
		// result, err := ops.Extract[TestData]("Extract name and age", ops.NewExtractOptions())
		
		// if err != nil {
		// 	t.Fatalf("Unexpected error: %v", err)
		// }
		
		// Local provider returns mock data
		// if result.Name == "" {
		// 	t.Error("Expected name to be extracted")
		// }
	})
	
	t.Run("MultipleProvidersInParallel", func(t *testing.T) {
		// Register providers globally
		openaiProvider, _ := NewOpenAIProvider(ProviderConfig{APIKey: "test"})
		localProvider, _ := NewLocalProvider(ProviderConfig{})
		anthropicProvider, _ := NewAnthropicProvider(ProviderConfig{APIKey: "test"})
		
		RegisterProvider("test-openai", openaiProvider)
		RegisterProvider("test-local", localProvider)
		RegisterProvider("test-anthropic", anthropicProvider)
		
		// Create multiple clients with different providers
		clients := []*Client{
			NewClient("").WithProvider("test-openai"),
			NewClient("").WithProvider("test-local"),
			NewClient("").WithProvider("test-anthropic"),
		}
		
		// Verify each client has the correct provider
		expectedNames := []string{"openai", "local", "anthropic"}
		for i, client := range clients {
			if client.provider == nil {
				t.Errorf("Client %d has no provider", i)
				continue
			}
			if client.provider.Name() != expectedNames[i] {
				t.Errorf("Client %d: expected provider %s, got %s", 
					i, expectedNames[i], client.provider.Name())
			}
		}
	})
}

func TestProviderCostEstimation(t *testing.T) {
	req := CompletionRequest{
		SystemPrompt: strings.Repeat("a", 1000), // ~250 tokens
		UserPrompt:   strings.Repeat("b", 1000), // ~250 tokens
		MaxTokens:    500,
	}
	
	t.Run("OpenAICost", func(t *testing.T) {
		provider, _ := NewOpenAIProvider(ProviderConfig{APIKey: "test"})
		cost := provider.EstimateCost(req)
		
		// Should be roughly (500 * 0.03 + 500 * 0.06) / 1000
		// = (15 + 30) / 1000 = 0.045
		if cost < 0.01 || cost > 0.1 {
			t.Errorf("Unexpected cost estimate: %f", cost)
		}
	})
	
	t.Run("AnthropicCost", func(t *testing.T) {
		provider, _ := NewAnthropicProvider(ProviderConfig{APIKey: "test"})
		cost := provider.EstimateCost(req)
		
		// Should be lower than OpenAI due to different pricing
		if cost < 0.001 || cost > 0.05 {
			t.Errorf("Unexpected cost estimate: %f", cost)
		}
	})
	
	t.Run("LocalCost", func(t *testing.T) {
		provider, _ := NewLocalProvider(ProviderConfig{})
		cost := provider.EstimateCost(req)
		
		if cost != 0 {
			t.Errorf("Expected 0 cost for local provider, got %f", cost)
		}
	})
}

func TestProviderTimeout(t *testing.T) {
	t.Run("TimeoutHandling", func(t *testing.T) {
		provider, _ := NewLocalProvider(ProviderConfig{})
		
		// Set handler that takes too long
		provider.WithHandler(func(ctx context.Context, req CompletionRequest) (string, error) {
			select {
			case <-time.After(1 * time.Second):
				return "too late", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		})
		
		// Create context with short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		
		_, err := provider.Complete(ctx, CompletionRequest{
			UserPrompt: "test",
		})
		
		if err == nil {
			t.Error("Expected timeout error")
		}
		
		if !strings.Contains(err.Error(), "context") {
			t.Errorf("Expected context error, got: %v", err)
		}
	})
}