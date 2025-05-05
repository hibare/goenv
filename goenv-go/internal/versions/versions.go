package versions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-nv/goenv/internal/constants"
	"github.com/go-nv/goenv/internal/installer"
	"github.com/go-nv/goenv/internal/utils"
)

// VersionManager handles Go version management.
type VersionManager struct {
	rootDir                  string
	versionsDir              string
	globalVersionFile        string
	legacyDefaultVersionFile string
	legacyGlobalVersionFile  string
}

// NewVersionManager creates a new VersionManager instance.
func NewVersionManager() (*VersionManager, error) {
	rootDir, err := utils.GetGoenvRootDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	return &VersionManager{
		rootDir:                  rootDir,
		versionsDir:              filepath.Join(rootDir, constants.VersionsDir),
		globalVersionFile:        filepath.Join(rootDir, constants.GlobalGoVersionFile),
		legacyDefaultVersionFile: filepath.Join(rootDir, constants.GlobalVersionFileLegacyDefault),
		legacyGlobalVersionFile:  filepath.Join(rootDir, constants.GlobalVersionFileLegacyGlobal),
	}, nil
}

// RootDir returns the root directory of goenv.
func (vm *VersionManager) RootDir() string {
	return vm.rootDir
}

// GetVersionsDir returns the versions directory.
func (vm *VersionManager) GetVersionsDir() string {
	return vm.versionsDir
}

// GetGlobalVersionFile returns the global version file.
func (vm *VersionManager) GetGlobalVersionFile() string {
	return vm.globalVersionFile
}

// GetLegacyDefaultVersionFile returns the legacy default version file.
func (vm *VersionManager) GetLegacyDefaultVersionFile() string {
	return vm.legacyDefaultVersionFile
}

// GetLegacyGlobalVersionFile returns the legacy global version file.
func (vm *VersionManager) GetLegacyGlobalVersionFile() string {
	return vm.legacyGlobalVersionFile
}

// GetLocalVersionFile returns the local version file.
func (vm *VersionManager) GetLocalVersionFile() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return filepath.Join(currentDir, constants.LocalGoVersionFile), nil
}

// GetLocalVersion returns the local version.
func (vm *VersionManager) GetLocalVersion() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	versionFilePath, err := vm.GetLocalVersionFile()
	if err != nil {
		return "", fmt.Errorf("failed to get local version file: %w", err)
	}
	if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("No local version found in %s", currentDir)
	}

	return vm.readVersionFile(versionFilePath)
}

// SetLocalVersion sets the local version.
func (vm *VersionManager) SetLocalVersion(version string) error {
	versionFilePath, err := vm.GetLocalVersionFile()
	if err != nil {
		return fmt.Errorf("failed to get local version file: %w", err)
	}

	// if the version is not installed, install it
	if !vm.IsVersionInstalled(version) {
		installer, err := installer.NewInstaller()
		if err != nil {
			return fmt.Errorf("failed to create installer: %w", err)
		}
		if err := installer.Install(version); err != nil {
			return fmt.Errorf("failed to install version: %w", err)
		}
	}

	// write the version to the local version file
	if err := vm.writeVersionFile(versionFilePath, version); err != nil {
		return fmt.Errorf("failed to write version file: %w", err)
	}

	return nil
}

// UnsetLocalVersion unsets the local version.
func (vm *VersionManager) UnsetLocalVersion() error {
	versionFilePath, err := vm.GetLocalVersionFile()
	if err != nil {
		return fmt.Errorf("failed to get local version file: %w", err)
	}
	if err := os.Remove(versionFilePath); err != nil {
		return fmt.Errorf("failed to remove version file: %w", err)
	}

	return nil
}

// readVersionFile reads the version file.
func (vm *VersionManager) readVersionFile(versionFilePath string) (string, error) {
	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// writeVersionFile writes the version file.
func (vm *VersionManager) writeVersionFile(versionFilePath string, version string) error {
	if err := os.WriteFile(versionFilePath, []byte(version), 0644); err != nil {
		return err
	}

	return nil
}

// versionFileExists checks if the version file exists.
func (vm *VersionManager) versionFileExists(versionFilePath string) bool {
	_, err := os.Stat(versionFilePath)
	return err == nil
}

// rmVersionFile removes the version file.
func (vm *VersionManager) rmVersionFile(versionFilePath string) error {
	return os.Remove(versionFilePath)
}

// GetGlobalVersion returns the global version.
func (vm *VersionManager) GetGlobalVersion() (string, error) {
	versionFilePath := vm.GetGlobalVersionFile()
	legacyDefaultVersionFilePath := vm.GetLegacyGlobalVersionFile()
	legacyGlobalVersionFilePath := vm.GetLegacyDefaultVersionFile()

	if _, err := os.Stat(versionFilePath); err == nil {
		return vm.readVersionFile(versionFilePath)
	}

	if _, err := os.Stat(legacyDefaultVersionFilePath); err == nil {
		return vm.readVersionFile(legacyDefaultVersionFilePath)
	}

	if _, err := os.Stat(legacyGlobalVersionFilePath); err == nil {
		return vm.readVersionFile(legacyGlobalVersionFilePath)
	}

	return constants.GoSystemVersion, nil
}

// SetGlobalVersion sets the global version.
func (vm *VersionManager) SetGlobalVersion(version string) error {
	// if the version is the system version, remove the global version file
	if version == constants.GoSystemVersion {
		var currentGlobalVersion string

		if vm.versionFileExists(filepath.Join(vm.rootDir, constants.GlobalGoVersionFile)) {
			var err error
			currentGlobalVersion, err = vm.GetGlobalVersion()
			if err != nil {
				return fmt.Errorf("failed to get current global version: %w", err)
			}
		}

		if err := vm.rmVersionFile(filepath.Join(vm.rootDir, constants.GlobalGoVersionFile)); err != nil {
			return fmt.Errorf("failed to remove global version file: %w", err)
		}

		if currentGlobalVersion != "" {
			fmt.Printf("using system version instead of %s now\n", currentGlobalVersion)
		}
		return nil
	}

	// if the version is not installed, install it
	if !vm.IsVersionInstalled(version) {
		installer, err := installer.NewInstaller()
		if err != nil {
			return fmt.Errorf("failed to create installer: %w", err)
		}
		if err := installer.Install(version); err != nil {
			return fmt.Errorf("failed to install version: %w", err)
		}
	}

	// write the version to the global version file
	versionFilePath := filepath.Join(vm.rootDir, constants.GlobalGoVersionFile)
	if err := os.WriteFile(versionFilePath, []byte(version), 0644); err != nil {
		return fmt.Errorf("failed to write version file: %w", err)
	}

	return nil
}

// ListVersions lists all versions.
func (vm *VersionManager) ListVersions() ([]string, error) {
	versionsDir := filepath.Join(vm.rootDir, constants.VersionsDir)
	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("No versions found in %s", versionsDir)
	}

	files, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions directory: %w", err)
	}

	versions := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			versions = append(versions, file.Name())
		}
	}

	return versions, nil
}

// IsVersionInstalled checks if a version is installed.
func (vm *VersionManager) IsVersionInstalled(version string) bool {
	versionsDir := filepath.Join(vm.GetVersionsDir(), version)
	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		return false
	}

	return true
}
