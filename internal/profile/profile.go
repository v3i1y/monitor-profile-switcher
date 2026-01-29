package profile

import (
	"encoding/json"
	"fmt"
	"os"
)

type Profile struct {
	PathInfo       []PathInfo       `json:"pathInfo"`
	ModeInfo       []ModeInfo       `json:"modeInfo"`
	AdditionalInfo []AdditionalInfo `json:"additionalInfo"`
}

type LUID struct {
	LowPart  uint32 `json:"lowPart"`
	HighPart uint32 `json:"highPart"`
}

type PathInfo struct {
	SourceInfo PathSourceInfo `json:"sourceInfo"`
	TargetInfo PathTargetInfo `json:"targetInfo"`
	Flags      uint32         `json:"flags"`
}

type PathSourceInfo struct {
	AdapterID   LUID   `json:"adapterId"`
	ID          uint32 `json:"id"`
	ModeInfoIdx uint32 `json:"modeInfoIdx"`
	StatusFlags uint32 `json:"statusFlags"`
}

type PathTargetInfo struct {
	AdapterID        LUID     `json:"adapterId"`
	ID               uint32   `json:"id"`
	ModeInfoIdx      uint32   `json:"modeInfoIdx"`
	OutputTechnology uint32   `json:"outputTechnology"`
	Rotation         uint32   `json:"rotation"`
	Scaling          uint32   `json:"scaling"`
	RefreshRate      Rational `json:"refreshRate"`
	ScanLineOrdering uint32   `json:"scanLineOrdering"`
	TargetAvailable  bool     `json:"targetAvailable"`
	StatusFlags      uint32   `json:"statusFlags"`
}

type ModeInfo struct {
	InfoType         uint32            `json:"infoType"`
	ID               uint32            `json:"id"`
	AdapterID        LUID              `json:"adapterId"`
	TargetMode       *TargetMode       `json:"targetMode,omitempty"`
	SourceMode       *SourceMode       `json:"sourceMode,omitempty"`
	DesktopImageInfo *DesktopImageInfo `json:"desktopImageInfo,omitempty"`
}

type TargetMode struct {
	TargetVideoSignalInfo VideoSignalInfo `json:"targetVideoSignalInfo"`
}

type SourceMode struct {
	Width       uint32 `json:"width"`
	Height      uint32 `json:"height"`
	PixelFormat uint32 `json:"pixelFormat"`
	Position    PointL `json:"position"`
}

type VideoSignalInfo struct {
	PixelRate        int64    `json:"pixelRate"`
	HSyncFreq        Rational `json:"hSyncFreq"`
	VSyncFreq        Rational `json:"vSyncFreq"`
	ActiveSize       Region   `json:"activeSize"`
	TotalSize        Region   `json:"totalSize"`
	VideoStandard    uint32   `json:"videoStandard"`
	ScanLineOrdering uint32   `json:"scanLineOrdering"`
}

type Region struct {
	Cx uint32 `json:"cx"`
	Cy uint32 `json:"cy"`
}

type Rational struct {
	Numerator   uint32 `json:"numerator"`
	Denominator uint32 `json:"denominator"`
}

type PointL struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type RectL struct {
	Left   int32 `json:"left"`
	Top    int32 `json:"top"`
	Right  int32 `json:"right"`
	Bottom int32 `json:"bottom"`
}

type DesktopImageInfo struct {
	PathSourceSize     PointL `json:"pathSourceSize"`
	DesktopImageRegion RectL  `json:"desktopImageRegion"`
	DesktopImageClip   RectL  `json:"desktopImageClip"`
}

type AdditionalInfo struct {
	ManufactureID         uint16 `json:"manufactureId"`
	ProductCodeID         uint16 `json:"productCodeId"`
	Valid                 bool   `json:"valid"`
	MonitorDevicePath     string `json:"monitorDevicePath"`
	MonitorFriendlyDevice string `json:"monitorFriendlyDevice"`
}

func Load(path string) (Profile, error) {
	var profile Profile
	data, err := os.ReadFile(path)
	if err != nil {
		return profile, fmt.Errorf("read profile: %w", err)
	}
	if err := json.Unmarshal(data, &profile); err != nil {
		return profile, fmt.Errorf("parse profile: %w", err)
	}
	return profile, nil
}

func Save(path string, profile Profile) error {
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize profile: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write profile: %w", err)
	}
	return nil
}
