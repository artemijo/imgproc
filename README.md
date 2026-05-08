# imgproc

A Go-based image processing pipeline with support for resizing, cropping, and color correction.

## Prerequisites

- Go 1.26.2 or later

## Compilation

### Windows (PowerShell)
```powershell
go build -o imgproc.exe
```

### Windows (Command Prompt)
```cmd
go build -o imgproc.exe
```

### Linux/macOS
```bash
go build -o imgproc
```

### Cross-Compilation

#### Build for Windows from Linux/macOS
```bash
GOOS=windows GOARCH=amd64 go build -o imgproc.exe
```

#### Build for Linux from Windows/macOS
```bash
GOOS=linux GOARCH=amd64 go build -o imgproc
```

#### Build for macOS from Windows/Linux
```bash
GOOS=darwin GOARCH=amd64 go build -o imgproc
```

## Usage

Simply double click on the executable program and done! 
But before that put all of the images need covertation in the *input* folder.

## Features

- Image resizing
- Image cropping
- Color correction
