package switcher

import (
	"monitor-profile-switcher/internal/ccd"
	"monitor-profile-switcher/internal/profile"
)

func profileFromCCD(paths []ccd.DisplayConfigPathInfo, modes []ccd.DisplayConfigModeInfo, additional []ccd.MonitorAdditionalInfo) profile.Profile {
	result := profile.Profile{
		PathInfo:       make([]profile.PathInfo, len(paths)),
		ModeInfo:       make([]profile.ModeInfo, len(modes)),
		AdditionalInfo: make([]profile.AdditionalInfo, len(additional)),
	}

	for i, path := range paths {
		result.PathInfo[i] = profile.PathInfo{
			SourceInfo: profile.PathSourceInfo{
				AdapterID:   toProfileLUID(path.SourceInfo.AdapterID),
				ID:          path.SourceInfo.ID,
				ModeInfoIdx: path.SourceInfo.ModeInfoIdx,
				StatusFlags: uint32(path.SourceInfo.StatusFlags),
			},
			TargetInfo: profile.PathTargetInfo{
				AdapterID:        toProfileLUID(path.TargetInfo.AdapterID),
				ID:               path.TargetInfo.ID,
				ModeInfoIdx:      path.TargetInfo.ModeInfoIdx,
				OutputTechnology: uint32(path.TargetInfo.OutputTechnology),
				Rotation:         uint32(path.TargetInfo.Rotation),
				Scaling:          uint32(path.TargetInfo.Scaling),
				RefreshRate:      toProfileRational(path.TargetInfo.RefreshRate),
				ScanLineOrdering: uint32(path.TargetInfo.ScanLineOrdering),
				TargetAvailable:  path.TargetInfo.TargetAvailable != 0,
				StatusFlags:      uint32(path.TargetInfo.StatusFlags),
			},
			Flags: path.Flags,
		}
	}

	for i, mode := range modes {
		modeProfile := profile.ModeInfo{
			InfoType:  uint32(mode.InfoType),
			ID:        mode.ID,
			AdapterID: toProfileLUID(mode.AdapterID),
		}
		switch mode.InfoType {
		case ccd.DisplayConfigModeInfoTypeTarget:
			target := mode.TargetMode()
			modeProfile.TargetMode = &profile.TargetMode{
				TargetVideoSignalInfo: toProfileVideoSignalInfo(target.TargetVideoSignalInfo),
			}
		case ccd.DisplayConfigModeInfoTypeSource:
			source := mode.SourceMode()
			modeProfile.SourceMode = &profile.SourceMode{
				Width:       source.Width,
				Height:      source.Height,
				PixelFormat: uint32(source.PixelFormat),
				Position: profile.PointL{
					X: source.Position.X,
					Y: source.Position.Y,
				},
			}
		case ccd.DisplayConfigModeInfoTypeDesktopImage:
			desktop := mode.DesktopImageInfo()
			modeProfile.DesktopImageInfo = &profile.DesktopImageInfo{
				PathSourceSize: profile.PointL{
					X: desktop.PathSourceSize.X,
					Y: desktop.PathSourceSize.Y,
				},
				DesktopImageRegion: profile.RectL{
					Left:   desktop.DesktopImageRegion.Left,
					Top:    desktop.DesktopImageRegion.Top,
					Right:  desktop.DesktopImageRegion.Right,
					Bottom: desktop.DesktopImageRegion.Bottom,
				},
				DesktopImageClip: profile.RectL{
					Left:   desktop.DesktopImageClip.Left,
					Top:    desktop.DesktopImageClip.Top,
					Right:  desktop.DesktopImageClip.Right,
					Bottom: desktop.DesktopImageClip.Bottom,
				},
			}
		}
		result.ModeInfo[i] = modeProfile
	}

	for i, info := range additional {
		result.AdditionalInfo[i] = profile.AdditionalInfo{
			ManufactureID:         info.ManufactureID,
			ProductCodeID:         info.ProductCodeID,
			Valid:                 info.Valid,
			MonitorDevicePath:     info.MonitorDevicePath,
			MonitorFriendlyDevice: info.MonitorFriendlyDevice,
		}
	}

	return result
}

func ccdFromProfile(prof profile.Profile) ([]ccd.DisplayConfigPathInfo, []ccd.DisplayConfigModeInfo, []ccd.MonitorAdditionalInfo) {
	paths := make([]ccd.DisplayConfigPathInfo, len(prof.PathInfo))
	modes := make([]ccd.DisplayConfigModeInfo, len(prof.ModeInfo))
	additional := make([]ccd.MonitorAdditionalInfo, len(prof.AdditionalInfo))

	for i, path := range prof.PathInfo {
		paths[i] = ccd.DisplayConfigPathInfo{
			SourceInfo: ccd.DisplayConfigPathSourceInfo{
				AdapterID:   toCCDLUID(path.SourceInfo.AdapterID),
				ID:          path.SourceInfo.ID,
				ModeInfoIdx: path.SourceInfo.ModeInfoIdx,
				StatusFlags: ccd.DisplayConfigSourceStatus(path.SourceInfo.StatusFlags),
			},
			TargetInfo: ccd.DisplayConfigPathTargetInfo{
				AdapterID:        toCCDLUID(path.TargetInfo.AdapterID),
				ID:               path.TargetInfo.ID,
				ModeInfoIdx:      path.TargetInfo.ModeInfoIdx,
				OutputTechnology: ccd.DisplayConfigVideoOutputTechnology(path.TargetInfo.OutputTechnology),
				Rotation:         ccd.DisplayConfigRotation(path.TargetInfo.Rotation),
				Scaling:          ccd.DisplayConfigScaling(path.TargetInfo.Scaling),
				RefreshRate:      toCCDRational(path.TargetInfo.RefreshRate),
				ScanLineOrdering: ccd.DisplayConfigScanLineOrdering(path.TargetInfo.ScanLineOrdering),
				TargetAvailable:  boolToUint32(path.TargetInfo.TargetAvailable),
				StatusFlags:      ccd.DisplayConfigTargetStatus(path.TargetInfo.StatusFlags),
			},
			Flags: path.Flags,
		}
	}

	for i, mode := range prof.ModeInfo {
		modeInfo := ccd.DisplayConfigModeInfo{
			InfoType:  ccd.DisplayConfigModeInfoType(mode.InfoType),
			ID:        mode.ID,
			AdapterID: toCCDLUID(mode.AdapterID),
		}
		switch modeInfo.InfoType {
		case ccd.DisplayConfigModeInfoTypeTarget:
			if mode.TargetMode != nil {
				modeInfo.SetTargetMode(ccd.DisplayConfigTargetMode{
					TargetVideoSignalInfo: toCCDVideoSignalInfo(mode.TargetMode.TargetVideoSignalInfo),
				})
			}
		case ccd.DisplayConfigModeInfoTypeSource:
			if mode.SourceMode != nil {
				modeInfo.SetSourceMode(ccd.DisplayConfigSourceMode{
					Width:       mode.SourceMode.Width,
					Height:      mode.SourceMode.Height,
					PixelFormat: ccd.DisplayConfigPixelFormat(mode.SourceMode.PixelFormat),
					Position: ccd.PointL{
						X: mode.SourceMode.Position.X,
						Y: mode.SourceMode.Position.Y,
					},
				})
			}
		case ccd.DisplayConfigModeInfoTypeDesktopImage:
			if mode.DesktopImageInfo != nil {
				modeInfo.SetDesktopImageInfo(ccd.DisplayConfigDesktopImageInfo{
					PathSourceSize: ccd.PointL{
						X: mode.DesktopImageInfo.PathSourceSize.X,
						Y: mode.DesktopImageInfo.PathSourceSize.Y,
					},
					DesktopImageRegion: ccd.RectL{
						Left:   mode.DesktopImageInfo.DesktopImageRegion.Left,
						Top:    mode.DesktopImageInfo.DesktopImageRegion.Top,
						Right:  mode.DesktopImageInfo.DesktopImageRegion.Right,
						Bottom: mode.DesktopImageInfo.DesktopImageRegion.Bottom,
					},
					DesktopImageClip: ccd.RectL{
						Left:   mode.DesktopImageInfo.DesktopImageClip.Left,
						Top:    mode.DesktopImageInfo.DesktopImageClip.Top,
						Right:  mode.DesktopImageInfo.DesktopImageClip.Right,
						Bottom: mode.DesktopImageInfo.DesktopImageClip.Bottom,
					},
				})
			}
		}
		modes[i] = modeInfo
	}

	for i, info := range prof.AdditionalInfo {
		additional[i] = ccd.MonitorAdditionalInfo{
			ManufactureID:         info.ManufactureID,
			ProductCodeID:         info.ProductCodeID,
			Valid:                 info.Valid,
			MonitorDevicePath:     info.MonitorDevicePath,
			MonitorFriendlyDevice: info.MonitorFriendlyDevice,
		}
	}

	return paths, modes, additional
}

func toProfileLUID(id ccd.LUID) profile.LUID {
	return profile.LUID{
		LowPart:  id.LowPart,
		HighPart: id.HighPart,
	}
}

func toCCDLUID(id profile.LUID) ccd.LUID {
	return ccd.LUID{
		LowPart:  id.LowPart,
		HighPart: id.HighPart,
	}
}

func toProfileRational(r ccd.DisplayConfigRational) profile.Rational {
	return profile.Rational{
		Numerator:   r.Numerator,
		Denominator: r.Denominator,
	}
}

func toCCDRational(r profile.Rational) ccd.DisplayConfigRational {
	return ccd.DisplayConfigRational{
		Numerator:   r.Numerator,
		Denominator: r.Denominator,
	}
}

func toProfileVideoSignalInfo(info ccd.DisplayConfigVideoSignalInfo) profile.VideoSignalInfo {
	return profile.VideoSignalInfo{
		PixelRate:        info.PixelRate,
		HSyncFreq:        toProfileRational(info.HSyncFreq),
		VSyncFreq:        toProfileRational(info.VSyncFreq),
		ActiveSize:       profile.Region{Cx: info.ActiveSize.Cx, Cy: info.ActiveSize.Cy},
		TotalSize:        profile.Region{Cx: info.TotalSize.Cx, Cy: info.TotalSize.Cy},
		VideoStandard:    uint32(info.VideoStandard),
		ScanLineOrdering: uint32(info.ScanLineOrdering),
	}
}

func toCCDVideoSignalInfo(info profile.VideoSignalInfo) ccd.DisplayConfigVideoSignalInfo {
	return ccd.DisplayConfigVideoSignalInfo{
		PixelRate:        info.PixelRate,
		HSyncFreq:        toCCDRational(info.HSyncFreq),
		VSyncFreq:        toCCDRational(info.VSyncFreq),
		ActiveSize:       ccd.DisplayConfig2DRegion{Cx: info.ActiveSize.Cx, Cy: info.ActiveSize.Cy},
		TotalSize:        ccd.DisplayConfig2DRegion{Cx: info.TotalSize.Cx, Cy: info.TotalSize.Cy},
		VideoStandard:    ccd.D3DkmdtVideoSignalStandard(info.VideoStandard),
		ScanLineOrdering: ccd.DisplayConfigScanLineOrdering(info.ScanLineOrdering),
	}
}

func boolToUint32(value bool) uint32 {
	if value {
		return 1
	}
	return 0
}
