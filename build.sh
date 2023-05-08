#!/bin/bash

# 设置项目名
PROJECT_NAME="rare_eth"

# 构建 amd64 macOS 二进制文件
echo "Building for amd64 macOS..."
GOOS=darwin
GOARCH=amd64
OUTPUT="${PROJECT_NAME}_${GOARCH}_macos"
env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
echo "Build complete: $OUTPUT"
echo ""

# 构建 arm64 macOS 二进制文件
echo "Building for arm64 macOS..."
GOOS=darwin
GOARCH=arm64
OUTPUT="${PROJECT_NAME}_${GOARCH}_macos"
env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
echo "Build complete: $OUTPUT"
echo ""

# 构建 Windows 二进制文件
echo "Building for Windows..."
GOOS=windows
GOARCH=amd64
OUTPUT="${PROJECT_NAME}_${GOARCH}_windows.exe"
env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
echo "Build complete: $OUTPUT"
echo ""

# 构建 amd64 Linux 二进制文件
echo "Building for amd64 Linux..."
GOOS=linux
GOARCH=amd64
OUTPUT="${PROJECT_NAME}_${GOARCH}_linux"
env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
echo "Build complete: $OUTPUT"
echo ""

# 构建 arm64 Linux 二进制文件
echo "Building for arm64 Linux..."
GOOS=linux
GOARCH=arm64
OUTPUT="${PROJECT_NAME}_${GOARCH}_linux"
env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
echo "Build complete: $OUTPUT"
echo ""
