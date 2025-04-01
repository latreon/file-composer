package archiver

import (
	"io"
)

// ProgressCallback is a function type that receives progress updates
type ProgressCallback func(bytesWritten, totalSize int64)

// ProgressTracker implements functionality for tracking progress
type ProgressTracker struct {
	callback  ProgressCallback
	totalSize int64
	written   int64
}

// NewProgressCallback creates a new ProgressTracker with the provided callback
func NewProgressCallback(callback ProgressCallback) *ProgressTracker {
	return &ProgressTracker{
		callback: callback,
	}
}

// SetTotalSize sets the total size for progress calculation
func (pt *ProgressTracker) SetTotalSize(size int64) {
	pt.totalSize = size
	// Report initial progress
	pt.ReportProgress(0)
}

// AddProgress adds to the progress counter and reports the current progress
func (pt *ProgressTracker) AddProgress(bytes int64) {
	pt.written += bytes
	pt.ReportProgress(pt.written)
}

// ReportProgress reports the current progress through the callback
func (pt *ProgressTracker) ReportProgress(bytesWritten int64) {
	if pt.callback != nil {
		pt.callback(bytesWritten, pt.totalSize)
	}
}

// SetComplete marks the progress as complete
func (pt *ProgressTracker) SetComplete() {
	if pt.callback != nil {
		pt.callback(pt.totalSize, pt.totalSize)
	}
}

// ProgressWriter is an io.Writer that reports progress
type ProgressWriter struct {
	Writer   io.Writer
	Tracker  *ProgressTracker
	Progress int64
}

// NewProgressWriter creates a new progress-tracking writer
func NewProgressWriter(writer io.Writer, tracker *ProgressTracker) *ProgressWriter {
	return &ProgressWriter{
		Writer:  writer,
		Tracker: tracker,
	}
}

// Write implements the io.Writer interface with progress reporting
func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)
	if n > 0 {
		pw.Progress += int64(n)
		pw.Tracker.AddProgress(int64(n))
	}
	return
}
