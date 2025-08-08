package types

type StatusMsg string

type CheckYtdlpMsg struct {
	Installed bool
}
type YtdlpInstalledMsg struct {
	Err error
}
