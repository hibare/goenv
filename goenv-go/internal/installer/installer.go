package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-nv/goenv/internal/constants"
	"github.com/go-nv/goenv/internal/utils"
)

// Installer handles Go version installation.
type Installer struct {
	rootDir string
}

// NewInstaller creates a new Installer instance.
func NewInstaller() (*Installer, error) {
	rootDir, err := utils.GetGoenvRootDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get goenv root directory: %w", err)
	}

	return &Installer{
		rootDir: rootDir,
	}, nil
}

func (i *Installer) makeGetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", utils.GetHTTPUserAgent())

	return http.DefaultClient.Do(req)
}

// func (i *Installer) resolveVersion(version string, versions []string) (string, error) {}

// Install downloads and installs a Go version.
func (i *Installer) Install(version string) error {
	systemOs := utils.GetOS()
	systemArch := utils.GetArch()
	versionDir := filepath.Join(i.rootDir, constants.VersionsDir)

	fmt.Printf("Installing Go %s for %s/%s\n", version, systemOs, systemArch)

	// Download the Go version
	url := fmt.Sprintf("https://golang.org/dl/go%s.%s-%s.%s", version, systemOs, systemArch, constants.ArchiveFormat)
	resp, err := i.makeGetRequest(url)
	if err != nil {
		return fmt.Errorf("failed to download Go version: %w", err)
	}
	defer resp.Body.Close()

	// Create a temporary file for the download
	tmpFile, err := os.CreateTemp("", "go-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Save the downloaded file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	// Extract the archive
	if err := utils.ExtractArchive(tmpFile.Name(), filepath.Join(versionDir, version)); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	return nil
}

// Uninstall removes a Go version.
func (i *Installer) Uninstall(version string) error {
	versionDir := filepath.Join(i.rootDir, constants.VersionsDir, version)
	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("failed to remove version directory: %w", err)
	}
	return nil
}

// ListAvailableVersions returns a list of available Go versions.
func (i *Installer) ListAvailableVersions() ([]string, error) {
	resp, err := i.makeGetRequest(constants.GoDevDl)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var _versions GoVersions
	if err := json.Unmarshal(body, &_versions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal versions: %w", err)
	}

	versions := make([]string, len(_versions))
	for i, version := range _versions {
		versions[i] = version.Version
	}

	return versions, nil
}
