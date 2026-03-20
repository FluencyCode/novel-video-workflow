package providers

import (
	"path/filepath"
	"testing"

	configpkg "novel-video-workflow/pkg/config"
)

func newTestWindowsSubtitleProvider(t *testing.T) WindowsSubtitleProvider {
	t.Helper()
	return WindowsSubtitleProvider{
		baseDir: t.TempDir(),
		config: configpkg.SubtitleConfig{
			Provider: "windows-aegisub",
			Aegisub: configpkg.SubtitleAegisubConfig{
				ExecutablePath: "C:/Program Files/Aegisub/aegisub64.exe",
				ScriptPath:     filepath.Join("pkg", "tools", "aegisub", "aegisub_subtitle_gen.sh"),
				UseAutomation:  true,
			},
		},
		generateFunc: func(audioPath, text, outputPath string) error {
			return writeWindowsSubtitleSRT(outputPath, text)
		},
	}
}

func TestWindowsSubtitleProvider_GeneratesSRTForMinimalSample(t *testing.T) {
	provider := newTestWindowsSubtitleProvider(t)
	result, err := provider.Generate(SubtitleRequest{
		ProjectID:     "demo",
		ChapterNumber: 1,
		AudioPath:     filepath.Join("testdata", "workflow", "reference.wav"),
		Text:          "第一章 测试文本",
	})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(result.SubtitlePath) != ".srt" {
		t.Fatalf("expected srt output, got %q", result.SubtitlePath)
	}
	if result.Format != "srt" {
		t.Fatalf("expected srt format, got %q", result.Format)
	}
}

func TestWindowsSubtitleProvider_UsesConfiguredChapterSubtitleDirectory(t *testing.T) {
	provider := newTestWindowsSubtitleProvider(t)
	result, err := provider.Generate(SubtitleRequest{
		ProjectID:     "demo",
		ChapterNumber: 2,
		AudioPath:     filepath.Join("testdata", "workflow", "reference.wav"),
		Text:          "第二章 测试文本",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(provider.baseDir, "projects", "demo", "chapters", "02", "subtitle", "chapter_02.srt")
	if result.SubtitlePath != want {
		t.Fatalf("expected %q, got %q", want, result.SubtitlePath)
	}
}

func TestWindowsSubtitleProvider_ReturnsCategorizedErrorWhenGenerationFails(t *testing.T) {
	provider := newTestWindowsSubtitleProvider(t)
	provider.generateFunc = func(audioPath, text, outputPath string) error {
		return NewProviderError(CategoryExecutionError, "boom", nil)
	}
	_, err := provider.Generate(SubtitleRequest{
		ProjectID:     "demo",
		ChapterNumber: 3,
		AudioPath:     filepath.Join("testdata", "workflow", "reference.wav"),
		Text:          "第三章 测试文本",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	providerErr, ok := err.(ProviderError)
	if !ok {
		t.Fatalf("expected ProviderError, got %T", err)
	}
	if providerErr.Category != CategoryExecutionError {
		t.Fatalf("expected execution error, got %q", providerErr.Category)
	}
}

func TestWindowsSubtitleProvider_HealthCheck_BlocksOnMissingScriptPath(t *testing.T) {
	provider := WindowsSubtitleProvider{
		baseDir: t.TempDir(),
		config: configpkg.SubtitleConfig{
			Provider: "windows-aegisub",
			Aegisub: configpkg.SubtitleAegisubConfig{
				ScriptPath: "",
			},
		},
	}
	result := provider.HealthCheck()
	if result.Severity != SeverityBlocking {
		t.Fatalf("expected blocking severity, got %q", result.Severity)
	}
}

func TestWindowsSubtitleProvider_HealthCheck_PassesWithScriptPath(t *testing.T) {
	provider := WindowsSubtitleProvider{
		baseDir: t.TempDir(),
		config: configpkg.SubtitleConfig{
			Provider: "windows-aegisub",
			Aegisub: configpkg.SubtitleAegisubConfig{
				ScriptPath: "./pkg/tools/aegisub/aegisub_subtitle_gen.sh",
			},
		},
	}
	result := provider.HealthCheck()
	if result.Severity != SeverityInfo {
		t.Fatalf("expected info severity, got %q", result.Severity)
	}
}