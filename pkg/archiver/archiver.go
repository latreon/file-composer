package archiver

import (
	"archive/zip"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/latreon/file-compressor/pkg/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"golang.org/x/image/draw"
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
	case "pdf":
		if !strings.HasSuffix(strings.ToLower(sourcePath), ".pdf") {
			return fmt.Errorf("source file must be a PDF for PDF compression")
		}
		return compressPDF(sourcePath, destPath)
	case "png":
		if !strings.HasSuffix(strings.ToLower(sourcePath), ".png") {
			return fmt.Errorf("source file must be a PNG for PNG compression")
		}
		return compressPNG(sourcePath, destPath)
	case "jpg", "jpeg":
		if !strings.HasSuffix(strings.ToLower(sourcePath), ".jpg") && !strings.HasSuffix(strings.ToLower(sourcePath), ".jpeg") {
			return fmt.Errorf("source file must be a JPEG for JPEG compression")
		}
		return compressJPEG(sourcePath, destPath)
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
	case "pdf":
		if !strings.HasSuffix(strings.ToLower(sourcePath), ".pdf") {
			return fmt.Errorf("source file must be a PDF for PDF compression")
		}
		return compressPDF(sourcePath, destPath)
	case "png":
		if !strings.HasSuffix(strings.ToLower(sourcePath), ".png") {
			return fmt.Errorf("source file must be a PNG for PNG compression")
		}
		return compressPNG(sourcePath, destPath)
	case "jpg", "jpeg":
		if !strings.HasSuffix(strings.ToLower(sourcePath), ".jpg") && !strings.HasSuffix(strings.ToLower(sourcePath), ".jpeg") {
			return fmt.Errorf("source file must be a JPEG for JPEG compression")
		}
		return compressJPEG(sourcePath, destPath)
	case "zip":
		return compressZipWithProgress(sourcePath, destPath, info.IsDir(), progressTracker)
	// TODO: Implement other formats (tar, gz, bz2, xz, 7z)
	default:
		return fmt.Errorf("unsupported compression format: %s", format)
	}
}

// compressPNG compresses a PNG image with high compression
func compressPNG(sourcePath, destPath string) error {
	// Open the source file
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Decode the PNG image
	img, err := png.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Get original dimensions
	originalBounds := img.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()

	// Calculate new dimensions (reduce to 80% if large enough)
	var newWidth, newHeight int
	if originalWidth > 1000 || originalHeight > 1000 {
		// For large images, reduce more aggressively
		newWidth = originalWidth * 7 / 10
		newHeight = originalHeight * 7 / 10
	} else if originalWidth > 500 || originalHeight > 500 {
		// For medium images, reduce moderately
		newWidth = originalWidth * 8 / 10
		newHeight = originalHeight * 8 / 10
	} else {
		// For small images, don't resize
		newWidth = originalWidth
		newHeight = originalHeight
	}

	// Only resize if dimensions changed
	if newWidth != originalWidth || newHeight != originalHeight {
		// Create a new RGBA image
		dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		
		// Resize the image using high-quality resampling
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, originalBounds, draw.Over, nil)
		
		// Use the resized image for compression
		img = dst
	}

	// Create the destination file
	dstFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Create a PNG encoder with best compression
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	// Encode the image with best compression
	if err := encoder.Encode(dstFile, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// compressJPEG compresses a JPEG image with high compression
func compressJPEG(sourcePath, destPath string) error {
	// Open the source file
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Decode the JPEG image
	img, err := jpeg.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("failed to decode JPEG: %w", err)
	}

	// Get original dimensions
	originalBounds := img.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()

	// Calculate new dimensions (reduce to 80% if large enough)
	var newWidth, newHeight int
	if originalWidth > 1000 || originalHeight > 1000 {
		// For large images, reduce more aggressively
		newWidth = originalWidth * 7 / 10
		newHeight = originalHeight * 7 / 10
	} else if originalWidth > 500 || originalHeight > 500 {
		// For medium images, reduce moderately
		newWidth = originalWidth * 8 / 10
		newHeight = originalHeight * 8 / 10
	} else {
		// For small images, don't resize
		newWidth = originalWidth
		newHeight = originalHeight
	}

	// Only resize if dimensions changed
	if newWidth != originalWidth || newHeight != originalHeight {
		// Create a new RGBA image
		dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		
		// Resize the image using high-quality resampling
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, originalBounds, draw.Over, nil)
		
		// Use the resized image for compression
		img = dst
	}

	// Create the destination file
	dstFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Encode the image with maximum compression (quality 1 for absolute maximum compression)
	options := jpeg.Options{
		Quality: 1, // Lowest quality = highest compression
	}

	if err := jpeg.Encode(dstFile, img, &options); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

// compressPDF compresses a PDF file with extreme compression
func compressPDF(sourcePath, destPath string) error {
	// Create temporary files for multi-stage optimization
	tempFile1 := destPath + ".temp1"
	tempDir := destPath + ".tempdir"
	defer os.Remove(tempFile1) // Clean up when done
	defer os.RemoveAll(tempDir) // Clean up temp directory when done
	
	var err error
	
	// Try using Ghostscript for better compression
	cmd := exec.Command("gs", 
		"-sDEVICE=pdfwrite",
		"-dPDFSETTINGS=/screen", // Options: /screen (72dpi), /ebook (150dpi), /printer (300dpi), /prepress (300dpi+)
		"-dCompatibilityLevel=1.4",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-dColorImageDownsampleType=/Bicubic",
		"-dColorImageResolution=72",
		"-dGrayImageDownsampleType=/Bicubic",
		"-dGrayImageResolution=72",
		"-dMonoImageDownsampleType=/Bicubic", 
		"-dMonoImageResolution=72",
		"-sOutputFile="+destPath,
		sourcePath)
	
	// Try Ghostscript first
	err = cmd.Run()
	if err == nil {
		// Ghostscript succeeded
		return nil
	}
	
	// Ghostscript not available or failed, create temp directory for processing with pdfcpu
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Stage 1: Extract and optimize images
	// Extract images from PDF (if possible)
	err = api.ExtractImagesFile(sourcePath, tempDir, nil, nil)
	if err != nil {
		// If image extraction fails, just proceed with regular optimization
		fmt.Printf("Warning: Image extraction failed, proceeding with standard optimization: %v\n", err)
	} else {
		// Recompress all extracted images with high compression
		err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Skip directories
			if info.IsDir() {
				return nil
			}
			
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".jpg", ".jpeg":
				// Compress JPEG with quality 1
				return compressJPEG(path, path)
			case ".png":
				// Compress PNG with maximum compression
				return compressPNG(path, path)
			}
			
			return nil
		})
		
		if err != nil {
			fmt.Printf("Warning: Image recompression failed: %v\n", err)
		}
	}

	// Create a configuration with maximum compression settings
	conf := model.NewDefaultConfiguration()
	
	// Enable PDF 1.5 features for better compression
	conf.Reader15 = true
	conf.WriteObjectStream = true
	conf.WriteXRefStream = true
	
	// Stage 2: Apply optimization
	err = api.OptimizeFile(sourcePath, tempFile1, conf)
	if err != nil {
		return fmt.Errorf("failed PDF compression: %w", err)
	}
	
	// Stage 3: Convert to PDF 1.5 for better compression
	finalConf := model.NewDefaultConfiguration()
	finalConf.Reader15 = true
	finalConf.WriteObjectStream = true
	finalConf.WriteXRefStream = true
	
	// Apply final optimization
	err = api.OptimizeFile(tempFile1, destPath, finalConf)
	if err != nil {
		return fmt.Errorf("failed final PDF optimization: %w", err)
	}

	return nil
}

// compressZip compresses files or directories into a ZIP archive
func compressZip(sourcePath, destPath string, isDir bool) error {
	// Create the ZIP file
	zipFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file: %w", err)
	}
	defer zipFile.Close()

	// Create a new ZIP writer with best compression
	zipWriter := zip.NewWriter(zipFile)
	
	// Set the default compression level to the best possible compression
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return utils.NewMaxCompressionWriter(out), nil
	})
	
	defer zipWriter.Close()

	// Use a larger buffer for better compression
	buffer := make([]byte, 4*1024*1024) // 4MB buffer

	// Walk through the source path
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a new file in the ZIP archive
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		// Create file header with best compression
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name = relPath

		zipFile, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Open the source file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Copy the file contents to the ZIP archive using our buffer
		_, err = io.CopyBuffer(zipFile, srcFile, buffer)
		return err
	})

	return err
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
	
	// Set the default compression level to the best possible compression
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return utils.NewMaxCompressionWriter(out), nil
	})
	
	defer zipWriter.Close()

	// Use a larger buffer for better compression
	buffer := make([]byte, 4*1024*1024) // 4MB buffer

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

			return addFileToZipWithBuffer(zipWriter, path, relPath, buffer)
		})

		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}
	} else {
		// Single file compression
		// Use the filename as the zip header name
		filename := filepath.Base(sourcePath)
		err = addFileToZipWithBuffer(zipWriter, sourcePath, filename, buffer)
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

// addFileToZipWithBuffer adds a file to the zip archive using a buffer
func addFileToZipWithBuffer(zipWriter *zip.Writer, filePath, zipPath string, buffer []byte) error {
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

	// Copy contents with progress tracking
	_, err = io.CopyBuffer(writer, file, buffer)
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
