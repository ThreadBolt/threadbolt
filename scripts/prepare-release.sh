# #!/bin/bash
# set -e

# VERSION=$1
# if [ -z "$VERSION" ]; then
#     echo "Usage: $0 <version>"
#     exit 1
# fi

# echo "🚀 Preparing release $VERSION"


# # Update version in code
# sed -i "s/version = \".*\"/version = \"$VERSION\"/" pkg/framework/version.go

# # Update CHANGELOG
# echo "## [$VERSION] - $(date +%Y-%m-%d)" > CHANGELOG.tmp
# echo "" >> CHANGELOG.tmp
# echo "### Added" >> CHANGELOG.tmp
# echo "### Changed" >> CHANGELOG.tmp
# echo "### Fixed" >> CHANGELOG.tmp
# echo "" >> CHANGELOG.tmp
# cat CHANGELOG.md >> CHANGELOG.tmp
# mv CHANGELOG.tmp CHANGELOG.md

# # Build for multiple platforms
# GOOS=linux GOARCH=amd64 go build -o dist/threadbolt-linux-amd64 cmd/threadbolt/main.go
# GOOS=darwin GOARCH=amd64 go build -o dist/threadbolt-darwin-amd64 cmd/threadbolt/main.go
# GOOS=windows GOARCH=amd64 go build -o dist/threadbolt-windows-amd64.exe cmd/threadbolt/main.go

# echo "✅ Release $VERSION prepared"
# echo "📝 Don't forget to update CHANGELOG.md with release notes"

#!/usr/bin/env bash
set -euo pipefail

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "🚀 Preparing release $VERSION"


# Create pkg/framework directory if it doesn't exist
mkdir -p pkg/framework

# Create version.go if it doesn't exist
if [ ! -f "pkg/framework/version.go" ]; then
    echo "📝 Creating pkg/framework/version.go..."
    cat << EOF > pkg/framework/version.go
// Package framework provides the core ThreadBolt application framework.
package framework

// Version holds the current version of the ThreadBolt framework.
const Version = "0.0.0"

// Get returns the current framework version.
func Get() string {
    return Version
}
EOF
fi

# Update version in code
echo "📝 Updating version in pkg/framework/version.go..."
sed -i "s/const Version = \".*\"/const Version = \"$VERSION\"/" pkg/framework/version.go

# Create CHANGELOG.md if it doesn't exist
if [ ! -f "CHANGELOG.md" ]; then
    echo "📝 Creating CHANGELOG.md..."
    touch CHANGELOG.md
fi

# Update CHANGELOG
echo "📝 Updating CHANGELOG.md..."
DATE=$(date +%Y-%m-%d)
{
    echo "## [$VERSION] - $DATE"
    echo ""
    echo "### Added"
    echo ""
    echo "### Changed"
    echo ""
    echo "### Fixed"
    echo ""
} > CHANGELOG.tmp
cat CHANGELOG.md >> CHANGELOG.tmp
mv CHANGELOG.tmp CHANGELOG.md

# Create dist directory if it doesn't exist
mkdir -p dist

# Build for multiple platforms
echo "🔨 Building binaries for multiple platforms..."
GOOS=linux GOARCH=amd64 go build -o dist/threadbolt-linux-amd64 cmd/threadbolt/main.go
GOOS=darwin GOARCH=amd64 go build -o dist/threadbolt-darwin-amd64 cmd/threadbolt/main.go
GOOS=windows GOARCH=amd64 go build -o dist/threadbolt-windows-amd64.exe cmd/threadbolt/main.go

# Update README.md with version badge
if [ -f "README.md" ]; then
    if grep -q "![Version]" README.md; then
        echo "📝 Updating existing version badge in README.md..."
        sed -i "s|version-[0-9.]\+|version-${VERSION}|g" README.md
    else
        echo "📝 Adding version badge to README.md..."
        sed -i "1a\\\n![Version](https://img.shields.io/badge/version-${VERSION}-blue.svg)" README.md
    fi
else
    echo "📝 Creating README.md with version badge..."
    cat << EOF > README.md
# ThreadBolt
A convention-over-configuration web framework for Go

![Version](https://img.shields.io/badge/version-${VERSION}-blue.svg)

## Installation
\`\`\`bash
go install github.com/ThreadBolt/threadbolt/cmd/threadbolt@$VERSION
\`\`\`

## Documentation
See [ThreadBolt Development Guide](ThreadBolt-Development-Guide.md) for details.
EOF
fi

# Run go mod tidy
if [ -f "go.mod" ]; then
    echo "🧹 Running go mod tidy in root..."
    go mod tidy
fi
if [ -f "example-app/go.mod" ]; then
    echo "🧹 Running go mod tidy in example-app..."
    (cd example-app && go mod tidy)
fi

# Commit changes
echo "📦 Committing release changes..."
git add .
git commit -m "chore: prepare release $VERSION" || echo "⚠️ Nothing to commit"

# Create and push tag
echo "🏷️ Creating and pushing git tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"
git push origin HEAD
git push origin "$VERSION"

echo "✅ Release $VERSION prepared"
echo "📝 Please review and update CHANGELOG.md with detailed release notes"