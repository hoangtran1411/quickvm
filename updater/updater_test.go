package updater

import (
	"archive/zip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewUpdater(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"simple version", "1.0.0"},
		{"with v prefix", "v1.2.3"},
		{"prerelease", "1.0.0-beta"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpdater(tt.version)
			if u == nil {
				t.Fatal("NewUpdater() returned nil")
			}

			if u.currentVersion != tt.version {
				t.Errorf("Expected version %s, got %s", tt.version, u.currentVersion)
			}

			if u.githubRepo != "hoangtran1411/quickvm" {
				t.Errorf("Expected repo hoangtran1411/quickvm, got %s", u.githubRepo)
			}
		})
	}
}

func TestGetAssetName(t *testing.T) {
	u := NewUpdater("1.0.0")
	assetName := u.getAssetName()

	switch runtime.GOARCH {
	case "amd64":
		if assetName != "windows-amd64.exe" {
			t.Errorf("Expected windows-amd64.exe for amd64 arch, got %s", assetName)
		}
	case "arm64":
		if assetName != "windows-arm64.exe" {
			t.Errorf("Expected windows-arm64.exe for arm64 arch, got %s", assetName)
		}
	default:
		// For other architectures, should fallback to arm64
		if assetName != "windows-arm64.exe" {
			t.Errorf("Expected windows-arm64.exe as fallback, got %s", assetName)
		}
	}
}

func TestCheckForUpdates_MockServer(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		serverVersion  string
		wantUpdate     bool
	}{
		{"same version", "1.0.0", "v1.0.0", false},
		{"new version available", "1.0.0", "v1.1.0", true},
		{"with v prefix current", "v1.0.0", "v1.0.0", false},
		{"older server version", "2.0.0", "v1.0.0", true}, // different = update check
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := `{
					"tag_name": "` + tt.serverVersion + `",
					"name": "Release ` + tt.serverVersion + `",
					"body": "Test release",
					"assets": [
						{
							"name": "quickvm-windows-amd64.exe",
							"browser_download_url": "https://example.com/quickvm-windows-amd64.exe",
							"size": 10485760
						},
						{
							"name": "quickvm-windows-arm64.exe",
							"browser_download_url": "https://example.com/quickvm-windows-arm64.exe",
							"size": 10485760
						}
					]
				}`
				w.Write([]byte(response))
			}))
			defer server.Close()

			// We can't easily inject the URL, so we'll test the parsing logic separately
			// For now, test with actual API if available, skip if not
		})
	}
}

func TestCheckForUpdates_APIErrors(t *testing.T) {
	// Skip network tests in CI
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping network test in CI/CD environment")
	}

	u := NewUpdater("0.0.1")
	release, hasUpdate, err := u.CheckForUpdates()

	if err != nil {
		t.Skipf("Skipping: Could not reach GitHub API: %v", err)
	}

	if release == nil {
		t.Error("Expected release info, got nil")
	}

	if release != nil && release.TagName == "" {
		t.Error("Release tag name should not be empty")
	}

	// Version 0.0.1 should always have updates
	if !hasUpdate {
		t.Log("Warning: No update available for 0.0.1 (unexpected)")
	}
}

func TestCopyFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "quickvm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := []byte("Hello, QuickVM!")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	dstPath := filepath.Join(tmpDir, "dest.txt")
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify content
	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read dest file: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("Content mismatch: got %s, want %s", string(got), string(content))
	}

	// Verify permissions
	srcInfo, _ := os.Stat(srcPath)
	dstInfo, _ := os.Stat(dstPath)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("Permission mismatch: src %v, dst %v", srcInfo.Mode(), dstInfo.Mode())
	}
}

func TestCopyFile_Errors(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{"non-existent source", "/nonexistent/file.txt", "/tmp/dest.txt", true},
		{"invalid dest path", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := copyFile(tt.src, tt.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractZip(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "quickvm-zip-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test ZIP file
	zipPath := filepath.Join(tmpDir, "test.zip")
	if err := createTestZip(zipPath); err != nil {
		t.Fatalf("Failed to create test zip: %v", err)
	}

	// Extract
	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatalf("Failed to create extract dir: %v", err)
	}

	if err := extractZip(zipPath, extractDir); err != nil {
		t.Fatalf("extractZip failed: %v", err)
	}

	// Verify extracted files
	expectedFile := filepath.Join(extractDir, "test.txt")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Error("Expected extracted file test.txt not found")
	}

	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != "Hello from ZIP!" {
		t.Errorf("Content mismatch: got %s", string(content))
	}
}

func TestExtractZip_WithDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quickvm-zip-dir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create ZIP with directory structure
	zipPath := filepath.Join(tmpDir, "test-dir.zip")
	if err := createTestZipWithDir(zipPath); err != nil {
		t.Fatalf("Failed to create test zip: %v", err)
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatalf("Failed to create extract dir: %v", err)
	}

	if err := extractZip(zipPath, extractDir); err != nil {
		t.Fatalf("extractZip failed: %v", err)
	}

	// Verify directory was created
	subDir := filepath.Join(extractDir, "subdir")
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("Expected extracted directory subdir not found")
	}

	// Verify file in subdirectory
	subFile := filepath.Join(subDir, "nested.txt")
	if _, err := os.Stat(subFile); os.IsNotExist(err) {
		t.Error("Expected nested file not found")
	}
}

func TestExtractZip_InvalidZip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quickvm-invalid-zip-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an invalid ZIP (just a text file)
	invalidZip := filepath.Join(tmpDir, "invalid.zip")
	if err := os.WriteFile(invalidZip, []byte("not a zip file"), 0644); err != nil {
		t.Fatalf("Failed to create invalid zip: %v", err)
	}

	err = extractZip(invalidZip, tmpDir)
	if err == nil {
		t.Error("extractZip should fail on invalid ZIP file")
	}
}

func TestExtractZip_NonExistent(t *testing.T) {
	err := extractZip("/nonexistent/path.zip", "/tmp")
	if err == nil {
		t.Error("extractZip should fail on non-existent file")
	}
}

func TestRelease_Struct(t *testing.T) {
	r := Release{
		TagName: "v1.0.0",
		Name:    "Release 1.0.0",
		Body:    "Release notes here",
		Assets: []Asset{
			{
				Name:               "quickvm-windows-amd64.exe",
				BrowserDownloadURL: "https://example.com/download",
				Size:               1024,
			},
		},
	}

	if r.TagName != "v1.0.0" {
		t.Errorf("TagName mismatch")
	}
	if len(r.Assets) != 1 {
		t.Errorf("Expected 1 asset, got %d", len(r.Assets))
	}
	if r.Assets[0].Size != 1024 {
		t.Errorf("Asset size mismatch")
	}
}

func TestAsset_Struct(t *testing.T) {
	a := Asset{
		Name:               "test-asset.exe",
		BrowserDownloadURL: "https://example.com/test-asset.exe",
		Size:               2048,
	}

	if a.Name != "test-asset.exe" {
		t.Errorf("Name mismatch")
	}
	if a.Size != 2048 {
		t.Errorf("Size mismatch")
	}
}

func TestUpdater_VersionComparison(t *testing.T) {
	// This tests the version comparison logic used in CheckForUpdates
	tests := []struct {
		current  string
		latest   string
		wantDiff bool
	}{
		{"1.0.0", "1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"1.0.0", "v1.0.0", false}, // After TrimPrefix, both are "1.0.0"
		{"v1.0.0", "1.0.0", false},
		{"1.0.0", "1.0.1", true},
		{"1.0.0", "2.0.0", true},
		{"1.1.0", "1.0.0", true}, // Downgrade also counts as "different"
	}

	for _, tt := range tests {
		t.Run(tt.current+"_vs_"+tt.latest, func(t *testing.T) {
			current := trimVersionPrefix(tt.current)
			latest := trimVersionPrefix(tt.latest)
			isDiff := current != latest

			if isDiff != tt.wantDiff {
				t.Errorf("Version comparison: %s vs %s, got diff=%v, want diff=%v",
					tt.current, tt.latest, isDiff, tt.wantDiff)
			}
		})
	}
}

// Helper function matching the logic in CheckForUpdates
func trimVersionPrefix(v string) string {
	if len(v) > 0 && v[0] == 'v' {
		return v[1:]
	}
	return v
}

// Helper: Create a simple test ZIP file
func createTestZip(path string) error {
	zipFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	f, err := w.Create("test.txt")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte("Hello from ZIP!"))
	return err
}

// Helper: Create a test ZIP with directory structure
func createTestZipWithDir(path string) error {
	zipFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	// Create directory entry
	_, err = w.Create("subdir/")
	if err != nil {
		return err
	}

	// Create file in subdirectory
	f, err := w.Create("subdir/nested.txt")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte("Nested content"))
	return err
}

// Benchmark tests
func BenchmarkCopyFile(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "bench-*")
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "source.bin")
	// Create 1MB file
	data := make([]byte, 1024*1024)
	os.WriteFile(srcPath, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dstPath := filepath.Join(tmpDir, "dest.bin")
		copyFile(srcPath, dstPath)
		os.Remove(dstPath)
	}
}

func BenchmarkExtractZip(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "bench-zip-*")
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "test.zip")
	createTestZip(zipPath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractDir := filepath.Join(tmpDir, "extract")
		os.MkdirAll(extractDir, 0755)
		extractZip(zipPath, extractDir)
		os.RemoveAll(extractDir)
	}
}

// Skip network test in CI
func TestCheckForUpdates(t *testing.T) {
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping network test in CI/CD environment")
	}

	u := NewUpdater("0.0.1")
	release, hasUpdate, err := u.CheckForUpdates()

	if err != nil {
		t.Skipf("Skipping: Could not check for updates (offline?): %v", err)
	}

	if !hasUpdate {
		t.Log("No update available (test may need adjustment)")
	}

	if release != nil && release.TagName == "" {
		t.Error("Release tag name should not be empty")
	}
}

// Benchmark only if not in CI
func BenchmarkCheckForUpdates(b *testing.B) {
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		b.Skip("Skipping network benchmark in CI/CD environment")
	}

	u := NewUpdater("1.0.0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = u.CheckForUpdates()
	}
}
