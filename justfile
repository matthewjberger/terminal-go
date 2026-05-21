set windows-shell := ["powershell.exe"]

# Displays the list of available commands
@just:
    just --list

# Runs the named game (default: zork). Example: `just run zork`.
run $project="zork":
    go run ./cmd/{{project}}

# Builds the named game's binary.
build $project="zork":
    go build ./cmd/{{project}}

# Runs go vet and fails on unformatted files (Windows)
[windows]
check:
    go vet ./...
    $unformatted = (gofmt -l . | Out-String).Trim(); if ($unformatted) { Write-Host $unformatted; exit 1 }

# Runs go vet and fails on unformatted files (Unix)
[unix]
check:
    go vet ./...
    unformatted="$(gofmt -l .)"; if [ -n "$unformatted" ]; then echo "$unformatted"; exit 1; fi

# Formats all Go files
format:
    gofmt -w .

# Runs all tests
test:
    go test ./...

# Runs check + test (use this before pushing)
ci: check test

# Lists all module dependencies with available updates
outdated:
    go list -m -u all

# Shows what `go mod tidy` would change without applying it
tidy-check:
    go mod tidy -diff

# Tidies go.mod / go.sum
tidy:
    go mod tidy

# Runs every read-only check: vet+fmt, tidy diff, outdated, tests
audit: check tidy-check outdated test

# Removes any built binaries (Windows)
[windows]
clean:
    Remove-Item -Force -ErrorAction SilentlyContinue zork.exe

# Removes any built binaries (Unix)
[unix]
clean:
    rm -f zork zork.exe

# Displays Go tool version
@versions:
    go version
