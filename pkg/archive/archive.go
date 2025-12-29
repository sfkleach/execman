package archive

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractBinary extracts a binary from a tar.gz archive
func ExtractBinary(archivePath, destPath string) error {
	// Open the archive
	// #nosec G304 -- Opening archive in temp directory
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	// Create gzip reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// Create tar reader
	tr := tar.NewReader(gzr)

	// Get the directory of the destination to create a root scope
	destDir := filepath.Dir(destPath)
	destName := filepath.Base(destPath)

	// Create a scoped root for the destination directory to prevent path traversal
	root, err := os.OpenRoot(destDir)
	if err != nil {
		return fmt.Errorf("failed to create root scope: %w", err)
	}
	defer root.Close()

	// Find and extract the first executable file
	var extracted bool
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Skip directories
		if header.Typeflag == tar.TypeDir {
			continue
		}

		// Check if file is executable
		if header.Mode&0111 != 0 {
			// Extract this file using scoped root
			destFile, err := root.OpenFile(destName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("failed to create destination file: %w", err)
			}
			defer destFile.Close()

			// #nosec G110 -- Decompression of trusted GitHub release assets
			if _, err := io.Copy(destFile, tr); err != nil {
				return fmt.Errorf("failed to extract file: %w", err)
			}

			extracted = true
			break
		}
	}

	if !extracted {
		return fmt.Errorf("no executable file found in archive")
	}

	return nil
}

// CalculateChecksum calculates the SHA256 checksum of a file
func CalculateChecksum(filePath string) (string, error) {
	// #nosec G304 -- Calculating checksum of controlled file path
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}

// VerifyChecksum verifies a file's checksum against an expected value
func VerifyChecksum(filePath, expectedChecksum string) error {
	actualChecksum, err := CalculateChecksum(filePath)
	if err != nil {
		return err
	}

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// FindChecksumInFile finds the checksum for a specific file in a checksums.txt file
func FindChecksumInFile(checksumsPath, targetFilename string) (string, error) {
	// #nosec G304 -- Reading checksums from temp directory
	data, err := os.ReadFile(checksumsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read checksums file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format can be: "checksum  filename" or "checksum filename"
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			checksum := parts[0]
			filename := parts[len(parts)-1]

			// Handle paths in filename (extract basename)
			filename = filepath.Base(filename)

			if filename == targetFilename {
				return "sha256:" + checksum, nil
			}
		}
	}

	return "", fmt.Errorf("checksum not found for %s in checksums file", targetFilename)
}
