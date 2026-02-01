# Monitor Profile Switcher (Go CLI)

A Windows-only CLI tool to save and switch multi-monitor display configurations using the Windows CCD (Connecting and Configuring Displays) APIs. This is a Go reimplementation of the original MonitorSwitcher CLI, without the GUI.

I built this version primarily because the original version wasn't working with VDD pretty well.

## Features

- Save the current active display configuration to a JSON profile.
- Load a profile and apply it using `SetDisplayConfig`.
- Human-readable summary output (`-print`).
- Virtual display (VDD) aware: includes virtual-mode metadata when saving profiles.
- Missing targets are ignored by default with warning logs.

## Requirements

- Windows 10/11.
- Go 1.20+ recommended (tested with 1.24).

## Build

```powershell
# From repo root
make build
# or
go build -o .\monitor-switcher.exe .\cmd\monitor-switcher
```

## Install

```powershell
make install
```

This copies the binary to `%USERPROFILE%\bin` and adds that folder to your user `PATH` (restart your terminal after running).

## Usage

```text
monitor-switcher.exe -save:Profile.monitorprofile
monitor-switcher.exe -load:Profile.monitorprofile
monitor-switcher.exe -debug -load:Profile.json
monitor-switcher.exe -print
```

### Flags

- `-save:{file}` Save the current active display configuration to a profile file.
- `-load:{file}` Load and apply a profile file.
- `-print` Print a human-readable summary of the current configuration.
- `-debug` Enable debug output (use before `-save`/`-load`).
- `-noidmatch` Disable adapter-ID matching (advanced).
- `-v` Enable virtual desktop injection (advanced).

### Missing targets

If a profile references a target that is not currently present, Windows may return an error. Virtual targets must be present for their paths to apply.

### Default profile location

If you pass a filename without a path (e.g. `-save:MyProfile`), profiles are stored under:

```
%USERPROFILE%\Monitor Profiles
```
If you omit the extension, `.monitorprofile` is added automatically.

## JSON profile format (overview)

Profiles contain three arrays:

- `pathInfo`: display paths and their settings (rotation, scaling, refresh rate, etc.)
- `modeInfo`: source/target modes and virtual desktop image entries
- `additionalInfo`: friendly names and device paths

The format mirrors the structures returned by `QueryDisplayConfig`.

## Development

Useful commands from the Makefile:

- `make build` Build the binary into `bin/`.
- `make fmt`   Run `go fmt`.
- `make vet`   Run `go vet`.
- `make test`  Run `go test`.
- `make tidy`  Run `go mod tidy`.

## Notes on Virtual Displays (VDD)

Virtual targets must be present/enumerated by Windows for their paths to apply. If a VDD is not active, its target will be ignored, and only remaining targets will be applied.

## License

Original project: MonitorSwitcher by Martin Kr?mer (MPL-2.0).
This Go reimplementation is provided as-is.
