package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/latreon/file-compressor/pkg/archiver"
	"github.com/rs/cors"
)

// Configuration
const (
	uploadDir       = "./uploads"
	compressedDir   = "./compressed"
	maxUploadSize   = 1024 * 1024 * 100 // 100MB
	serverPort      = 8080
	cleanupInterval = 1 * time.Hour // Cleanup uploaded files after 1 hour
)

type CompressResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	DownloadLink string `json:"downloadLink,omitempty"`
	OutputSize   int64  `json:"outputSize,omitempty"`
	InputSize    int64  `json:"inputSize,omitempty"`
}

// Progress structure for websocket updates
type ProgressUpdate struct {
	Percentage float64 `json:"percentage"`
	Stage      string  `json:"stage"`
	FileName   string  `json:"fileName"`
}

func main() {
	// Ensure upload and compressed directories exist
	ensureDirectories()

	// Start cleanup routine
	go cleanupRoutine()

	// Create router and add routes
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/compress", handleCompressFile).Methods("POST")
	r.HandleFunc("/api/formats", handleGetFormats).Methods("GET")
	r.HandleFunc("/download/{filename}", handleDownload).Methods("GET")

	// Serve the Next.js app later (if we want to serve it from the same Go server)
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web-ui/out")))

	// Apply CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	// Start server
	fmt.Printf("API server running at http://localhost:%d\n", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), handler))
}

// ensureDirectories creates necessary directories if they don't exist
func ensureDirectories() {
	log.Println("Creating necessary directories...")

	// Create upload directory
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Could not create upload directory: %v", err)
	}
	absUploadPath, err := filepath.Abs(uploadDir)
	if err != nil {
		log.Fatalf("Could not get absolute path for upload directory: %v", err)
	}
	log.Printf("Upload directory created at: %s", absUploadPath)

	// Create compressed directory
	if err := os.MkdirAll(compressedDir, 0755); err != nil {
		log.Fatalf("Could not create compressed directory: %v", err)
	}
	absCompressedPath, err := filepath.Abs(compressedDir)
	if err != nil {
		log.Fatalf("Could not get absolute path for compressed directory: %v", err)
	}
	log.Printf("Compressed directory created at: %s", absCompressedPath)

	// Verify write permissions
	testUploadPath := filepath.Join(uploadDir, "test_write_permissions.txt")
	testFile, err := os.Create(testUploadPath)
	if err != nil {
		log.Fatalf("Could not write to upload directory: %v", err)
	}
	testFile.Close()
	os.Remove(testUploadPath)

	testCompressedPath := filepath.Join(compressedDir, "test_write_permissions.txt")
	testFile, err = os.Create(testCompressedPath)
	if err != nil {
		log.Fatalf("Could not write to compressed directory: %v", err)
	}
	testFile.Close()
	os.Remove(testCompressedPath)

	log.Println("Directories created and write permissions verified successfully")
}

// cleanupRoutine periodically removes old files
func cleanupRoutine() {
	for {
		time.Sleep(cleanupInterval)
		cleanup()
	}
}

// cleanup removes files older than the cleanup interval
func cleanup() {
	removeOldFiles(uploadDir)
	removeOldFiles(compressedDir)
}

// removeOldFiles deletes files older than the cleanup interval
func removeOldFiles(dir string) {
	cutoff := time.Now().Add(-cleanupInterval)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				log.Printf("Error removing old file %s: %v", path, err)
			} else {
				log.Printf("Removed old file: %s", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking directory %s: %v", dir, err)
	}
}

// handleGetFormats returns the supported compression formats
func handleGetFormats(w http.ResponseWriter, r *http.Request) {
	formats := []string{"pdf", "zip", "png", "jpg", "jpeg"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"formats": formats,
	})
}

// handleCompressFile compresses an uploaded file
func handleCompressFile(w http.ResponseWriter, r *http.Request) {
	// Set JSON content type
	w.Header().Set("Content-Type", "application/json")

	// Enforce size limit
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("File too large or invalid form: %v", err))
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving the file: %v", err))
		return
	}
	defer file.Close()

	// Get compression format from form, if not specified, determine from file type
	format := r.FormValue("format")
	if format == "" {
		// If no format specified (auto-select), use the original file extension
		ext := strings.ToLower(filepath.Ext(handler.Filename))
		if ext != "" {
			format = ext[1:] // Remove the dot from extension
		} else {
			format = "zip" // Fallback to zip if no extension
		}
	}

	log.Printf("Using compression format: %s", format)

	// Check if the file is a PDF when PDF compression is selected
	if format == "pdf" && !strings.HasSuffix(strings.ToLower(handler.Filename), ".pdf") {
		respondWithError(w, http.StatusBadRequest, "PDF compression can only be used with PDF files")
		return
	}

	// Check if the file is an image when image compression is selected
	if (format == "png" || format == "jpg" || format == "jpeg") && 
	   !strings.HasSuffix(strings.ToLower(handler.Filename), "."+format) {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s compression can only be used with %s files", format, format))
		return
	}

	log.Printf("Received file: %s (%d bytes)", handler.Filename, handler.Size)

	// Generate secure uploaded filename (using timestamp to avoid collisions)
	timestamp := time.Now().UnixNano()
	uploadPath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", timestamp, handler.Filename))

	// Save the uploaded file
	outFile, err := os.Create(uploadPath)
	if err != nil {
		log.Printf("Error creating temporary file: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error saving the file: %v", err))
		return
	}
	defer outFile.Close()

	// Copy the file to our local file
	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Printf("Error saving uploaded file: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error saving the file: %v", err))
		return
	}

	// Get the original file size
	fileInfo, err := os.Stat(uploadPath)
	if err != nil {
		log.Printf("Error reading file info: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error reading file info: %v", err))
		return
	}
	inputSize := fileInfo.Size()

	// Generate output filename with "compressed" prefix
	baseName := strings.TrimSuffix(handler.Filename, filepath.Ext(handler.Filename))
	outputFilename := fmt.Sprintf("%d_%s_compressed.%s", timestamp, baseName, format)
	outputPath := filepath.Join(compressedDir, outputFilename)

	// Create progress tracker
	progressTracker := archiver.NewProgressCallback(func(bytesWritten, totalSize int64) {
		if totalSize > 0 {
			percentage := float64(bytesWritten) / float64(totalSize) * 100
			log.Printf("Compression progress: %.2f%% (%d/%d bytes)", percentage, bytesWritten, totalSize)
		}
	})

	// Compress the file
	log.Printf("Compressing file from %s to %s with format %s", uploadPath, outputPath, format)
	err = archiver.CompressWithProgress(uploadPath, outputPath, format, progressTracker)
	if err != nil {
		log.Printf("Error compressing file: %v", err)
		// Try to get more detailed error information
		if os.IsNotExist(err) {
			respondWithError(w, http.StatusInternalServerError, "Source file not found")
		} else if os.IsPermission(err) {
			respondWithError(w, http.StatusInternalServerError, "Permission denied while compressing file")
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error compressing file: %v", err))
		}
		return
	}

	// Get the compressed file size
	compressedInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Printf("Error reading compressed file info: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error reading compressed file info: %v", err))
		return
	}
	outputSize := compressedInfo.Size()

	log.Printf("Successfully compressed %s to %s. Original: %d bytes, Compressed: %d bytes",
		handler.Filename, outputPath, inputSize, outputSize)

	// Create download link
	downloadLink := fmt.Sprintf("/download/%s", outputFilename)

	// Return success response
	resp := CompressResponse{
		Success:      true,
		Message:      "File compressed successfully",
		DownloadLink: downloadLink,
		InputSize:    inputSize,
		OutputSize:   outputSize,
	}

	json.NewEncoder(w).Encode(resp)
}

// handleDownload serves a compressed file for download
func handleDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Validate filename to prevent directory traversal
	if filepath.Base(filename) != filename {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(compressedDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate headers for download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filePath)
}

// respondWithError sends an error response in JSON format
func respondWithError(w http.ResponseWriter, code int, message string) {
	resp := CompressResponse{
		Success: false,
		Message: message,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
