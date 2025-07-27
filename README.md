# GPMF Parser

A Go-based parser for extracting and processing GoPro Metadata Format (GPMF) data from GoPro MP4 files. This tool allows you to extract various types of telemetry and metadata from GoPro video files.

## Features

The parser supports extraction of the following telemetry data:

- 🌍 **GPS Data**: Location coordinates and movement tracking
- 🚀 **Acceleration**: Device acceleration measurements
- ☀️ **Luminance**: Light level measurements
- 🎨 **Hue**: Color information from the video
- 😀 **Face Detection**: Face detection metadata
- 🎬 **Scene Analysis**: Scene detection and analysis
- 📊 **Gyroscope**: Device orientation and movement data

## Installation

```bash
go get github.com/chrigu/gopro-meta
```

## Usage

### Command Line Interface

#### Extract raw GPMF

To extract binary GPMF data from a GoPro MP4 file:

```bash
go run ./gopro-mf extract <mp4 file>
```

#### Print telemetry data for debugging

Print telemetry data from a GoPro MP4 file:

```bash
go run ./gopro-mf print <mp4 file>
```

#### Help

See help

```bash
go run ./gopro-mf help
```

### WebAssembly Support

The parser can be compiled to WebAssembly for use in web applications:

```bash
GOOS=js GOARCH=wasm go build -o wasm/main.wasm ./wasm
```

Two functions are exposed:

- Export of binary GPMF file
- Export of telemetry data as JSON

See https://github.com/chrigu/trailtrace for a complete example.

#### Binary GPMF file

```typescript
const gpmfData = await window.exportGPMF(file) as Uint8Array;
```

#### Telemetry data

```typescript
const metadata = await window.processFile(file) as {gpsData: GpsData[], gyroData: any[], faceData: any[], lumaData: any[], hueData: any[], sceneData: any[]};
```

## Project Structure

- `/parser`: Core parsing functionality for different types of metadata
- `/telemetry`: Data structures and types for various telemetry data
- `/wasm`: WebAssembly implementation
- `/internal`: Internal utilities and helpers
- `/mp4`: MP4 file handling utilities

## Development

The project is written in Go and uses the following key dependencies:
- `github.com/abema/go-mp4`: For MP4 file handling
- `github.com/urfave/cli`: For command-line interface
- Various other utilities for data processing and formatting

## Todos

- Rename `processFile` to something meaningful
- Add types for WASM export
- Refactor timed data, extraction
- Test older GoPros
- Select metadata to export
- Tests
- Refactoring
- Optimize performance https://goperf.dev/01-common-patterns/mem-prealloc/#why-preallocation-matters

## Resources
- https://github.com/gopro/gpmf-parser
- https://www.trekview.org/blog/injecting-camm-gpmd-telemetry-videos-part-3-mp4-structure-telemetry-trak/
- https://developer.apple.com/documentation/quicktime-file-format/sample-to-chunk_atom/sample-to-chunk_table

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
