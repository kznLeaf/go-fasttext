package gofasttext

import (
	"os"
	"path/filepath"
	"testing"
)

func modelPath(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("TEST_LID_MODEL"); p != "" {
		if _, err := os.Stat(p); err != nil {
			t.Skipf("TEST_LID_MODEL not found: %s", p)
		}
		return p
	}
	candidates := []string{
		filepath.Join("models", "lid.176.ftz"),
		filepath.Join("models", "lid.176.bin"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	t.Skip("LID model not found; set TEST_LID_MODEL or place models/lid.176.ftz")
	return ""
}

func TestPredictLanguages(t *testing.T) {
	path := modelPath(t)
	if err := LoadModel(path); err != nil {
		t.Fatalf("LoadModel: %v", err)
	}
	t.Cleanup(Close)

	cases := []struct {
		text string
		lang string
	}{
		{"Hello, how are you doing today?", "en"},
		{"你好世界，今天天气怎么样？", "zh"},
		{"こんにちは、世界！", "ja"},
	}

	for _, tc := range cases {
		lang, conf, err := Predict(tc.text)
		if err != nil {
			t.Fatalf("Predict(%q): %v", tc.text, err)
		}
		if lang != tc.lang {
			t.Errorf("Predict(%q) lang = %q, want %q", tc.text, lang, tc.lang)
		}
		if conf <= 0 || conf > 1 {
			t.Errorf("Predict(%q) confidence = %v, want in (0,1]", tc.text, conf)
		}
	}
}

func TestPredictRequiresModel(t *testing.T) {
	Close()
	_, _, err := Predict("hello")
	if err == nil {
		t.Fatal("expected error when model is not loaded")
	}
}

func TestPredictEmptyText(t *testing.T) {
	_, _, err := Predict("")
	if err == nil {
		t.Fatal("expected error for empty text")
	}
}

func TestLoadModelEmptyPath(t *testing.T) {
	if err := LoadModel(""); err == nil {
		t.Fatal("expected error for empty path")
	}
}
