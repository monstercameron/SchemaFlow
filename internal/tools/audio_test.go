package tools

import (
	"context"
	"testing"
)

func TestTextToSpeechToolStub(t *testing.T) {
	result, _ := TextToSpeechTool.Execute(context.Background(), map[string]any{
		"text": "Hello world",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected TTS to be stubbed")
	}
}

func TestSpeechToTextToolStub(t *testing.T) {
	result, _ := SpeechToTextTool.Execute(context.Background(), map[string]any{
		"audio": "speech.mp3",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected STT to be stubbed")
	}
}

func TestAudioInfoToolStub(t *testing.T) {
	result, _ := AudioInfoTool.Execute(context.Background(), map[string]any{
		"path": "audio.mp3",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected audio_info to be stubbed")
	}
}

func TestAudioConvertToolStub(t *testing.T) {
	result, _ := AudioConvertTool.Execute(context.Background(), map[string]any{
		"input":  "input.mp3",
		"output": "output.wav",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected audio_convert to be stubbed")
	}
}

func TestAudioTrimToolStub(t *testing.T) {
	result, _ := AudioTrimTool.Execute(context.Background(), map[string]any{
		"input":  "input.mp3",
		"output": "output.mp3",
		"start":  "0:30",
		"end":    "1:00",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected audio_trim to be stubbed")
	}
}

func TestAudioAnalyzeToolStub(t *testing.T) {
	result, _ := AudioAnalyzeTool.Execute(context.Background(), map[string]any{
		"path":     "audio.mp3",
		"features": []any{"volume", "silence"},
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected audio_analyze to be stubbed")
	}
}
