package progressbar

import (
	"io"
	"time"
)

type IProgressBar interface {
	RenderBlank() error
	Reset()
	Finish() error
	Add(num int) error
	Add64(num int64) error
	Clear() error
	Describe(description string)
	State() IProgressState
}
type IProgressState interface {
}
type IProgressStateTheme interface {
}
type IProgressBarOption func(p IProgressBar)
type IProgressBarManage interface {
	NewOptions(max int, options ...IProgressBarOption) IProgressBar
	NewOptions64(max int64, options ...IProgressBarOption) IProgressBar
	New(max int) IProgressBar
	New64(max int64) IProgressBar
	OptionSetWidth(s int) IProgressBarOption
	OptionSetTheme(t IProgressStateTheme) IProgressBarOption
	OptionSetWriter(w io.Writer) IProgressBarOption
	OptionSetRenderBlankState(r bool) IProgressBarOption
	OptionSetDescription(description string) IProgressBarOption
	OptionEnableColorCodes(colorCodes bool) IProgressBarOption
	OptionSetBytes(maxBytes int) IProgressBarOption
	OptionSetBytes64(maxBytes int64) IProgressBarOption
	OptionShowCount() IProgressBarOption
	OptionShowIts() IProgressBarOption
	OptionThrottle(duration time.Duration) IProgressBarOption
	OptionClearOnFinish() IProgressBarOption
	OptionOnCompletion(cmpl func()) IProgressBarOption
}
