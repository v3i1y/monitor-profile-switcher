package switcher

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"monitor-profile-switcher/internal/ccd"
	"monitor-profile-switcher/internal/profile"
)

const (
	applyFlags = ccd.SdcFlagsApply | ccd.SdcFlagsUseSuppliedDisplayConfig | ccd.SdcFlagsSaveToDatabase | ccd.SdcFlagsNoOptimization | ccd.SdcFlagsAllowChanges
)

func SaveProfile(path string, debug bool) error {
	debugf(debug, "Saving profile to: %s", path)

	paths, modes, additional, err := ccd.GetDisplaySettingsWithFlags(ccd.QueryDisplayFlagsOnlyActivePaths | ccd.QueryDisplayFlagsVirtualModeAware)
	if err != nil {
		debugf(debug, "VirtualModeAware query failed, falling back to standard query: %v", err)
		paths, modes, additional, err = ccd.GetDisplaySettings(true)
		if err != nil {
			return fmt.Errorf("get display settings: %w", err)
		}
	}

	prof := profileFromCCD(paths, modes, additional)
	if err := profile.Save(path, prof); err != nil {
		return err
	}
	return nil
}

func LoadProfile(path string, debug bool, noIDMatch bool, virtualInject bool) error {
	debugf(debug, "Loading profile from: %s", path)

	prof, err := profile.Load(path)
	if err != nil {
		return err
	}

	paths, modes, additional := ccdFromProfile(prof)
	origPaths := append([]ccd.DisplayConfigPathInfo(nil), paths...)
	origModes := append([]ccd.DisplayConfigModeInfo(nil), modes...)

	virtualAware := profileHasVirtualDisplay(prof)

	currentPaths, currentModes, currentAdditional, err := ccd.GetDisplaySettingsWithFlags(queryFlagsForProfile(false, virtualAware))
	if err != nil {
		return fmt.Errorf("get current display settings: %w", err)
	}

	if !noIDMatch {
		debugf(debug, "Matching adapter IDs for path info")
		for i := range paths {
			for j := range currentPaths {
				if paths[i].SourceInfo.ID == currentPaths[j].SourceInfo.ID &&
					paths[i].TargetInfo.ID == currentPaths[j].TargetInfo.ID {
					paths[i].SourceInfo.AdapterID = currentPaths[j].SourceInfo.AdapterID
					paths[i].TargetInfo.AdapterID = currentPaths[j].TargetInfo.AdapterID
					break
				}
			}
		}

		debugf(debug, "Matching adapter IDs for mode info")
		for i := range modes {
			for j := range paths {
				if modes[i].ID == paths[j].TargetInfo.ID && modes[i].InfoType == ccd.DisplayConfigModeInfoTypeTarget {
					for k := range modes {
						if modes[k].ID == paths[j].SourceInfo.ID &&
							modes[k].AdapterID.LowPart == modes[i].AdapterID.LowPart &&
							modes[k].InfoType == ccd.DisplayConfigModeInfoTypeSource {
							modes[k].AdapterID = paths[j].SourceInfo.AdapterID
							break
						}
					}
					modes[i].AdapterID = paths[j].TargetInfo.AdapterID
					break
				}
			}
		}
	}

	if virtualInject {
		if ensureDesktopImageModes(&paths, &modes, currentModes) {
			debugf(debug, "Injected missing desktop image info from current configuration")
		}
	}

	flags := applyFlags
	if virtualAware {
		flags |= ccd.SdcFlagsVirtualModeAware
	}

	if err := ccd.SetDisplayConfig(paths, modes, flags); err != nil {
		debugf(debug, "Primary SetDisplayConfig failed: %v", err)
		if len(currentAdditional) > 0 && len(additional) > 0 {
			debugf(debug, "Trying alternative matching method")
			paths = append([]ccd.DisplayConfigPathInfo(nil), origPaths...)
			modes = append([]ccd.DisplayConfigModeInfo(nil), origModes...)

			for i := range modes {
				for j := range currentAdditional {
					if currentAdditional[j].MonitorFriendlyDevice == "" || additional[i].MonitorFriendlyDevice == "" {
						continue
					}
					if currentAdditional[j].MonitorFriendlyDevice == additional[i].MonitorFriendlyDevice {
						originalID := modes[i].AdapterID
						for p := range paths {
							if paths[p].TargetInfo.AdapterID.LowPart == originalID.LowPart &&
								paths[p].TargetInfo.AdapterID.HighPart == originalID.HighPart {
								paths[p].TargetInfo.AdapterID = currentModes[j].AdapterID
								paths[p].SourceInfo.AdapterID = currentModes[j].AdapterID
								paths[p].TargetInfo.ID = currentModes[j].ID
							}
						}
						for k := range modes {
							if modes[k].AdapterID.LowPart == originalID.LowPart &&
								modes[k].AdapterID.HighPart == originalID.HighPart {
								modes[k].AdapterID = currentModes[j].AdapterID
							}
						}
						modes[i].AdapterID = currentModes[j].AdapterID
						modes[i].ID = currentModes[j].ID
						break
					}
				}
			}

			if err := ccd.SetDisplayConfig(paths, modes, flags); err != nil {
				if virtualAware {
					mergedPaths, mergedModes, ok := mergeProfileWithCurrent(origPaths, origModes, currentPaths, currentModes)
					if ok {
						debugf(debug, "Trying virtual-mode merge fallback")
						if mergeErr := ccd.SetDisplayConfig(mergedPaths, mergedModes, flags); mergeErr == nil {
							return nil
						} else {
							debugf(debug, "Merge fallback failed: %v", mergeErr)
						}
					}
				}
				return fmt.Errorf("SetDisplayConfig failed (alternative): %w", err)
			}
			return nil
		}
		return fmt.Errorf("SetDisplayConfig failed: %w", err)
	}
	return nil
}

func PrintSummary(w io.Writer) error {
	paths, modes, additional, err := ccd.GetDisplaySettingsWithFlags(ccd.QueryDisplayFlagsOnlyActivePaths | ccd.QueryDisplayFlagsVirtualModeAware)
	if err != nil {
		paths, modes, additional, err = ccd.GetDisplaySettings(true)
		if err != nil {
			return err
		}
	}

	summary := formatSummary(paths, modes, additional)
	_, err = io.WriteString(w, summary)
	return err
}

func debugf(enabled bool, format string, args ...any) {
	if !enabled {
		return
	}
	fmt.Printf(format+"\n", args...)
}

func formatSummary(paths []ccd.DisplayConfigPathInfo, modes []ccd.DisplayConfigModeInfo, additional []ccd.MonitorAdditionalInfo) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Active display paths: %d\n", len(paths))

	for i, path := range paths {
		active := (path.Flags & uint32(ccd.DisplayConfigFlagPathActive)) != 0
		state := "inactive"
		if active {
			state = "active"
		}
		fmt.Fprintf(&builder, "Path %d (%s)\n", i+1, state)

		targetIdx := int(path.TargetInfo.ModeInfoIdx)
		targetName := "Unknown"
		if targetIdx >= 0 && targetIdx < len(additional) {
			if additional[targetIdx].Valid && additional[targetIdx].MonitorFriendlyDevice != "" {
				targetName = additional[targetIdx].MonitorFriendlyDevice
			}
		}

		fmt.Fprintf(&builder, "  Target: %s (id %d, adapter %s)\n", targetName, path.TargetInfo.ID, formatAdapterID(path.TargetInfo.AdapterID))

		if targetIdx >= 0 && targetIdx < len(modes) && modes[targetIdx].InfoType == ccd.DisplayConfigModeInfoTypeTarget {
			targetMode := modes[targetIdx].TargetMode()
			refresh := formatRefreshRate(targetMode.TargetVideoSignalInfo.VSyncFreq)
			fmt.Fprintf(&builder, "  Refresh: %s\n", refresh)
			fmt.Fprintf(&builder, "  Active size: %dx%d\n", targetMode.TargetVideoSignalInfo.ActiveSize.Cx, targetMode.TargetVideoSignalInfo.ActiveSize.Cy)
		}

		sourceIdx := int(path.SourceInfo.ModeInfoIdx)
		if sourceIdx >= 0 && sourceIdx < len(modes) && modes[sourceIdx].InfoType == ccd.DisplayConfigModeInfoTypeSource {
			sourceMode := modes[sourceIdx].SourceMode()
			fmt.Fprintf(&builder, "  Source: %dx%d @ (%d,%d), pixel format %d\n", sourceMode.Width, sourceMode.Height, sourceMode.Position.X, sourceMode.Position.Y, sourceMode.PixelFormat)
		}

		fmt.Fprintf(&builder, "  Rotation: %d, Scaling: %d, TargetAvailable: %t\n", path.TargetInfo.Rotation, path.TargetInfo.Scaling, path.TargetInfo.TargetAvailable != 0)
	}

	return builder.String()
}

func formatRefreshRate(r ccd.DisplayConfigRational) string {
	if r.Denominator == 0 {
		return fmt.Sprintf("%d/%d Hz", r.Numerator, r.Denominator)
	}
	hz := float64(r.Numerator) / float64(r.Denominator)
	return fmt.Sprintf("%.2f Hz (%d/%d)", hz, r.Numerator, r.Denominator)
}

func formatAdapterID(id ccd.LUID) string {
	return fmt.Sprintf("%08X:%08X", id.HighPart, id.LowPart)
}

func ValidateProfilePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("profile path is empty")
	}
	if strings.HasSuffix(strings.TrimSpace(path), string(os.PathSeparator)) {
		return errors.New("profile path must be a file")
	}
	return nil
}

func ResolveProfilePath(input string, createDir bool) (string, error) {
	if err := ValidateProfilePath(input); err != nil {
		return "", err
	}

	cleaned := strings.TrimSpace(input)
	if filepath.Ext(cleaned) == "" {
		cleaned += ".monitorprofile"
	}
	if isExplicitPath(cleaned) {
		return cleaned, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}

	profileDir := filepath.Join(homeDir, "Monitor Profiles")
	if createDir {
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			return "", fmt.Errorf("create profile dir: %w", err)
		}
	}

	return filepath.Join(profileDir, cleaned), nil
}

func isExplicitPath(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}
	if filepath.VolumeName(path) != "" {
		return true
	}
	if strings.HasPrefix(path, `\\`) {
		return true
	}
	if strings.Contains(path, string(os.PathSeparator)) {
		return true
	}
	return false
}

func queryFlagsForProfile(activeOnly bool, virtualAware bool) ccd.QueryDisplayFlags {
	if activeOnly {
		if virtualAware {
			return ccd.QueryDisplayFlagsOnlyActivePaths | ccd.QueryDisplayFlagsVirtualModeAware
		}
		return ccd.QueryDisplayFlagsOnlyActivePaths
	}
	if virtualAware {
		return ccd.QueryDisplayFlagsAllPaths | ccd.QueryDisplayFlagsVirtualModeAware
	}
	return ccd.QueryDisplayFlagsAllPaths
}

func profileHasVirtualDisplay(prof profile.Profile) bool {
	for _, mode := range prof.ModeInfo {
		if mode.InfoType == uint32(ccd.DisplayConfigModeInfoTypeDesktopImage) {
			return true
		}
		if mode.TargetMode != nil && mode.TargetMode.TargetVideoSignalInfo.VideoStandard == 65791 {
			return true
		}
	}
	for _, path := range prof.PathInfo {
		if path.Flags&uint32(ccd.DisplayConfigFlagPathSupportVirtualMode) != 0 {
			return true
		}
	}
	for _, info := range prof.AdditionalInfo {
		name := strings.ToLower(info.MonitorFriendlyDevice)
		if strings.Contains(name, "vdd") || strings.Contains(name, "virtual") {
			return true
		}
	}
	return false
}

func ensureDesktopImageModes(paths *[]ccd.DisplayConfigPathInfo, modes *[]ccd.DisplayConfigModeInfo, currentModes []ccd.DisplayConfigModeInfo) bool {
	if len(*modes) == 0 || len(currentModes) == 0 {
		return false
	}

	hasDesktop := false
	for _, mode := range *modes {
		if mode.InfoType == ccd.DisplayConfigModeInfoTypeDesktopImage {
			hasDesktop = true
			break
		}
	}
	if hasDesktop {
		return false
	}

	currentDesktopByID := make(map[uint32]ccd.DisplayConfigModeInfo)
	for _, mode := range currentModes {
		if mode.InfoType == ccd.DisplayConfigModeInfoTypeDesktopImage {
			currentDesktopByID[mode.ID] = mode
		}
	}
	if len(currentDesktopByID) == 0 {
		return false
	}

	targetIndexByID := make(map[uint32]int)
	for i, mode := range *modes {
		if mode.InfoType == ccd.DisplayConfigModeInfoTypeTarget {
			targetIndexByID[mode.ID] = i
		}
	}

	changed := false
	for i := range *paths {
		targetID := (*paths)[i].TargetInfo.ID
		desktopMode, ok := currentDesktopByID[targetID]
		if !ok {
			continue
		}
		targetIdx, ok := targetIndexByID[targetID]
		if !ok {
			continue
		}
		desktopIdx := len(*modes)
		*modes = append(*modes, desktopMode)
		(*paths)[i].Flags |= uint32(ccd.DisplayConfigFlagPathSupportVirtualMode)
		(*paths)[i].TargetInfo.ModeInfoIdx = packTargetModeIndices(targetIdx, desktopIdx)
		changed = true
	}

	return changed
}

func packTargetModeIndices(targetIdx int, desktopIdx int) uint32 {
	return uint32(uint32(uint16(desktopIdx)) | (uint32(uint16(targetIdx)) << 16))
}

type modeKey struct {
	infoType ccd.DisplayConfigModeInfoType
	id       uint32
}

type pathKey struct {
	sourceID uint32
	targetID uint32
}

func mergeProfileWithCurrent(profilePaths []ccd.DisplayConfigPathInfo, profileModes []ccd.DisplayConfigModeInfo, currentPaths []ccd.DisplayConfigPathInfo, currentModes []ccd.DisplayConfigModeInfo) ([]ccd.DisplayConfigPathInfo, []ccd.DisplayConfigModeInfo, bool) {
	if len(currentPaths) == 0 || len(currentModes) == 0 {
		return nil, nil, false
	}

	paths := append([]ccd.DisplayConfigPathInfo(nil), currentPaths...)
	modes := append([]ccd.DisplayConfigModeInfo(nil), currentModes...)

	profileModeMap := make(map[modeKey]ccd.DisplayConfigModeInfo, len(profileModes))
	for _, mode := range profileModes {
		profileModeMap[modeKey{infoType: mode.InfoType, id: mode.ID}] = mode
	}

	for i := range modes {
		if mode, ok := profileModeMap[modeKey{infoType: modes[i].InfoType, id: modes[i].ID}]; ok {
			switch modes[i].InfoType {
			case ccd.DisplayConfigModeInfoTypeTarget:
				modes[i].SetTargetMode(*mode.TargetMode())
			case ccd.DisplayConfigModeInfoTypeSource:
				modes[i].SetSourceMode(*mode.SourceMode())
			case ccd.DisplayConfigModeInfoTypeDesktopImage:
				modes[i].SetDesktopImageInfo(*mode.DesktopImageInfo())
			}
		}
	}

	profilePathMap := make(map[pathKey]ccd.DisplayConfigPathInfo, len(profilePaths))
	for _, path := range profilePaths {
		profilePathMap[pathKey{sourceID: path.SourceInfo.ID, targetID: path.TargetInfo.ID}] = path
	}

	for i := range paths {
		if path, ok := profilePathMap[pathKey{sourceID: paths[i].SourceInfo.ID, targetID: paths[i].TargetInfo.ID}]; ok {
			paths[i].TargetInfo.OutputTechnology = path.TargetInfo.OutputTechnology
			paths[i].TargetInfo.Rotation = path.TargetInfo.Rotation
			paths[i].TargetInfo.Scaling = path.TargetInfo.Scaling
			paths[i].TargetInfo.RefreshRate = path.TargetInfo.RefreshRate
			paths[i].TargetInfo.ScanLineOrdering = path.TargetInfo.ScanLineOrdering
		}
	}

	return paths, modes, true
}
