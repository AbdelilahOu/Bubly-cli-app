package types

type StatusMsg string

type CheckYtdlpMsg struct {
	Installed bool
}

type CheckFfmpegMsg struct {
	Installed bool
}

type YtdlpInstalledMsg struct {
	Err error
}

type FfmpegInstalledMsg struct {
	Err error
}
