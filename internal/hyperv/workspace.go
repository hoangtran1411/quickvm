package hyperv

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Workspace represents a group of virtual machines
type Workspace struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	VMs         []string `yaml:"vms"` // List of VM names
}

// GetWorkspaceDir returns the directory where workspace files are stored
func GetWorkspaceDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	dir := filepath.Join(home, ".quickvm", "workspaces")
	// gosec G301: Expect directory permissions to be 0750 or less
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return "", fmt.Errorf("failed to create workspace directory: %w", err)
		}
	}

	return dir, nil
}

// SaveWorkspace saves a workspace to a YAML file
func SaveWorkspace(ws *Workspace) error {
	dir, err := GetWorkspaceDir()
	if err != nil {
		return err
	}

	filename := filepath.Join(dir, ws.Name+".yaml")
	data, err := yaml.Marshal(ws)
	if err != nil {
		return fmt.Errorf("failed to marshal workspace: %w", err)
	}

	// gosec G306: Expect WriteFile permissions to be 0600 or less
	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write workspace file: %w", err)
	}
	return nil
}

// LoadWorkspace loads a workspace by name
func LoadWorkspace(name string) (*Workspace, error) {
	dir, err := GetWorkspaceDir()
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(dir, name+".yaml")
	//nolint:gosec // G304: Path is constructed from trusted dir and name + literal extension.
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read workspace file '%s': %w", name, err)
	}

	var ws Workspace
	if err := yaml.Unmarshal(data, &ws); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workspace: %w", err)
	}

	return &ws, nil
}

// ListWorkspaces returns a list of all available workspace names
func ListWorkspaces() ([]string, error) {
	dir, err := GetWorkspaceDir()
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list workspace files: %w", err)
	}

	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, filepath.Base(file[:len(file)-5])) // Strip .yaml
	}

	return names, nil
}

// DeleteWorkspace deletes a workspace file
func DeleteWorkspace(name string) error {
	dir, err := GetWorkspaceDir()
	if err != nil {
		return err
	}

	filename := filepath.Join(dir, name+".yaml")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("workspace '%s' does not exist", name)
	}

	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete workspace file: %w", err)
	}
	return nil
}
