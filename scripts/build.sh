#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Version handling
get_version() {
    local version

    # First try to get the latest tag
    version=$(git describe --tags 2>/dev/null)

    # If no tags exist, use the commit hash
    if [ $? -ne 0 ]; then
        version="v0.1.0-$(git rev-parse --short HEAD)"
        if [ -n "$(git status --porcelain)" ]; then
            version="${version}-dev"
        fi
    fi

    echo "$version"
}

# Version from git
VERSION=$(get_version)

# Platforms to build for
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Build directory
BUILD_DIR="build"
BINARY_NAME="apex"

# Print step with color
print_step() {
    echo -e "${BLUE}==> $1${NC}"
}

# Print success with color
print_success() {
    echo -e "${GREEN}==> $1${NC}"
}

# Print error with color
print_error() {
    echo -e "${RED}==> Error: $1${NC}"
}

# Print warning with color
print_warning() {
    echo -e "${RED}==> Warning: $1${NC}"
}

# Check if garble is installed
check_garble() {
    if ! command -v garble &> /dev/null; then
        print_step "Installing garble..."
        go clean -cache -modcache
        go install mvdan.cc/garble@latest
        if [ $? -ne 0 ]; then
            print_error "Failed to install garble"
            exit 1
        fi
    fi
}

# Clean build directory and caches
clean_build_dir() {
    print_step "Cleaning build directory and Go caches..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    go clean -cache -modcache
    rm -rf ~/Library/Caches/garble
}

# Compress binary with UPX if possible
compress_binary() {
    local binary=$1
    local os=$2
    local arch=$3

    # Skip compression for unsupported platforms
    if [ "$os" = "darwin" ] || [ "$arch" = "arm64" ]; then
        print_warning "Skipping UPX compression for $os/$arch (not supported)"
        return 0
    fi

    if command -v upx &> /dev/null; then
        print_step "Compressing ${binary} with UPX..."
        if upx -9 --no-progress "$binary" &> /dev/null; then
            print_success "Successfully compressed ${binary}"
        else
            print_warning "Failed to compress ${binary} (compression skipped)"
        fi
    fi
}

# Build for all platforms
build_all() {
    print_step "Building APEX version: $VERSION"
    
    # Create build directory if it doesn't exist
    mkdir -p "$BUILD_DIR"
    
    # Build for each platform
    for PLATFORM in "${PLATFORMS[@]}"; do
        # Split platform into OS and ARCH
        IFS="/" read -r OS ARCH <<< "$PLATFORM"
        
        # Set binary name based on OS
        if [ "$OS" = "windows" ]; then
            BINARY_SUFFIX=".exe"
        else
            BINARY_SUFFIX=""
        fi
        
        # Set the output binary path
        OUTPUT_NAME="${BUILD_DIR}/${BINARY_NAME}-${OS}-${ARCH}-${VERSION}${BINARY_SUFFIX}"
        
        print_step "Building for ${OS}/${ARCH}..."
        
        # Set GOOS and GOARCH
        export GOOS=$OS
        export GOARCH=$ARCH
        
        # Clean Go cache before each build
        go clean -cache
        
        # Build with garble
        garble -tiny -seed=random build -mod=mod -trimpath -ldflags="-s -w -X main.Version=${VERSION}" -o "$OUTPUT_NAME" .
        
        if [ $? -eq 0 ]; then
            print_success "Successfully built ${OUTPUT_NAME}"
            
            # Create SHA256 checksum
            if command -v shasum &> /dev/null; then
                shasum -a 256 "$OUTPUT_NAME" > "${OUTPUT_NAME}.sha256"
                print_success "Generated SHA256 checksum for ${OUTPUT_NAME}"
            fi
            
            # Try to compress the binary
            compress_binary "$OUTPUT_NAME" "$OS" "$ARCH"
        else
            print_error "Failed to build for ${OS}/${ARCH}"
            exit 1
        fi
    done
}

# Main execution
main() {
    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        print_error "Please run this script from the root of the project"
        exit 1
    fi
    
    # Check dependencies
    check_garble
    
    # Clean build directory and caches
    clean_build_dir
    
    # Build for all platforms
    build_all
    
    print_success "Build process completed successfully!"
    print_success "Binaries are available in the ${BUILD_DIR} directory"
}

# Run main function
main 