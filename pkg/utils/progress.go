package utils

import (
	"fmt"
	"io"
	"time"
)

// ProgressWriter is an io.Writer wrapper that displays progress information
type ProgressWriter struct {
	Writer       io.Writer
	TotalSize    int64
	BytesWritten int64
	StartTime    time.Time
	LastUpdate   time.Time
	Description  string
}

// NewProgressWriter creates a new ProgressWriter
func NewProgressWriter(writer io.Writer, totalSize int64, description string) *ProgressWriter {
	return &ProgressWriter{
		Writer:      writer,
		TotalSize:   totalSize,
		StartTime:   time.Now(),
		LastUpdate:  time.Now(),
		Description: description,
	}
}

// Write implements the io.Writer interface and updates progress information
func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.Writer.Write(p)
	pw.BytesWritten += int64(n)

	// Update progress every 500ms to avoid too frequent updates
	if time.Since(pw.LastUpdate) > 500*time.Millisecond {
		pw.DisplayProgress()
		pw.LastUpdate = time.Now()
	}

	return n, err
}

// DisplayProgress prints the current progress to the console
func (pw *ProgressWriter) DisplayProgress() {
	// Calculate percentage
	var percentage float64
	if pw.TotalSize > 0 {
		percentage = float64(pw.BytesWritten) / float64(pw.TotalSize) * 100
	} else {
		percentage = 0
	}

	// Calculate speed
	elapsed := time.Since(pw.StartTime).Seconds()
	var bytesPerSecond int64
	if elapsed > 0 {
		bytesPerSecond = int64(float64(pw.BytesWritten) / elapsed)
	}

	// Calculate ETA
	var eta string
	if bytesPerSecond > 0 && pw.TotalSize > 0 {
		secondsLeft := float64(pw.TotalSize-pw.BytesWritten) / float64(bytesPerSecond)
		if secondsLeft > 0 {
			eta = formatDuration(secondsLeft)
		} else {
			eta = "0s"
		}
	} else {
		eta = "Unknown"
	}

	// Build progress bar (50 chars width)
	const barWidth = 50
	completedWidth := int(percentage / 100 * barWidth)
	progressBar := "["
	for i := 0; i < barWidth; i++ {
		if i < completedWidth {
			progressBar += "="
		} else if i == completedWidth {
			progressBar += ">"
		} else {
			progressBar += " "
		}
	}
	progressBar += "]"

	// Print progress information
	fmt.Printf("\r%s: %s %.1f%% %s/s ETA: %s",
		pw.Description,
		progressBar,
		percentage,
		formatBytes(bytesPerSecond),
		eta,
	)
}

// Complete finalizes the progress display
func (pw *ProgressWriter) Complete() {
	pw.BytesWritten = pw.TotalSize
	pw.DisplayProgress()
	fmt.Printf(" - Done!\n")
}

// Helper functions to format time and sizes
func formatDuration(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, secs)
	}
	return fmt.Sprintf("%ds", secs)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
