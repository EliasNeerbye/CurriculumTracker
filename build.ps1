# Curriculum Tracker Build Script for Windows
param(
    [Parameter(Position=0)]
    [string]$Target = "build"
)

$ErrorActionPreference = "Stop"

# Variables
$BinaryName = "curriculum-tracker.exe"
$ServerCmd = "./cmd/server"
$WasmCmd = "./cmd/wasm"
$WebDir = "./web"
$BinDir = "./bin"

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  .\build.ps1 build     - Build both server and WebAssembly"
    Write-Host "  .\build.ps1 server    - Build only the server"
    Write-Host "  .\build.ps1 wasm      - Build only WebAssembly frontend"
    Write-Host "  .\build.ps1 run       - Build and run the application"
    Write-Host "  .\build.ps1 clean     - Clean build artifacts"
    Write-Host "  .\build.ps1 setup     - Set up the database"
    Write-Host "  .\build.ps1 test      - Run tests"
    Write-Host "  .\build.ps1 prod      - Build optimized production version"
    Write-Host "  .\build.ps1 help      - Show this help message"
}

function Build-Server {
    Write-Host "Building server..." -ForegroundColor Yellow
    
    if (!(Test-Path $BinDir)) {
        New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
    }
    
    go build -o "$BinDir/$BinaryName" $ServerCmd
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Server built successfully!" -ForegroundColor Green
    } else {
        throw "Server build failed"
    }
}

function Build-Wasm {
    Write-Host "Building WebAssembly..." -ForegroundColor Yellow
    
    if (!(Test-Path $WebDir)) {
        New-Item -ItemType Directory -Path $WebDir -Force | Out-Null
    }
    
    $env:GOOS = "js"
    $env:GOARCH = "wasm"
    
    go build -o "$WebDir/main.wasm" $WasmCmd
    if ($LASTEXITCODE -ne 0) {
        throw "WebAssembly build failed"
    }
      # Copy wasm_exec.js from Go installation
    $goRoot = go env GOROOT
    
    # Try multiple possible locations for wasm_exec.js
    $wasmExecPaths = @(
        (Join-Path $goRoot "lib/wasm/wasm_exec.js"),      # Windows Go 1.21+
        (Join-Path $goRoot "misc/wasm/wasm_exec.js")      # Older versions or Linux/Mac
    )
    
    $wasmExecFound = $false
    foreach ($wasmExecPath in $wasmExecPaths) {
        if (Test-Path $wasmExecPath) {
            Copy-Item $wasmExecPath "$WebDir/wasm_exec.js" -Force
            Write-Host "WebAssembly built successfully!" -ForegroundColor Green
            $wasmExecFound = $true
            break
        }
    }
    
    if (-not $wasmExecFound) {
        throw "Could not find wasm_exec.js in Go installation at any expected location"
    }
    
    # Reset environment variables
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
}

function Clean-Build {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
    
    if (Test-Path $BinDir) {
        Remove-Item $BinDir -Recurse -Force
    }
    
    if (Test-Path "$WebDir/main.wasm") {
        Remove-Item "$WebDir/main.wasm" -Force
    }
    
    if (Test-Path "$WebDir/wasm_exec.js") {
        Remove-Item "$WebDir/wasm_exec.js" -Force
    }
    
    Write-Host "Clean complete!" -ForegroundColor Green
}

function Setup-Database {
    Write-Host "Setting up database..." -ForegroundColor Yellow
    
    # Check if PostgreSQL is available
    try {
        $null = Get-Command psql -ErrorAction Stop
    } catch {
        throw "PostgreSQL client (psql) is required but not found in PATH"
    }
    
    # Create database (ignore error if it already exists)
    try {
        createdb curriculum_tracker 2>$null
    } catch {
        Write-Host "Database may already exist, continuing..." -ForegroundColor Yellow
    }
    
    # Run schema
    psql curriculum_tracker -f database/schema.sql
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Database setup complete!" -ForegroundColor Green
    } else {
        throw "Database setup failed"
    }
}

function Run-Tests {
    Write-Host "Running tests..." -ForegroundColor Yellow
    go test ./...
}

function Run-Application {
    Build-Server
    Build-Wasm
    
    Write-Host "Starting server..." -ForegroundColor Yellow
    Write-Host "Server will be available at http://localhost:8080" -ForegroundColor Cyan
    & "$BinDir/$BinaryName"
}

function Build-Production {
    Clean-Build
    
    Write-Host "Building for production..." -ForegroundColor Yellow
    
    if (!(Test-Path $BinDir)) {
        New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
    }
    
    # Build optimized server
    $env:CGO_ENABLED = "0"
    go build -ldflags="-w -s" -o "$BinDir/$BinaryName" $ServerCmd
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    
    # Build optimized WebAssembly
    $env:GOOS = "js"
    $env:GOARCH = "wasm"
    go build -ldflags="-w -s" -o "$WebDir/main.wasm" $WasmCmd
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
      # Copy wasm_exec.js
    $goRoot = go env GOROOT
    
    # Try multiple possible locations for wasm_exec.js
    $wasmExecPaths = @(
        (Join-Path $goRoot "lib/wasm/wasm_exec.js"),      # Windows Go 1.21+
        (Join-Path $goRoot "misc/wasm/wasm_exec.js")      # Older versions or Linux/Mac
    )
    
    $wasmExecFound = $false
    foreach ($wasmExecPath in $wasmExecPaths) {
        if (Test-Path $wasmExecPath) {
            Copy-Item $wasmExecPath "$WebDir/wasm_exec.js" -Force
            $wasmExecFound = $true
            break
        }
    }
    
    if (-not $wasmExecFound) {
        throw "Could not find wasm_exec.js in Go installation at any expected location"
    }
    
    Write-Host "Production build complete!" -ForegroundColor Green
}

# Main execution
try {
    switch ($Target.ToLower()) {
        "build" { 
            Build-Server
            Build-Wasm
        }
        "server" { Build-Server }
        "wasm" { Build-Wasm }
        "run" { Run-Application }
        "clean" { Clean-Build }
        "setup" { Setup-Database }
        "test" { Run-Tests }
        "prod" { Build-Production }
        "help" { Show-Help }
        default { 
            Write-Host "Unknown target: $Target" -ForegroundColor Red
            Show-Help
            exit 1
        }
    }
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}
