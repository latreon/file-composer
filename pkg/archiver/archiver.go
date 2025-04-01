package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/latreon/file-compressor/pkg/utils"
)

// Compress compresses files or directories at sourcePath to destPath using the specified format
func Compress(sourcePath, destPath, format string) error {
	// Check if source exists
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("source path error: %w", err)
	}

	// Validate and select compression format
	switch strings.ToLower(format) {
	case "zip":
		return compressZip(sourcePath, destPath, info.IsDir())
	// TODO: Implement other formats (tar, gz, bz2, xz, 7z)
	default:
		return fmt.Errorf("unsupported compression format: %s", format)
	}
}

// CompressWithProgress compresses files with progress reporting through a ProgressTracker
func CompressWithProgress(sourcePath, destPath, format string, progressTracker *ProgressTracker) error {
	// Check if source exists
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("source path error: %w", err)
	}

	// Calculate total size for progress reporting
	var totalSize int64
	if info.IsDir() {
		err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				totalSize += info.Size()
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error calculating total size: %w", err)
		}
	} else {
		totalSize = info.Size()
	}

	// Set total size in progress tracker
	progressTracker.SetTotalSize(totalSize)

	// Validate and select compression format
	switch strings.ToLower(format) {
	case "zip":
		return compressZipWithProgress(sourcePath, destPath, info.IsDir(), progressTracker)
	// TODO: Implement other formats (tar, gz, bz2, xz, 7z)
	default:
		return fmt.Errorf("unsupported compression format: %s", format)
	}
}

// compressZip compresses files using the ZIP format
func compressZip(sourcePath, destPath string, isDir bool) error {
	// First, calculate total size for progress reporting
	var totalSize int64
	if isDir {
		err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				totalSize += info.Size()
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error calculating total size: %w", err)
		}
	} else {
		info, err := os.Stat(sourcePath)
		if err != nil {
			return fmt.Errorf("error getting file info: %w", err)
		}
		totalSize = info.Size()
	}

	// Create the destination file
	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	// Create a progress writer
	progressWriter := utils.NewProgressWriter(dest, totalSize, "Compressing")
	defer progressWriter.Complete()

	// Create a new zip writer with progress tracking
	zipWriter := zip.NewWriter(progressWriter)
	defer zipWriter.Close()

	if isDir {
		// Walk through all files in the directory
		err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories themselves
			if info.IsDir() {
				return nil
			}

			// Create a relative path as zip header name
			relPath, err := filepath.Rel(sourcePath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			return addFileToZip(zipWriter, path, relPath)
		})

		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}
	} else {
		// Single file compression
		// Use the filename as the zip header name
		filename := filepath.Base(sourcePath)
		err = addFileToZip(zipWriter, sourcePath, filename)
		if err != nil {
			return fmt.Errorf("failed to add file to zip: %w", err)
		}
	}

	return nil
}

// compressZipWithProgress compresses files using the ZIP format with progress reporting
func compressZipWithProgress(sourcePath, destPath string, isDir bool, progressTracker *ProgressTracker) error {
	// Create the destination file
	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	// Create a progress writer
	progressWriter := NewProgressWriter(dest, progressTracker)

	// Create a new zip writer with progress tracking
	zipWriter := zip.NewWriter(progressWriter)
	defer zipWriter.Close()

	if isDir {
		// Walk through all files in the directory
		err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories themselves
			if info.IsDir() {
				return nil
			}

			// Create a relative path as zip header name
			relPath, err := filepath.Rel(sourcePath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			return addFileToZip(zipWriter, path, relPath)
		})

		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}
	} else {
		// Single file compression
		// Use the filename as the zip header name
		filename := filepath.Base(sourcePath)
		err = addFileToZip(zipWriter, sourcePath, filename)
		if err != nil {
			return fmt.Errorf("failed to add file to zip: %w", err)
		}
	}

	// Mark progress as complete
	progressTracker.SetComplete()

	return nil
}

// addFileToZip adds a file to the zip archive
func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Get file info for header
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Create zip header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("failed to create zip header: %w", err)
	}

	// Use the provided zip path for the header name
	header.Name = zipPath
	// Use best compression
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create zip header: %w", err)
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("failed to write file to zip: %w", err)
	}

	return nil
}

// Extract extracts the archive at sourcePath to destPath
func Extract(sourcePath, destPath string) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Determine archive type based on file extension
	ext := strings.ToLower(filepath.Ext(sourcePath))
	switch ext {
	case ".zip":
		return extractZip(sourcePath, destPath)
	// TODO: Implement other formats (tar, gz, bz2, xz, 7z)
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}
}

// extractZip extracts a ZIP archive
func extractZip(sourcePath, destPath string) error {
	// Open the zip file
	reader, err := zip.OpenReader(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Calculate total uncompressed size for progress tracking
	var totalSize int64
	for _, file := range reader.File {
		totalSize += int64(file.UncompressedSize64)
	}

	// Create a progress tracker (doesn't use a writer directly)
	progress := utils.NewProgressWriter(io.Discard, totalSize, "Extracting")
	defer progress.Complete()

	// Extract each file
	for _, file := range reader.File {
		err := extractZipFileWithProgress(file, destPath, progress)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractZipFileWithProgress extracts a single file from a ZIP archive with progress tracking
func extractZipFileWithProgress(file *zip.File, destPath string, progress *utils.ProgressWriter) error {
	// Prepare full path for extraction
	filePath := filepath.Join(destPath, file.Name)

	// Check for zip slip vulnerability (traversal attack)
	if !strings.HasPrefix(filePath, filepath.Clean(destPath)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", filePath)
	}

	// Create directory tree
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		return nil
	}

	// Create the directory tree for the file
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Create the file
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Open the file inside the archive
	inFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer inFile.Close()

	// Copy contents with progress tracking
	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		bytesRead, err := inFile.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from archive: %w", err)
		}

		if bytesRead > 0 {
			_, err := outFile.Write(buffer[:bytesRead])
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}

			// Update progress
			progress.BytesWritten += int64(bytesRead)
			progress.DisplayProgress()
		}

		if err == io.EOF {
			break
		}
	}

	// Restore file permissions
	if err := os.Chmod(filePath, file.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// ExtractWithProgress extracts an archive with progress reporting
func ExtractWithProgress(sourcePath, destPath string, progressTracker *ProgressTracker) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Determine archive type based on file extension
	ext := strings.ToLower(filepath.Ext(sourcePath))
	switch ext {
	case ".zip":
		return extractZipWithProgress(sourcePath, destPath, progressTracker)
	// TODO: Implement other formats (tar, gz, bz2, xz, 7z)
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}
}

// extractZipWithProgress extracts a ZIP archive with progress reporting
func extractZipWithProgress(sourcePath, destPath string, progressTracker *ProgressTracker) error {
	// Open the zip file
	reader, err := zip.OpenReader(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Calculate total uncompressed size for progress tracking
	var totalSize int64
	for _, file := range reader.File {
		totalSize += int64(file.UncompressedSize64)
	}

	// Set total size in progress tracker
	progressTracker.SetTotalSize(totalSize)

	// Extract each file
	for _, file := range reader.File {
		err := extractZipFileWithProgressTracker(file, destPath, progressTracker)
		if err != nil {
			return err
		}
	}

	// Mark progress as complete
	progressTracker.SetComplete()

	return nil
}

// extractZipFileWithProgressTracker extracts a single file from a ZIP archive with progress tracking
func extractZipFileWithProgressTracker(file *zip.File, destPath string, progressTracker *ProgressTracker) error {
	// Prepare full path for extraction
	filePath := filepath.Join(destPath, file.Name)

	// Check for zip slip vulnerability (traversal attack)
	if !strings.HasPrefix(filePath, filepath.Clean(destPath)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", filePath)
	}

	// Create directory tree
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		return nil
	}

	// Create the directory tree for the file
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Create the file
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Open the file inside the archive
	inFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer inFile.Close()

	// Copy contents with progress tracking
	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		bytesRead, err := inFile.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from archive: %w", err)
		}

		if bytesRead > 0 {
			_, err := outFile.Write(buffer[:bytesRead])
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}

			// Update progress
			progressTracker.AddProgress(int64(bytesRead))
		}

		if err == io.EOF {
			break
		}
	}

	// Restore file permissions
	if err := os.Chmod(filePath, file.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}
