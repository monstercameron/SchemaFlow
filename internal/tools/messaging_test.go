package tools

import (
	"context"
	"testing"
)

func TestEmailToolStub(t *testing.T) {
	result, _ := EmailTool.Execute(context.Background(), map[string]any{
		"action":  "send",
		"to":      "test@example.com",
		"subject": "Test",
		"body":    "Hello",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected email to be stubbed")
	}
}

func TestSMSToolStub(t *testing.T) {
	result, _ := SMSTool.Execute(context.Background(), map[string]any{
		"action":  "send",
		"to":      "+15551234567",
		"message": "Hello",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected SMS to be stubbed")
	}
}

func TestPushToolStub(t *testing.T) {
	result, _ := PushTool.Execute(context.Background(), map[string]any{
		"action": "send",
		"token":  "device-token",
		"title":  "Test",
		"body":   "Hello",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected push to be stubbed")
	}
}

func TestSlackToolStub(t *testing.T) {
	result, _ := SlackTool.Execute(context.Background(), map[string]any{
		"action":  "message",
		"channel": "general",
		"text":    "Hello",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected slack to be stubbed")
	}
}

func TestDiscordToolStub(t *testing.T) {
	result, _ := DiscordTool.Execute(context.Background(), map[string]any{
		"action":     "message",
		"channel_id": "123456789",
		"content":    "Hello",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected discord to be stubbed")
	}
}

func TestWebhookNotifyToolStub(t *testing.T) {
	result, _ := WebhookNotifyTool.Execute(context.Background(), map[string]any{
		"url":     "https://example.com/webhook",
		"payload": map[string]any{"test": true},
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected webhook_notify to be stubbed")
	}
}
