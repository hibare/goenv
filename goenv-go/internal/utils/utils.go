package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-nv/goenv/internal/constants"
	"github.com/go-nv/goenv/internal/version"
)

func InitDirs() error {
	goenvRootDir, err := GetGoenvRootDir()
	if err != nil {
		return fmt.Errorf("failed to get goenv root directory: %w", err)
	}

	dirs := []string{
		constants.VersionsDir,
		constants.GlobalGoVersionFile,
		constants.LocalGoVersionFile,
	}

	// Ensure goenvRootDir exists, else create it
	if _, err := os.Stat(goenvRootDir); os.IsNotExist(err) {
		if err := os.MkdirAll(goenvRootDir, 0755); err != nil {
			return fmt.Errorf("failed to create goenv root directory: %w", err)
		}
	}

	// Ensure all dirs exist, else create them
	for _, dir := range dirs {
		p := filepath.Join(goenvRootDir, dir)
		if err := os.MkdirAll(p, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

func GetGoenvRootDir() (string, error) {
	// ToDo: @hibare handle GO_ENV_ROOT env variable
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	rootDir := filepath.Join(homeDir, constants.GoenvRootDir)
	return rootDir, nil
}

func GetHTTPUserAgent() string {
	return fmt.Sprintf("%s/%s", constants.ProjectName, version.CurrentVersion)
}

func GetOS() string {
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}

// ExtractArchive extracts a tar.gz archive to the specified directory.
func ExtractArchive(archivePath, targetDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	// Create the target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Extract all files from the archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// The Go distribution has a top-level "go" directory
		// We need to strip this prefix and place files directly in the target directory
		name := header.Name
		if strings.HasPrefix(name, "go/") {
			name = strings.TrimPrefix(name, "go/")
		} else if name == "go" {
			// Skip the top-level directory entry
			continue
		}

		target := filepath.Join(targetDir, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			// Create parent directories if they don't exist
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", target, err)
			}

			// Create the file
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}

			// Copy the file contents
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
			f.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink %s -> %s: %w", target, header.Linkname, err)
			}
		default:
			// Skip other types (e.g., hardlinks, character devices, etc.)
			continue
		}
	}

	return nil
}
