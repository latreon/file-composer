# Advanced File Compressor

A high-performance, feature-rich file compression and extraction utility written in Go, with multiple user interfaces:

1. **Command-line interface** for scripting and power users
2. **Native GUI** built with Fyne for desktop integration
3. **Web-based UI** built with Next.js for browser access

## Features

Current features:
- Compress files and directories to ZIP format
- Extract files from ZIP archives
- Progress tracking with ETA and speed information
- Multiple user interfaces (CLI, GUI, Web)
- Drag-and-drop file uploads in GUI and Web interfaces
- Auto-detection of archive types for extraction
- Auto-generated filenames based on the source
- Format selection (ZIP, TAR, GZ, BZ2, XZ)
- Real-time progress visualization

Planned features:
- Support for additional formats (TAR, GZ, BZ2, XZ, 7z)
- Encryption with AES-256
- Archive splitting
- Selective extraction
- Cloud integration
- Smart compression settings based on file types
- Batch processing and compression profiles
- Incremental archiving
- Archive integrity checking
- File preview

## Installation

### Prerequisites

- Go 1.17 or higher
- For Native GUI: A working graphics environment (X11, Wayland, macOS, or Windows)
- For Web UI: Node.js 18.x or higher and npm/yarn

### Building from source

1. Clone the repository:
```
git clone https://github.com/your-username/file-compressor.git
cd file-compressor
```

2. Build the command-line application:
```
make build
```

3. Build the GUI application:
```
make build-gui
```

4. For the Web UI, install dependencies:
```
cd web-ui
npm install
```

## Usage

### Command-line Interface

Compress a file or directory:
```
./build/file-compressor compress <source> <destination> [format]
```

Extract an archive:
```
./build/file-compressor extract <source> <destination>
```

### Native GUI Application

Run the GUI application:
```
./build/file-compressor-gui
```

### Web-based Interface

1. Start the API server:
```
go run ./cmd/api
```

2. Start the Next.js frontend (in a separate terminal):
```
cd web-ui
npm run dev
```

3. Open your browser and navigate to `http://localhost:3000`

Alternatively, you can use the provided script to run both the API server and the web frontend:
```
./run.sh
```

## Project Structure

- `cmd/file-compressor`: CLI application
- `cmd/gui`: Desktop GUI application (using Fyne)
- `cmd/api`: REST API server for the web interface
- `pkg/archiver`: Core compression/extraction functionality
- `pkg/utils`: Utility functions like progress tracking
- `web-ui`: Next.js-based web interface
- `build`: Compiled executables
- `uploads`: Temporary storage for uploaded files
- `compressed`: Storage for compressed output files

## Maintenance

To clean up temporary files and build artifacts:

```bash
./clean.sh
```

## Development

Run the applications in development mode:

```
# Run CLI
make run

# Run GUI
make run-gui

# Run Web UI (both API and frontend)
./run.sh
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 