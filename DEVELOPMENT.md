# ThreadBolt Framework - Development Guide

This guide covers how to set up ThreadBolt for development, run tests, contribute to the framework, and publish releases.

## 🏗️ Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for convenience commands)
- Docker (for integration tests)

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/ThreadBolt/threadbolt.git
cd threadbolt

# Install dependencies
go mod tidy

# Build the CLI tool
go build -o bin/threadbolt cmd/threadbolt/main.go

# Add to PATH (optional)
export PATH=$PWD/bin:$PATH
```

### Development Workflow

#### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

#### 2. Make Changes

Follow the project structure:
- Framework core: `pkg/framework/`
- CLI commands: `pkg/cli/`
- Code generation: `pkg/generator/`
- ORM helpers: `pkg/orm/`
- Configuration: `pkg/config/`

#### 3. Test Your Changes

```bash
# Run unit tests
go test ./...

# Run with coverage
go test ./... -cover

# Run integration tests
make test-integration

# Test CLI manually
./bin/threadbolt new test-app
cd test-app
../bin/threadbolt run
```

## 🧪 Testing

ThreadBolt has multiple levels of testing to ensure reliability.

### Unit Tests

Unit tests are located alongside the source code:

```bash
# Run all unit tests
go test ./...

# Run tests for specific package
go test ./pkg/framework

# Run with verbose output
go test ./... -v

# Run with coverage report
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests

Integration tests verify the framework works with real applications:

```bash
# Run integration tests
make test-integration

# Or manually:
./scripts/integration-test.sh
```

### Test Structure

#### Unit Test Example (`pkg/framework/app_test.go`)

```go
package framework

import (
    "testing"
    "github.com/spf13/viper"
)

func TestLoadApp(t *testing.T) {
    // Setup test environment
    setupTestStructure()
    defer cleanupTestStructure()
    
    app, err := LoadApp()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if app.Router == nil {
        t.Error("Expected router to be initialized")
    }
    
    if app.Config == nil {
        t.Error("Expected config to be loaded")
    }
}

func setupTestStructure() {
    // Create minimal test structure
    dirs := []string{
        "controllers", "models", "config", "routes",
        "internal/middleware", "internal/services",
        "migrations", "templates", "public", "tests",
    }
    
    for _, dir := range dirs {
        os.MkdirAll(dir, 0755)
    }
    
    // Create required files
    files := map[string]string{
        "main.go": "package main\nfunc main() {}",
        "go.mod": "module test\ngo 1.21",
        "config/config.yaml": "server:\n  port: 8080",
        "routes/routes.go": "package routes",
    }
    
    for path, content := range files {
        os.WriteFile(path, []byte(content), 0644)
    }
}
```

#### CLI Integration Test (`scripts/integration-test.sh`)

```bash
#!/bin/bash
set -e

echo "🧪 Running ThreadBolt Integration Tests"

# Build CLI
go build -o bin/threadbolt cmd/threadbolt/main.go

# Test project creation
echo "📁 Testing project creation..."
rm -rf test-integration-app
./bin/threadbolt new test-integration-app
cd test-integration-app

# Test model generation
echo "🏗️ Testing model generation..."
../bin/threadbolt generate model TestUser

# Test controller generation  
echo "🎮 Testing controller generation..."
../bin/threadbolt generate controller TestUser

# Test build
echo "🔨 Testing build..."
go mod tidy
go build -o test-app main.go

# Test server start (background)
echo "🚀 Testing server start..."
timeout 10s go run main.go &
SERVER_PID=$!

sleep 3

# Test endpoints
echo "🌐 Testing endpoints..."
curl -f http://localhost:8080/health || (echo "Health check failed" && exit 1)
curl -f http://localhost:8080/api/v1/status || (echo "Status check failed" && exit 1)

# Cleanup
kill $SERVER_PID 2>/dev/null || true
cd ..
rm -rf test-integration-app

echo "✅ All integration tests passed!"
```

### Performance Tests

Benchmark critical framework components:

```go
func BenchmarkAppStart(b *testing.B) {
    setupTestStructure()
    defer cleanupTestStructure()
    
    for i := 0; i < b.N; i++ {
        app, err := LoadApp()
        if err != nil {
            b.Fatal(err)
        }
        
        // Simulate app startup overhead
        _ = app.Router
    }
}

func BenchmarkRouteMatching(b *testing.B) {
    app := setupTestApp()
    
    req := httptest.NewRequest("GET", "/api/v1/users", nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rr := httptest.NewRecorder()
        app.Router.ServeHTTP(rr, req)
    }
}
```

## 🔧 Local Development Tools

### Makefile Commands

Create a `Makefile` for convenience:

```makefile
.PHONY: build test test-integration lint clean install

# Build CLI tool
build:
	go build -o bin/threadbolt cmd/threadbolt/main.go

# Run unit tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Run integration tests
test-integration:
	./scripts/integration-test.sh

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install CLI tool locally
install: build
	cp bin/threadbolt $(GOPATH)/bin/threadbolt

# Generate example application for testing
example:
	rm -rf example-app
	./bin/threadbolt new example-app
	cd example-app && ../bin/threadbolt generate model User
	cd example-app && ../bin/threadbolt generate controller User

# Run all checks (CI pipeline)
ci: test test-integration lint
```

### Development Scripts

#### Auto-rebuild Script (`scripts/dev-watch.sh`)

```bash
#!/bin/bash

# Watch for changes and rebuild
while true; do
    inotifywait -r -e modify --include='\.go$' pkg/ cmd/
    echo "🔨 Rebuilding..."
    make build
    echo "✅ Build complete"
done
```

#### Release Preparation (`scripts/prepare-release.sh`)

```bash
#!/bin/bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "🚀 Preparing release $VERSION"

# Run all tests
make ci

# Update version in code
sed -i "s/version = \".*\"/version = \"$VERSION\"/" pkg/framework/version.go

# Update CHANGELOG
echo "## [$VERSION] - $(date +%Y-%m-%d)" > CHANGELOG.tmp
echo "" >> CHANGELOG.tmp
echo "### Added" >> CHANGELOG.tmp
echo "### Changed" >> CHANGELOG.tmp
echo "### Fixed" >> CHANGELOG.tmp
echo "" >> CHANGELOG.tmp
cat CHANGELOG.md >> CHANGELOG.tmp
mv CHANGELOG.tmp CHANGELOG.md

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o dist/threadbolt-linux-amd64 cmd/threadbolt/main.go
GOOS=darwin GOARCH=amd64 go build -o dist/threadbolt-darwin-amd64 cmd/threadbolt/main.go
GOOS=windows GOARCH=amd64 go build -o dist/threadbolt-windows-amd64.exe cmd/threadbolt/main.go

echo "✅ Release $VERSION prepared"
echo "📝 Don't forget to update CHANGELOG.md with release notes"
```

## 📝 Code Style and Standards

### Go Code Standards

1. **Follow Go conventions**: Use `go fmt`, `go vet`
2. **Error handling**: Always handle errors explicitly
3. **Documentation**: All public functions must have comments
4. **Testing**: Maintain >80% test coverage
5. **Naming**: Use descriptive names following Go conventions

### Example Code Style

```go
// Package framework provides the core ThreadBolt application framework.
package framework

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
    "gorm.io/gorm"
)

// App represents a ThreadBolt application instance.
type App struct {
    Router    *mux.Router
    DB        *gorm.DB
    Config    *viper.Viper
    Container *di.Container
}

// NewApp creates a new ThreadBolt application instance.
// It initializes the router, database connection, and configuration.
func NewApp() (*App, error) {
    app := &App{
        Router:    mux.NewRouter(),
        Container: di.NewContainer(),
    }
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }
    app.Config = cfg
    
    return app, nil
}
```

### Commit Message Format

Follow conventional commits:

```
type(scope): description

[optional body]

[optional footer]
```

Examples:
- `feat(cli): add model generation command`
- `fix(orm): resolve connection pool issue`
- `docs(readme): update installation instructions`
- `test(framework): add integration tests for app loading`

## 🚀 Release Process

### Version Management

ThreadBolt follows semantic versioning (semver):
- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

### Release Steps

#### 1. Prepare Release Branch

```bash
git checkout -b release/v1.2.0
```

#### 2. Update Version and Changelog

```bash
# Update version
./scripts/prepare-release.sh v1.2.0

# Review and edit CHANGELOG.md
$EDITOR CHANGELOG.md
```

#### 3. Test Release

```bash
# Run comprehensive tests
make ci

# Test installation process
go install ./cmd/threadbolt
threadbolt new test-release-app
cd test-release-app
threadbolt run &
sleep 3
curl http://localhost:8080/health
pkill -f "threadbolt run"
cd .. && rm -rf test-release-app
```

#### 4. Create Pull Request

```bash
git add .
git commit -m "chore: prepare release v1.2.0"
git push origin release/v1.2.0
```

Create PR from release branch to main.

#### 5. Tag and Release

After PR is merged:

```bash
git checkout main
git pull origin main

# Create and push tag
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

#### 6. GitHub Release

Create GitHub release:
1. Go to GitHub Releases page
2. Click "Create a new release"
3. Select tag v1.2.0
4. Title: "ThreadBolt v1.2.0"
5. Description: Copy from CHANGELOG.md
6. Attach pre-built binaries from `dist/`
7. Publish release

#### 7. Update Documentation

```bash
# Update go.mod examples in documentation
find . -name "*.md" -exec sed -i 's/v1.1.0/v1.2.0/g' {} \;

git add .
git commit -m "docs: update version references to v1.2.0"
git push origin main
```

## 📦 Publishing

### Go Module Publishing

ThreadBolt is published as a Go module. After tagging:

```bash
# Verify module can be downloaded
go list -m github.com/ThreadBolt/threadbolt@v1.2.0

# Test installation
go install github.com/ThreadBolt/threadbolt/cmd/threadbolt@v1.2.0
```

### Distribution

#### Homebrew Formula (Optional)

Create a Homebrew formula for easy installation:

```ruby
# Formula/threadbolt.rb
class Threadbolt < Formula
  desc "Convention-over-configuration web framework for Go"
  homepage "https://github.com/ThreadBolt/threadbolt"
  url "https://github.com/ThreadBolt/threadbolt/archive/v1.2.0.tar.gz"
  sha256 "..."
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args, "./cmd/threadbolt"
  end

  test do
    system "#{bin}/threadbolt", "new", "test-app"
    assert_predicate testpath/"test-app/main.go", :exist?
  end
end
```

#### Docker Image (Optional)

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o threadbolt cmd/threadbolt/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/threadbolt /usr/local/bin/threadbolt
ENTRYPOINT ["threadbolt"]
```

## 🤝 Contributing Guidelines

### Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Create a feature branch** from main
4. **Make your changes** following the style guide
5. **Add tests** for new functionality
6. **Run the test suite** to ensure nothing breaks
7. **Commit your changes** with clear messages
8. **Push to your fork** and create a pull request

### Pull Request Process

1. **Description**: Provide a clear description of changes
2. **Testing**: Include tests for new features
3. **Documentation**: Update relevant documentation
4. **Breaking Changes**: Clearly document any breaking changes
5. **Review**: Address feedback from maintainers

### Code Review Checklist

- [ ] Code follows Go conventions and style guide
- [ ] All tests pass (`make ci`)
- [ ] New functionality is tested
- [ ] Documentation is updated
- [ ] No breaking changes (or clearly documented)
- [ ] Performance impact considered
- [ ] Security implications reviewed

## 🐛 Debugging

### Debug Mode

Enable debug mode for verbose logging:

```bash
export THREADBOLT_DEBUG=true
threadbolt run
```

### Common Issues

#### Project Structure Validation

```go
// Add debug logging to framework/app.go
func validateProjectStructure() error {
    requiredDirs := []string{...}
    
    for _, dir := range requiredDirs {
        if _, err := os.Stat(dir); os.IsNotExist(err) {
            log.Printf("DEBUG: Missing directory: %s", dir)
            return fmt.Errorf("required directory '%s' is missing", dir)
        }
        log.Printf("DEBUG: Found directory: %s", dir)
    }
    
    return nil
}
```

#### Database Connection Issues

```go
// Add connection testing
func testDatabaseConnection(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    return sqlDB.Ping()
}
```

### Performance Profiling

```bash
# Build with profiling
go build -o threadbolt-debug cmd/threadbolt/main.go

# Run with profiling
./threadbolt-debug run --cpuprofile=cpu.prof --memprofile=mem.prof

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

## 📊 Metrics and Monitoring

### Framework Metrics

Track framework usage and performance:

```go
// pkg/framework/metrics.go
type Metrics struct {
    AppStartTime    time.Time
    RequestCount    int64
    ErrorCount      int64
    DatabaseQueries int64
}

func (m *Metrics) RecordRequest() {
    atomic.AddInt64(&m.RequestCount, 1)
}

func (m *Metrics) RecordError() {
    atomic.AddInt64(&m.ErrorCount, 1)
}
```

### Health Checks

Implement comprehensive health checks:

```go
func (a *App) HealthCheck() map[string]interface{} {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now(),
        "version":   version.Get(),
        "uptime":    time.Since(a.startTime),
    }
    
    // Database health
    if sqlDB, err := a.DB.DB(); err == nil {
        if err := sqlDB.Ping(); err == nil {
            health["database"] = "healthy"
        } else {
            health["database"] = "unhealthy"
            health["status"] = "degraded"
        }
    }
    
    return health
}
```

## 📚 Documentation

### API Documentation

Generate API documentation:

```bash
# Install godoc
go install golang.org/x/tools/cmd/godoc@latest

# Generate docs
godoc -http=:6060

# Visit http://localhost:6060/pkg/github.com/ThreadBolt/threadbolt/
```

### Example Applications

Maintain comprehensive example applications in `examples/`:

- `examples/basic-api/` - Simple REST API
- `examples/blog-app/` - Full blog application
- `examples/microservice/` - Microservice example
- `examples/web-app/` - Traditional web application

Each example should include:
- Complete source code
- README with setup instructions
- Docker configuration
- Deployment instructions

---

**Happy contributing to ThreadBolt! 🎉**