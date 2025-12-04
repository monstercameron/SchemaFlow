package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVisionToolStub(t *testing.T) {
	result, _ := VisionTool.Execute(context.Background(), map[string]any{
		"action": "describe",
		"image":  "test.jpg",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected vision to be stubbed")
	}
}

func TestOCRToolStub(t *testing.T) {
	result, _ := OCRTool.Execute(context.Background(), map[string]any{
		"image": "test.jpg",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected OCR to be stubbed")
	}
}

func TestImageInfoTool(t *testing.T) {
	// Create a minimal PNG file for testing
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.png")

	// Minimal 1x1 PNG (red pixel)
	pngData := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, // PNG signature
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x18, 0xdd,
		0x8d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
		0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}

	if err := os.WriteFile(imgPath, pngData, 0644); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	result, _ := ImageInfoTool.Execute(context.Background(), map[string]any{
		"path": imgPath,
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	data := result.Data.(map[string]any)
	if data["format"] != "png" {
		t.Errorf("Expected png, got %v", data["format"])
	}
	// Size should be > 0
	if data["size"].(int64) <= 0 {
		t.Errorf("Expected positive size, got %v", data["size"])
	}
}

func TestImageInfoToolMissing(t *testing.T) {
	result, _ := ImageInfoTool.Execute(context.Background(), map[string]any{
		"path": "/nonexistent/image.png",
	})
	if result.Success {
		t.Error("Expected failure for missing file")
	}
}

func TestImageResizeToolStub(t *testing.T) {
	result, _ := ImageResizeTool.Execute(context.Background(), map[string]any{
		"input":  "input.jpg",
		"output": "output.jpg",
		"width":  100,
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected image_resize to be stubbed")
	}
}

func TestImageCropToolStub(t *testing.T) {
	result, _ := ImageCropTool.Execute(context.Background(), map[string]any{
		"input":  "input.jpg",
		"output": "output.jpg",
		"x":      0,
		"y":      0,
		"width":  100,
		"height": 100,
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected image_crop to be stubbed")
	}
}

func TestImageConvertToolStub(t *testing.T) {
	result, _ := ImageConvertTool.Execute(context.Background(), map[string]any{
		"input":  "input.png",
		"output": "output.jpg",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected image_convert to be stubbed")
	}
}

func TestImageBase64ToolEncode(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.png")

	// Create test file
	testData := []byte("test image data")
	if err := os.WriteFile(imgPath, testData, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, _ := ImageBase64Tool.Execute(context.Background(), map[string]any{
		"action": "encode",
		"path":   imgPath,
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	data := result.Data.(map[string]any)
	encoded := data["base64"].(string)
	if encoded == "" {
		t.Error("Expected non-empty base64 string")
	}
}

func TestImageBase64ToolEncodeDataURL(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.png")

	testData := []byte("test")
	os.WriteFile(imgPath, testData, 0644)

	result, _ := ImageBase64Tool.Execute(context.Background(), map[string]any{
		"action": "encode",
		"path":   imgPath,
	})

	data := result.Data.(map[string]any)
	dataUri := data["data_uri"].(string)
	if len(dataUri) < 20 || !strings.HasPrefix(dataUri, "data:image") {
		t.Errorf("Expected data URL format, got: %s", dataUri)
	}
}

func TestImageBase64ToolDecode(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "decoded.bin")

	// Base64 of "Hello"
	testBase64 := "SGVsbG8="

	result, _ := ImageBase64Tool.Execute(context.Background(), map[string]any{
		"action": "decode",
		"path":   outPath,
		"data":   testBase64,
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Failed to read decoded file: %v", err)
	}
	if string(content) != "Hello" {
		t.Errorf("Expected 'Hello', got %q", string(content))
	}
}

func TestImageBase64ToolDecodeDataURL(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "decoded.bin")

	// Data URL with base64 of "Test"
	testDataURL := "data:image/png;base64,VGVzdA=="

	result, _ := ImageBase64Tool.Execute(context.Background(), map[string]any{
		"action": "decode",
		"path":   outPath,
		"data":   testDataURL,
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	content, _ := os.ReadFile(outPath)
	if string(content) != "Test" {
		t.Errorf("Expected 'Test', got %q", string(content))
	}
}

func TestThumbnailToolStub(t *testing.T) {
	result, _ := ThumbnailTool.Execute(context.Background(), map[string]any{
		"input":  "input.jpg",
		"output": "thumb.jpg",
	})
	if result.Metadata["stubbed"] != true {
		t.Error("Expected thumbnail to be stubbed")
	}
}
