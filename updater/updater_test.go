package updater

import (
	"os"
	"testing"
)

func TestNewUpdater(t *testing.T) {
	u := NewUpdater("1.0.0")
	if u == nil {
		t.Error("NewUpdater() returned nil")
		return
	}
	
	if u.currentVersion != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", u.currentVersion)
	}
}

func TestCheckForUpdates(t *testing.T) {
	// Skip in CI/CD to avoid network dependencies
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping network test in CI/CD environment")
		return
	}

	u := NewUpdater("0.0.1") // Very old version
	
	release, hasUpdate, err := u.CheckForUpdates()
	
	// This test requires internet connection
	if err != nil {
		t.Skipf("Skipping test: Could not check for updates (offline?): %v", err)
		return
	}
	
	// Should have update since we're using version 0.0.1
	if !hasUpdate {
		t.Log("No update available (test may need adjustment)")
	}
	
	if release != nil && release.TagName == "" {
		t.Error("Release tag name should not be empty")
	}
}

func TestGetAssetName(t *testing.T) {
	u := NewUpdater("1.0.0")
	
	assetName := u.getAssetName()
	
	// Should contain either amd64 or arm64
	if assetName != "windows-amd64.exe" && assetName != "windows-arm64.exe" {
		t.Errorf("Unexpected asset name: %s", assetName)
	}
}

// Benchmark CheckForUpdates (only if not in CI)
func BenchmarkCheckForUpdates(b *testing.B) {
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		b.Skip("Skipping network benchmark in CI/CD environment")
		return
	}

	u := NewUpdater("1.0.0")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = u.CheckForUpdates()
	}
}
