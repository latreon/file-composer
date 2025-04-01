package main

import (
	"fmt"
	"log"
	"os"

	"github.com/latreon/file-compressor/pkg/archiver"
)

func main() {
	// Initialize logger
	log.SetPrefix("File Compressor: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// For now, we'll implement a simple CLI interface
	// Later, we can replace this with a proper GUI
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "compress":
		if len(os.Args) < 4 {
			fmt.Println("Insufficient arguments for compression")
			printUsage()
			return
		}
		sourcePath := os.Args[2]
		destPath := os.Args[3]

		// Default to ZIP format if not specified
		format := "zip"
		if len(os.Args) > 4 {
			format = os.Args[4]
		}

		err := archiver.Compress(sourcePath, destPath, format)
		if err != nil {
			log.Fatalf("Compression failed: %v", err)
		}
		fmt.Println("Compression completed successfully")

	case "extract":
		if len(os.Args) < 4 {
			fmt.Println("Insufficient arguments for extraction")
			printUsage()
			return
		}
		sourcePath := os.Args[2]
		destPath := os.Args[3]

		err := archiver.Extract(sourcePath, destPath)
		if err != nil {
			log.Fatalf("Extraction failed: %v", err)
		}
		fmt.Println("Extraction completed successfully")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  file-compressor compress <source> <destination> [format]")
	fmt.Println("  file-compressor extract <source> <destination>")
	fmt.Println()
	fmt.Println("Supported formats: zip, tar, gz, bz2, xz")
}
