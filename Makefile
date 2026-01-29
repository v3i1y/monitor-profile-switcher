APP := monitor-switcher
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP).exe

.PHONY: all build install fmt vet test tidy clean

all: build

build:
	go build -o $(BIN) ./cmd/monitor-switcher

install: build
	powershell -NoProfile -Command '& { $$dest = Join-Path $$env:USERPROFILE "bin"; if (!(Test-Path $$dest)) { New-Item -ItemType Directory -Path $$dest | Out-Null }; Copy-Item -Force "$(BIN)" (Join-Path $$dest "$(APP).exe"); $$path = [Environment]::GetEnvironmentVariable("Path","User"); if ($$path -notlike ("*" + $$dest + "*")) { [Environment]::SetEnvironmentVariable("Path", $$path + ";" + $$dest, "User"); Write-Output "Added to user PATH. Restart your terminal." } else { Write-Output "Already on user PATH." } }'

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

clean:
	powershell -NoProfile -Command "if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force '$(BIN_DIR)' }"
