//go:build windows

package ccd

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	errorSuccess = 0
)

type LUID struct {
	LowPart  uint32
	HighPart uint32
}

type DisplayConfigVideoOutputTechnology uint32

const (
	DisplayConfigVideoOutputTechnologyOther           DisplayConfigVideoOutputTechnology = 0xFFFFFFFF
	DisplayConfigVideoOutputTechnologyHd15            DisplayConfigVideoOutputTechnology = 0
	DisplayConfigVideoOutputTechnologySVideo          DisplayConfigVideoOutputTechnology = 1
	DisplayConfigVideoOutputTechnologyCompositeVideo  DisplayConfigVideoOutputTechnology = 2
	DisplayConfigVideoOutputTechnologyComponentVideo  DisplayConfigVideoOutputTechnology = 3
	DisplayConfigVideoOutputTechnologyDvi             DisplayConfigVideoOutputTechnology = 4
	DisplayConfigVideoOutputTechnologyHdmi            DisplayConfigVideoOutputTechnology = 5
	DisplayConfigVideoOutputTechnologyLvds            DisplayConfigVideoOutputTechnology = 6
	DisplayConfigVideoOutputTechnologyDJpn            DisplayConfigVideoOutputTechnology = 8
	DisplayConfigVideoOutputTechnologySdi             DisplayConfigVideoOutputTechnology = 9
	DisplayConfigVideoOutputTechnologyDisplayPortExt  DisplayConfigVideoOutputTechnology = 10
	DisplayConfigVideoOutputTechnologyDisplayPortEmb  DisplayConfigVideoOutputTechnology = 11
	DisplayConfigVideoOutputTechnologyUdiExternal     DisplayConfigVideoOutputTechnology = 12
	DisplayConfigVideoOutputTechnologyUdiEmbedded     DisplayConfigVideoOutputTechnology = 13
	DisplayConfigVideoOutputTechnologySdtvDongle      DisplayConfigVideoOutputTechnology = 14
	DisplayConfigVideoOutputTechnologyMiracast        DisplayConfigVideoOutputTechnology = 15
	DisplayConfigVideoOutputTechnologyIndirectWired   DisplayConfigVideoOutputTechnology = 16
	DisplayConfigVideoOutputTechnologyIndirectVirtual DisplayConfigVideoOutputTechnology = 17
	DisplayConfigVideoOutputTechnologyInternal        DisplayConfigVideoOutputTechnology = 0x80000000
	DisplayConfigVideoOutputTechnologyForceUint32     DisplayConfigVideoOutputTechnology = 0xFFFFFFFF
)

type SdcFlags uint32

const (
	SdcFlagsZero                     SdcFlags = 0
	SdcFlagsTopologyInternal         SdcFlags = 0x00000001
	SdcFlagsTopologyClone            SdcFlags = 0x00000002
	SdcFlagsTopologyExtend           SdcFlags = 0x00000004
	SdcFlagsTopologyExternal         SdcFlags = 0x00000008
	SdcFlagsTopologySupplied         SdcFlags = 0x00000010
	SdcFlagsUseSuppliedDisplayConfig SdcFlags = 0x00000020
	SdcFlagsValidate                 SdcFlags = 0x00000040
	SdcFlagsApply                    SdcFlags = 0x00000080
	SdcFlagsNoOptimization           SdcFlags = 0x00000100
	SdcFlagsSaveToDatabase           SdcFlags = 0x00000200
	SdcFlagsAllowChanges             SdcFlags = 0x00000400
	SdcFlagsPathPersistIfRequired    SdcFlags = 0x00000800
	SdcFlagsForceModeEnumeration     SdcFlags = 0x00001000
	SdcFlagsAllowPathOrderChanges    SdcFlags = 0x00002000
	SdcFlagsVirtualModeAware         SdcFlags = 0x00008000
	SdcFlagsUseDatabaseCurrent       SdcFlags = SdcFlagsTopologyInternal | SdcFlagsTopologyClone | SdcFlagsTopologyExtend | SdcFlagsTopologyExternal
)

type DisplayConfigFlags uint32

const (
	DisplayConfigFlagZero                   DisplayConfigFlags = 0x0
	DisplayConfigFlagPathActive             DisplayConfigFlags = 0x00000001
	DisplayConfigFlagPathPreferredUnscaled  DisplayConfigFlags = 0x00000004
	DisplayConfigFlagPathSupportVirtualMode DisplayConfigFlags = 0x00000008
	DisplayConfigFlagPathValidFlags         DisplayConfigFlags = 0x0000000D
)

type DisplayConfigSourceStatus uint32

const (
	DisplayConfigSourceStatusZero  DisplayConfigSourceStatus = 0x0
	DisplayConfigSourceStatusInUse DisplayConfigSourceStatus = 0x00000001
)

type DisplayConfigTargetStatus uint32

const (
	DisplayConfigTargetStatusZero                     DisplayConfigTargetStatus = 0x0
	DisplayConfigTargetStatusInUse                    DisplayConfigTargetStatus = 0x00000001
	DisplayConfigTargetStatusForcible                 DisplayConfigTargetStatus = 0x00000002
	DisplayConfigTargetStatusForcedAvailabilityBoot   DisplayConfigTargetStatus = 0x00000004
	DisplayConfigTargetStatusForcedAvailabilityPath   DisplayConfigTargetStatus = 0x00000008
	DisplayConfigTargetStatusForcedAvailabilitySystem DisplayConfigTargetStatus = 0x00000010
	DisplayConfigTargetStatusIsHMD                    DisplayConfigTargetStatus = 0x00000020
)

type DisplayConfigRotation uint32

const (
	DisplayConfigRotationZero        DisplayConfigRotation = 0x0
	DisplayConfigRotationIdentity    DisplayConfigRotation = 1
	DisplayConfigRotationRotate90    DisplayConfigRotation = 2
	DisplayConfigRotationRotate180   DisplayConfigRotation = 3
	DisplayConfigRotationRotate270   DisplayConfigRotation = 4
	DisplayConfigRotationForceUint32 DisplayConfigRotation = 0xFFFFFFFF
)

type DisplayConfigPixelFormat uint32

const (
	DisplayConfigPixelFormatZero        DisplayConfigPixelFormat = 0x0
	DisplayConfigPixelFormat8Bpp        DisplayConfigPixelFormat = 1
	DisplayConfigPixelFormat16Bpp       DisplayConfigPixelFormat = 2
	DisplayConfigPixelFormat24Bpp       DisplayConfigPixelFormat = 3
	DisplayConfigPixelFormat32Bpp       DisplayConfigPixelFormat = 4
	DisplayConfigPixelFormatNongdi      DisplayConfigPixelFormat = 5
	DisplayConfigPixelFormatForceUint32 DisplayConfigPixelFormat = 0xFFFFFFFF
)

type DisplayConfigScaling uint32

const (
	DisplayConfigScalingZero                   DisplayConfigScaling = 0x0
	DisplayConfigScalingIdentity               DisplayConfigScaling = 1
	DisplayConfigScalingCentered               DisplayConfigScaling = 2
	DisplayConfigScalingStretched              DisplayConfigScaling = 3
	DisplayConfigScalingAspectRatioCenteredMax DisplayConfigScaling = 4
	DisplayConfigScalingCustom                 DisplayConfigScaling = 5
	DisplayConfigScalingPreferred              DisplayConfigScaling = 128
	DisplayConfigScalingForceUint32            DisplayConfigScaling = 0xFFFFFFFF
)

type DisplayConfigRational struct {
	Numerator   uint32
	Denominator uint32
}

type DisplayConfigScanLineOrdering uint32

const (
	DisplayConfigScanLineOrderingUnspecified               DisplayConfigScanLineOrdering = 0
	DisplayConfigScanLineOrderingProgressive               DisplayConfigScanLineOrdering = 1
	DisplayConfigScanLineOrderingInterlaced                DisplayConfigScanLineOrdering = 2
	DisplayConfigScanLineOrderingInterlacedUpperFieldFirst DisplayConfigScanLineOrdering = DisplayConfigScanLineOrderingInterlaced
	DisplayConfigScanLineOrderingInterlacedLowerFieldFirst DisplayConfigScanLineOrdering = 3
	DisplayConfigScanLineOrderingForceUint32               DisplayConfigScanLineOrdering = 0xFFFFFFFF
)

type DisplayConfigPathInfo struct {
	SourceInfo DisplayConfigPathSourceInfo
	TargetInfo DisplayConfigPathTargetInfo
	Flags      uint32
}

type DisplayConfigModeInfoType uint32

const (
	DisplayConfigModeInfoTypeZero         DisplayConfigModeInfoType = 0
	DisplayConfigModeInfoTypeSource       DisplayConfigModeInfoType = 1
	DisplayConfigModeInfoTypeTarget       DisplayConfigModeInfoType = 2
	DisplayConfigModeInfoTypeDesktopImage DisplayConfigModeInfoType = 3
	DisplayConfigModeInfoTypeForceUint32  DisplayConfigModeInfoType = 0xFFFFFFFF
)

const displayConfigModeInfoUnionSize = 48

type DisplayConfigModeInfo struct {
	InfoType  DisplayConfigModeInfoType
	ID        uint32
	AdapterID LUID
	Mode      [displayConfigModeInfoUnionSize]byte
}

func (m *DisplayConfigModeInfo) TargetMode() *DisplayConfigTargetMode {
	return (*DisplayConfigTargetMode)(unsafe.Pointer(&m.Mode[0]))
}

func (m *DisplayConfigModeInfo) SourceMode() *DisplayConfigSourceMode {
	return (*DisplayConfigSourceMode)(unsafe.Pointer(&m.Mode[0]))
}

func (m *DisplayConfigModeInfo) SetTargetMode(target DisplayConfigTargetMode) {
	*m.TargetMode() = target
}

func (m *DisplayConfigModeInfo) SetSourceMode(source DisplayConfigSourceMode) {
	*m.SourceMode() = source
}

func (m *DisplayConfigModeInfo) DesktopImageInfo() *DisplayConfigDesktopImageInfo {
	return (*DisplayConfigDesktopImageInfo)(unsafe.Pointer(&m.Mode[0]))
}

func (m *DisplayConfigModeInfo) SetDesktopImageInfo(info DisplayConfigDesktopImageInfo) {
	*m.DesktopImageInfo() = info
}

type DisplayConfig2DRegion struct {
	Cx uint32
	Cy uint32
}

type D3DkmdtVideoSignalStandard uint32

const (
	D3DkmdtVideoSignalStandardUninitialized D3DkmdtVideoSignalStandard = 0
	D3DkmdtVideoSignalStandardVesaDmt       D3DkmdtVideoSignalStandard = 1
	D3DkmdtVideoSignalStandardVesaGtf       D3DkmdtVideoSignalStandard = 2
	D3DkmdtVideoSignalStandardVesaCvt       D3DkmdtVideoSignalStandard = 3
	D3DkmdtVideoSignalStandardIbm           D3DkmdtVideoSignalStandard = 4
	D3DkmdtVideoSignalStandardApple         D3DkmdtVideoSignalStandard = 5
	D3DkmdtVideoSignalStandardNtscM         D3DkmdtVideoSignalStandard = 6
	D3DkmdtVideoSignalStandardNtscJ         D3DkmdtVideoSignalStandard = 7
	D3DkmdtVideoSignalStandardNtsc443       D3DkmdtVideoSignalStandard = 8
	D3DkmdtVideoSignalStandardPalB          D3DkmdtVideoSignalStandard = 9
	D3DkmdtVideoSignalStandardPalB1         D3DkmdtVideoSignalStandard = 10
	D3DkmdtVideoSignalStandardPalG          D3DkmdtVideoSignalStandard = 11
	D3DkmdtVideoSignalStandardPalH          D3DkmdtVideoSignalStandard = 12
	D3DkmdtVideoSignalStandardPalI          D3DkmdtVideoSignalStandard = 13
	D3DkmdtVideoSignalStandardPalD          D3DkmdtVideoSignalStandard = 14
	D3DkmdtVideoSignalStandardPalN          D3DkmdtVideoSignalStandard = 15
	D3DkmdtVideoSignalStandardPalNc         D3DkmdtVideoSignalStandard = 16
	D3DkmdtVideoSignalStandardSecamB        D3DkmdtVideoSignalStandard = 17
	D3DkmdtVideoSignalStandardSecamD        D3DkmdtVideoSignalStandard = 18
	D3DkmdtVideoSignalStandardSecamG        D3DkmdtVideoSignalStandard = 19
	D3DkmdtVideoSignalStandardSecamH        D3DkmdtVideoSignalStandard = 20
	D3DkmdtVideoSignalStandardSecamK        D3DkmdtVideoSignalStandard = 21
	D3DkmdtVideoSignalStandardSecamK1       D3DkmdtVideoSignalStandard = 22
	D3DkmdtVideoSignalStandardSecamL        D3DkmdtVideoSignalStandard = 23
	D3DkmdtVideoSignalStandardSecamL1       D3DkmdtVideoSignalStandard = 24
	D3DkmdtVideoSignalStandardEia861        D3DkmdtVideoSignalStandard = 25
	D3DkmdtVideoSignalStandardEia861A       D3DkmdtVideoSignalStandard = 26
	D3DkmdtVideoSignalStandardEia861B       D3DkmdtVideoSignalStandard = 27
	D3DkmdtVideoSignalStandardPalK          D3DkmdtVideoSignalStandard = 28
	D3DkmdtVideoSignalStandardPalK1         D3DkmdtVideoSignalStandard = 29
	D3DkmdtVideoSignalStandardPalL          D3DkmdtVideoSignalStandard = 30
	D3DkmdtVideoSignalStandardPalM          D3DkmdtVideoSignalStandard = 31
	D3DkmdtVideoSignalStandardOther         D3DkmdtVideoSignalStandard = 255
	D3DkmdtVideoSignalStandardUSB           D3DkmdtVideoSignalStandard = 65791
)

type DisplayConfigVideoSignalInfo struct {
	PixelRate        int64
	HSyncFreq        DisplayConfigRational
	VSyncFreq        DisplayConfigRational
	ActiveSize       DisplayConfig2DRegion
	TotalSize        DisplayConfig2DRegion
	VideoStandard    D3DkmdtVideoSignalStandard
	ScanLineOrdering DisplayConfigScanLineOrdering
}

type DisplayConfigTargetMode struct {
	TargetVideoSignalInfo DisplayConfigVideoSignalInfo
}

type PointL struct {
	X int32
	Y int32
}

type RectL struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type DisplayConfigSourceMode struct {
	Width       uint32
	Height      uint32
	PixelFormat DisplayConfigPixelFormat
	Position    PointL
}

type DisplayConfigDesktopImageInfo struct {
	PathSourceSize     PointL
	DesktopImageRegion RectL
	DesktopImageClip   RectL
}

type DisplayConfigPathSourceInfo struct {
	AdapterID   LUID
	ID          uint32
	ModeInfoIdx uint32
	StatusFlags DisplayConfigSourceStatus
}

type DisplayConfigPathTargetInfo struct {
	AdapterID        LUID
	ID               uint32
	ModeInfoIdx      uint32
	OutputTechnology DisplayConfigVideoOutputTechnology
	Rotation         DisplayConfigRotation
	Scaling          DisplayConfigScaling
	RefreshRate      DisplayConfigRational
	ScanLineOrdering DisplayConfigScanLineOrdering
	TargetAvailable  uint32
	StatusFlags      DisplayConfigTargetStatus
}

type QueryDisplayFlags uint32

const (
	QueryDisplayFlagsZero             QueryDisplayFlags = 0x0
	QueryDisplayFlagsAllPaths         QueryDisplayFlags = 0x00000001
	QueryDisplayFlagsOnlyActivePaths  QueryDisplayFlags = 0x00000002
	QueryDisplayFlagsDatabaseCurrent  QueryDisplayFlags = 0x00000004
	QueryDisplayFlagsVirtualModeAware QueryDisplayFlags = 0x00000010
	QueryDisplayFlagsIncludeHMD       QueryDisplayFlags = 0x00000020
)

type DisplayConfigDeviceInfoType uint32

const (
	DisplayConfigDeviceInfoTypeGetSourceName               DisplayConfigDeviceInfoType = 1
	DisplayConfigDeviceInfoTypeGetTargetName               DisplayConfigDeviceInfoType = 2
	DisplayConfigDeviceInfoTypeGetTargetPreferredMode      DisplayConfigDeviceInfoType = 3
	DisplayConfigDeviceInfoTypeGetAdapterName              DisplayConfigDeviceInfoType = 4
	DisplayConfigDeviceInfoTypeSetTargetPersistence        DisplayConfigDeviceInfoType = 5
	DisplayConfigDeviceInfoTypeGetTargetBaseType           DisplayConfigDeviceInfoType = 6
	DisplayConfigDeviceInfoTypeGetSupportVirtualResolution DisplayConfigDeviceInfoType = 7
	DisplayConfigDeviceInfoTypeSetSupportVirtualResolution DisplayConfigDeviceInfoType = 8
	DisplayConfigDeviceInfoTypeAdvancedColorInfo           DisplayConfigDeviceInfoType = 9
	DisplayConfigDeviceInfoTypeAdvancedColorState          DisplayConfigDeviceInfoType = 10
	DisplayConfigDeviceInfoTypeSDRWhiteLevel               DisplayConfigDeviceInfoType = 11
	DisplayConfigDeviceInfoTypeForceUint32                 DisplayConfigDeviceInfoType = 0xFFFFFFFF
)

type DisplayConfigTargetDeviceNameFlags struct {
	Value uint32
}

type DisplayConfigDeviceInfoHeader struct {
	Type      DisplayConfigDeviceInfoType
	Size      uint32
	AdapterID LUID
	ID        uint32
}

type DisplayConfigTargetDeviceName struct {
	Header                    DisplayConfigDeviceInfoHeader
	Flags                     DisplayConfigTargetDeviceNameFlags
	OutputTechnology          DisplayConfigVideoOutputTechnology
	EdidManufactureID         uint16
	EdidProductCodeID         uint16
	ConnectorInstance         uint32
	MonitorFriendlyDeviceName [64]uint16
	MonitorDevicePath         [128]uint16
}

type MonitorAdditionalInfo struct {
	ManufactureID         uint16
	ProductCodeID         uint16
	Valid                 bool
	MonitorDevicePath     string
	MonitorFriendlyDevice string
}

var (
	user32                          = windows.NewLazySystemDLL("user32.dll")
	procSetDisplayConfig            = user32.NewProc("SetDisplayConfig")
	procQueryDisplayConfig          = user32.NewProc("QueryDisplayConfig")
	procGetDisplayConfigBufferSizes = user32.NewProc("GetDisplayConfigBufferSizes")
	procDisplayConfigGetDeviceInfo  = user32.NewProc("DisplayConfigGetDeviceInfo")
)

func SetDisplayConfig(paths []DisplayConfigPathInfo, modes []DisplayConfigModeInfo, flags SdcFlags) error {
	var pathPtr *DisplayConfigPathInfo
	var modePtr *DisplayConfigModeInfo
	if len(paths) > 0 {
		pathPtr = &paths[0]
	}
	if len(modes) > 0 {
		modePtr = &modes[0]
	}

	r1, _, _ := procSetDisplayConfig.Call(
		uintptr(uint32(len(paths))),
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(uint32(len(modes))),
		uintptr(unsafe.Pointer(modePtr)),
		uintptr(flags),
	)
	if r1 != errorSuccess {
		return fmt.Errorf("SetDisplayConfig failed: %d", r1)
	}
	return nil
}

func GetDisplaySettings(activeOnly bool) ([]DisplayConfigPathInfo, []DisplayConfigModeInfo, []MonitorAdditionalInfo, error) {
	flags := QueryDisplayFlagsAllPaths
	if activeOnly {
		flags = QueryDisplayFlagsOnlyActivePaths
	}
	return GetDisplaySettingsWithFlags(flags)
}

func GetDisplaySettingsWithFlags(flags QueryDisplayFlags) ([]DisplayConfigPathInfo, []DisplayConfigModeInfo, []MonitorAdditionalInfo, error) {
	var numPaths uint32
	var numModes uint32
	if err := getDisplayConfigBufferSizes(flags, &numPaths, &numModes); err != nil {
		return nil, nil, nil, err
	}

	pathInfo := make([]DisplayConfigPathInfo, numPaths)
	modeInfo := make([]DisplayConfigModeInfo, numModes)
	if err := queryDisplayConfig(flags, &numPaths, pathInfo, &numModes, modeInfo); err != nil {
		return nil, nil, nil, err
	}
	pathInfo = pathInfo[:numPaths]
	modeInfo = modeInfo[:numModes]

	filteredModes := make([]DisplayConfigModeInfo, 0, len(modeInfo))
	for _, mode := range modeInfo {
		if mode.InfoType != DisplayConfigModeInfoTypeZero {
			filteredModes = append(filteredModes, mode)
		}
	}
	modeInfo = filteredModes

	filteredPaths := make([]DisplayConfigPathInfo, 0, len(pathInfo))
	for _, path := range pathInfo {
		if path.TargetInfo.TargetAvailable != 0 {
			filteredPaths = append(filteredPaths, path)
		}
	}
	pathInfo = filteredPaths

	additional := make([]MonitorAdditionalInfo, len(modeInfo))
	for i := range modeInfo {
		if modeInfo[i].InfoType == DisplayConfigModeInfoTypeTarget {
			info, err := GetMonitorAdditionalInfo(modeInfo[i].AdapterID, modeInfo[i].ID)
			if err == nil {
				additional[i] = info
			} else {
				additional[i] = MonitorAdditionalInfo{Valid: false}
			}
		}
	}

	return pathInfo, modeInfo, additional, nil
}

func GetMonitorAdditionalInfo(adapterID LUID, targetID uint32) (MonitorAdditionalInfo, error) {
	var result MonitorAdditionalInfo

	deviceName := DisplayConfigTargetDeviceName{}
	deviceName.Header.Type = DisplayConfigDeviceInfoTypeGetTargetName
	deviceName.Header.Size = uint32(unsafe.Sizeof(deviceName))
	deviceName.Header.AdapterID = adapterID
	deviceName.Header.ID = targetID

	r1, _, _ := procDisplayConfigGetDeviceInfo.Call(uintptr(unsafe.Pointer(&deviceName)))
	if r1 != errorSuccess {
		return result, fmt.Errorf("DisplayConfigGetDeviceInfo failed: %d", r1)
	}

	result.Valid = true
	result.ManufactureID = deviceName.EdidManufactureID
	result.ProductCodeID = deviceName.EdidProductCodeID
	result.MonitorDevicePath = windows.UTF16ToString(deviceName.MonitorDevicePath[:])
	result.MonitorFriendlyDevice = windows.UTF16ToString(deviceName.MonitorFriendlyDeviceName[:])

	return result, nil
}

func getDisplayConfigBufferSizes(flags QueryDisplayFlags, numPaths *uint32, numModes *uint32) error {
	r1, _, _ := procGetDisplayConfigBufferSizes.Call(
		uintptr(flags),
		uintptr(unsafe.Pointer(numPaths)),
		uintptr(unsafe.Pointer(numModes)),
	)
	if r1 != errorSuccess {
		return fmt.Errorf("GetDisplayConfigBufferSizes failed: %d", r1)
	}
	return nil
}

func queryDisplayConfig(flags QueryDisplayFlags, numPaths *uint32, paths []DisplayConfigPathInfo, numModes *uint32, modes []DisplayConfigModeInfo) error {
	var pathPtr *DisplayConfigPathInfo
	var modePtr *DisplayConfigModeInfo
	if len(paths) > 0 {
		pathPtr = &paths[0]
	}
	if len(modes) > 0 {
		modePtr = &modes[0]
	}

	r1, _, _ := procQueryDisplayConfig.Call(
		uintptr(flags),
		uintptr(unsafe.Pointer(numPaths)),
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(numModes)),
		uintptr(unsafe.Pointer(modePtr)),
		uintptr(0),
	)
	if r1 != errorSuccess {
		return fmt.Errorf("QueryDisplayConfig failed: %d", r1)
	}
	return nil
}
