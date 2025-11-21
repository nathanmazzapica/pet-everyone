package handler

import (
	"testing"
)

func TestMediaTypeToExtension(t *testing.T) {
	tests := []struct {
		mimeType string
		expected string
	}{
		{"image/jpeg", ".jpeg"},
		{"image/png", ".png"},
	}

	for _, tt := range tests {
		t.Run(tt.mimeType, func(t *testing.T) {
			result := mediaTypeToExtension(tt.mimeType)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetAssetDiskPath(t *testing.T) {
	cfg := APIConfig{assetsRoot: "/tmp/assets"}

	tests := []struct {
		inputPath string
		shouldErr bool
	}{
		{"../../../etc/passwd", true},
		{"6q4U3nRManjLSElhew9ttYL4TDHj6mJCZwzn88qUVP0.jpeg", false},
	}

	for _, tt := range tests {
		t.Run(tt.inputPath, func(t *testing.T) {
			_, err := cfg.getAssetDiskPath(tt.inputPath)
			if (err != nil) != tt.shouldErr {
				t.Errorf("Expected error: %v, got: %v", tt.shouldErr, err)
			}
		})
	}
}
