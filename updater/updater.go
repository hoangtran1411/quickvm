package updater

import (
	"archive/zip"
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
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
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
	defer resp.Body.Close()

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

// DownloadAndInstall downloads and installs the latest version
func (u *Updater) DownloadAndInstall(release *Release) error {
	// Determine the correct asset based on architecture
	assetName := u.getAssetName()
	
	var downloadURL string
	var assetSize int64
	
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, assetName) {
			downloadURL = asset.BrowserDownloadURL
			assetSize = asset.Size
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no suitable release asset found for your platform")
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "quickvm-update-*.exe")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Download with progress
	_, err = io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to save update: %w", err)
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
	_ = os.Remove(oldPath)
	
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
	if err := createCleanupScript(exePath, oldPath); err != nil {
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
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// createCleanupScript creates a script to clean up old version after update
func createCleanupScript(exePath, oldPath string) error {
	// Create a PowerShell script that will delete the old version after a delay
	scriptPath := oldPath + ".cleanup.ps1"
	
	scriptContent := fmt.Sprintf(`# QuickVM Update Cleanup Script
# This script will delete itself after cleaning up

Start-Sleep -Seconds 2

# Try to remove old version
$oldFile = "%s"
if (Test-Path $oldFile) {
    try {
        Remove-Item $oldFile -Force -ErrorAction Stop
        Write-Host "✅ Cleaned up old version" -ForegroundColor Green
    } catch {
        # Silently fail if file is still locked
    }
}

# Delete this cleanup script
$scriptFile = "%s"
Start-Sleep -Milliseconds 500
Remove-Item $scriptFile -Force -ErrorAction SilentlyContinue
`, oldPath, scriptPath)

	// Write script to file
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return err
	}

	// Execute cleanup script in background
	// Start the process detached (don't wait for it)
	go func() {
		// Small delay to ensure current process can exit first
		time.Sleep(100 * time.Millisecond)
		_ = executeCommand(scriptPath)
	}()

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
	return cmd.Start()
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
		return err
	}
	defer resp.Body.Close()

	// Save ZIP file
	zipPath := filepath.Join(destPath, "quickvm-update.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return err
	}

	// Extract ZIP
	return extractZip(zipPath, destPath)
}

// extractZip extracts a ZIP file
func extractZip(zipPath, destPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destPath, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
