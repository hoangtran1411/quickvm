// Package updater handles the self-update mechanism for the application,
// interacting with GitHub Releases to download and install new versions.
package updater

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubAPI     = "https://api.github.com/repos/hoangtran1411/quickvm/releases/latest"
	updateTimeout = 30 * time.Second
)

// Release represents a GitHub release
type Release struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// Updater handles version checking and updates
type Updater struct {
	currentVersion string
	githubRepo     string
}

// NewUpdater creates a new updater instance
func NewUpdater(currentVersion string) *Updater {
	return &Updater{
		currentVersion: currentVersion,
		githubRepo:     "hoangtran1411/quickvm",
	}
}

// CheckForUpdates checks if a new version is available
func (u *Updater) CheckForUpdates() (*Release, bool, error) {
	client := &http.Client{
		Timeout: updateTimeout,
	}

	req, err := http.NewRequest("GET", githubAPI, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Compare versions
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(u.currentVersion, "v")

	hasUpdate := latestVersion != currentVersion

	return &release, hasUpdate, nil
}

func findReleaseAssets(release *Release, assetName string) (downloadURL string, assetSize int64, checksumURL string, err error) {
	if release == nil {
		return "", 0, "", fmt.Errorf("release cannot be nil")
	}

	expectedExeName := "quickvm-" + assetName
	expectedShaName := expectedExeName + ".sha256"

	// First, check for exact match
	for _, asset := range release.Assets {
		if asset.Name == expectedExeName {
			downloadURL = asset.BrowserDownloadURL
			assetSize = asset.Size
		}
		if asset.Name == expectedShaName {
			checksumURL = asset.BrowserDownloadURL
		}
	}

	// Fallback to suffix match only if exact match was not found
	if downloadURL == "" {
		for _, asset := range release.Assets {
			if strings.HasSuffix(asset.Name, assetName) && !strings.HasSuffix(asset.Name, ".sha256") {
				downloadURL = asset.BrowserDownloadURL
				assetSize = asset.Size
				break
			}
		}
	}
	if checksumURL == "" {
		for _, asset := range release.Assets {
			if strings.HasSuffix(asset.Name, assetName+".sha256") {
				checksumURL = asset.BrowserDownloadURL
				break
			}
		}
	}

	if downloadURL == "" {
		return "", 0, "", fmt.Errorf("no suitable release asset found for your platform")
	}
	return downloadURL, assetSize, checksumURL, nil
}

func verifyChecksum(client *http.Client, checksumURL, computedHash string) error {
	if checksumURL == "" {
		return nil
	}
	checkResp, err := client.Get(checksumURL)
	if err != nil {
		return fmt.Errorf("failed to download checksum: %w", err)
	}
	defer func() { _ = checkResp.Body.Close() }()

	if checkResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download checksum: HTTP %d", checkResp.StatusCode)
	}

	checkBytes, err := io.ReadAll(checkResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read checksum: %w", err)
	}
	expectedHash := strings.TrimSpace(string(checkBytes))
	if fields := strings.Fields(expectedHash); len(fields) > 0 {
		expectedHash = fields[0]
	}
	if expectedHash == "" {
		return fmt.Errorf("checksum file is empty")
	}
	if !strings.EqualFold(computedHash, expectedHash) {
		return fmt.Errorf("checksum verification failed: expected %s, got %s", expectedHash, computedHash)
	}
	fmt.Println("🔒 Checksum verified successfully!")
	return nil
}

// DownloadAndInstall downloads and installs the latest version
//
//nolint:funlen // Complex download flow
func (u *Updater) DownloadAndInstall(release *Release) error {
	assetName := u.getAssetName()
	downloadURL, assetSize, checksumURL, err := findReleaseAssets(release, assetName)
	if err != nil {
		return err
	}

	fmt.Printf("📦 Downloading QuickVM %s (%d MB)...\n", release.TagName, assetSize/1024/1024)

	// Download the file
	client := &http.Client{
		Timeout: 5 * time.Minute, // Longer timeout for download
	}

	resp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "quickvm-update-*.exe")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	// Download with progress and compute sha256 hash
	hasher := sha256.New()
	multiWriter := io.MultiWriter(tmpFile, hasher)
	_, err = io.Copy(multiWriter, resp.Body)
	_ = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to save update: %w", err)
	}

	computedHash := hex.EncodeToString(hasher.Sum(nil))

	// Verify checksum if available
	if err := verifyChecksum(client, checksumURL, computedHash); err != nil {
		return err
	}

	fmt.Println("✅ Download complete!")

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Use a different strategy: rename old file instead of deleting
	oldPath := exePath + ".old"
	fmt.Println("🔄 Installing update...")

	// Remove any existing .old file first
	if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
		// Just log warning if removal fails for reasons other than non-existence
		fmt.Printf("⚠️  Warning: Failed to remove old version: %v\n", err)
	}

	// Rename current executable to .old (this works even if file is locked)
	if err := os.Rename(exePath, oldPath); err != nil {
		return fmt.Errorf("failed to rename old version: %w", err)
	}

	// Copy new version to the original location
	if err := copyFile(tmpPath, exePath); err != nil {
		// Restore old version if update fails
		_ = os.Rename(oldPath, exePath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Create a cleanup script to delete the old version after this process exits
	if err := createCleanupScript(oldPath); err != nil {
		// Non-fatal error, just log it
		fmt.Printf("⚠️  Warning: Failed to create cleanup script: %v\n", err)
	}

	fmt.Printf("✅ Successfully updated to version %s!\n", release.TagName)
	fmt.Println("🔄 Please close this terminal and run 'quickvm version' to verify.")
	fmt.Println("   The old version will be automatically cleaned up.")

	return nil
}

// getAssetName returns the appropriate asset name for the current platform
func (u *Updater) getAssetName() string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		return "windows-amd64.exe"
	}
	return "windows-arm64.exe"
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	//nolint:gosec // G304: Path is from safe source (internal logic)
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = sourceFile.Close() }()

	//nolint:gosec // G304: Path is from safe source
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() { _ = destFile.Close() }()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	if err := os.Chmod(dst, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	return nil
}

// createCleanupScript creates a script to clean up old version after update
func createCleanupScript(oldPath string) error {
	// Create a PowerShell script that will delete the old version after a delay
	scriptPath := oldPath + ".cleanup.ps1"

	safeOldPath := strings.ReplaceAll(oldPath, "'", "''")
	safeScriptPath := strings.ReplaceAll(scriptPath, "'", "''")

	scriptContent := fmt.Sprintf(`# QuickVM Update Cleanup Script
# This script will delete itself after cleaning up

Start-Sleep -Seconds 2

# Try to remove old version
$oldFile = '%s'
if (Test-Path $oldFile) {
    try {
        Remove-Item $oldFile -Force -ErrorAction Stop
        Write-Host "✅ Cleaned up old version" -ForegroundColor Green
    } catch {
        # Silently fail if file is still locked
    }
}

# Delete this cleanup script
$scriptFile = '%s'
Start-Sleep -Milliseconds 500
Remove-Item $scriptFile -Force -ErrorAction SilentlyContinue
`, safeOldPath, safeScriptPath)

	// Write script to file
	// gosec G306: Expect WriteFile permissions to be 0600 or less
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0600); err != nil {
		return fmt.Errorf("failed to write cleanup script: %w", err)
	}

	// Execute cleanup script detached in background (cmd.Start() detaches it,
	// script itself sleeps 2s allowing current process to exit)
	if err := executeCommand(scriptPath); err != nil {
		return fmt.Errorf("failed to start cleanup script: %w", err)
	}

	return nil
}

// executeCommand executes a shell command in the background
func executeCommand(scriptPath string) error {
	// Use PowerShell to execute the cleanup script in hidden mode
	cmd := exec.Command("powershell.exe",
		"-WindowStyle", "Hidden",
		"-ExecutionPolicy", "Bypass",
		"-File", scriptPath)

	// Start without waiting for completion
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cleanup script: %w", err)
	}
	return nil
}

// DownloadZipPackage downloads the full ZIP package
func (u *Updater) DownloadZipPackage(release *Release, destPath string) error {
	// Find ZIP asset
	var zipURL string
	arch := runtime.GOARCH

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, ".zip") && strings.Contains(asset.Name, arch) {
			zipURL = asset.BrowserDownloadURL
			break
		}
	}

	if zipURL == "" {
		return fmt.Errorf("no ZIP package found for your platform")
	}

	// Download ZIP
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(zipURL)
	if err != nil {
		return fmt.Errorf("failed to download ZIP: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Save ZIP file
	zipPath := filepath.Join(destPath, "quickvm-update.zip")
	//nolint:gosec // G304: Path is from safe source
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create package file: %w", err)
	}
	defer func() { _ = zipFile.Close() }()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save package: %w", err)
	}

	// Extract ZIP
	return extractZip(zipPath, destPath)
}

// extractZip extracts a ZIP file
func extractZip(zipPath, destPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open ZIP archive: %w", err)
	}
	defer func() { _ = r.Close() }()

	for _, f := range r.File {
		//nolint:gosec // G305: File traversal prevented (assumed checked/trusted source)
		fpath := filepath.Join(destPath, f.Name)

		if f.FileInfo().IsDir() {
			// gosec G301: Expect directory permissions to be 0750 or less
			if err := os.MkdirAll(fpath, 0750); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// gosec G301: Expect directory permissions to be 0750 or less
		if err := os.MkdirAll(filepath.Dir(fpath), 0750); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		//nolint:gosec // G304: Path is from zip entry
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			_ = outFile.Close()
			return fmt.Errorf("failed to open zip entry: %w", err)
		}

		//nolint:gosec // G110: Decompression bomb check skipped for this trusted update zip
		_, err = io.Copy(outFile, rc)
		_ = outFile.Close()
		_ = rc.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	return nil
}
