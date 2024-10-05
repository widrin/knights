package progressbar

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ProgressBar is a thread-safe, simple
// progress bar
type ProgressBar struct {
	state  state
	config config

	lock sync.Mutex
}

// State is the basic properties of the bar
type State struct {
	CurrentPercent float64
	CurrentBytes   float64
	MaxBytes       int64
	SecondsSince   float64
	SecondsLeft    float64
	KBsPerSecond   float64
}

type state struct {
	currentNum        int64
	currentPercent    int
	lastPercent       int
	currentSaucerSize int

	lastShown time.Time
	startTime time.Time

	counterTime         time.Time
	counterNumSinceLast int64
	counterLastTenRates []float64

	maxLineWidth int
	currentBytes float64

	finished bool
}

type config struct {
	max                  int64 // max number of the counter
	width                int
	writer               io.Writer
	theme                Theme
	renderWithBlankState bool
	description          string
	// whether the output is expected to contain color codes
	colorCodes bool
	maxBytes   int64
	// show the iterations per second
	showIterationsPerSecond bool
	showIterationsCount     bool

	// minimum time to wait in between updates
	throttleDuration time.Duration

	// clear bar once finished
	clearOnFinish bool

	onCompletion func()
}

// Theme defines the elements of the bar
type Theme struct {
	Saucer        string
	SaucerHead    string
	SaucerPadding string
	BarStart      string
	BarEnd        string
}

type Manage struct {
}

var Default = &Manage{}

// OptionSetWidth sets the width of the bar
func (*Manage) OptionSetWidth(s int) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.width = s
	}
}

// OptionSetTheme sets the elements the bar is constructed of
func (*Manage) OptionSetTheme(t IProgressStateTheme) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.theme = *(t.(*Theme))
	}
}

// OptionSetWriter sets the output writer (defaults to os.StdOut)
func (*Manage) OptionSetWriter(w io.Writer) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.writer = w
	}
}

// OptionSetRenderBlankState sets whether or not to render a 0% bar on construction
func (*Manage) OptionSetRenderBlankState(r bool) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.renderWithBlankState = r
	}
}

// OptionSetDescription sets the description of the bar to render in front of it
func (*Manage) OptionSetDescription(description string) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.description = description
	}
}

// OptionEnableColorCodes enables or disables support for color codes
// using mitchellh/colorstring
func (*Manage) OptionEnableColorCodes(colorCodes bool) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.colorCodes = colorCodes
	}
}

// OptionSetBytes will also print the bytes/second
func (m *Manage) OptionSetBytes(maxBytes int) IProgressBarOption {
	return m.OptionSetBytes64(int64(maxBytes))
}

// OptionSetBytes64 will also print the bytes/second
func (*Manage) OptionSetBytes64(maxBytes int64) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.maxBytes = maxBytes
	}
}

// OptionShowCount will also print current count out of total
func (*Manage) OptionShowCount() IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.showIterationsCount = true
	}
}

// OptionShowIts will also print the iterations/second
func (*Manage) OptionShowIts() IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.showIterationsPerSecond = true
	}
}

// OptionThrottle will wait the specified duration before updating again. The default
// duration is 0 seconds.
func (*Manage) OptionThrottle(duration time.Duration) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.throttleDuration = duration
	}
}

// OptionClearOnFinish will clear the bar once its finished
func (*Manage) OptionClearOnFinish() IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.clearOnFinish = true
	}
}

func (*Manage) OptionOnCompletion(cmpl func()) IProgressBarOption {
	return func(p IProgressBar) {
		pp := p.(*ProgressBar)
		pp.config.onCompletion = cmpl
	}
}

var defaultTheme = Theme{Saucer: "█", SaucerPadding: " ", BarStart: "|", BarEnd: "|"}

// NewOptions constructs a new instance of ProgressBar, with any options you specify
func (m *Manage) NewOptions(max int, options ...IProgressBarOption) IProgressBar {
	return m.NewOptions64(int64(max), options...)
}

// NewOptions64 constructs a new instance of ProgressBar, with any options you specify
func (*Manage) NewOptions64(max int64, options ...IProgressBarOption) IProgressBar {
	b := &ProgressBar{
		state: getBlankState(),
		config: config{
			writer:           os.Stdout,
			theme:            defaultTheme,
			width:            10,
			max:              max,
			throttleDuration: 30 * time.Millisecond,
		},
	}

	for _, o := range options {
		o(b)
	}

	if b.config.renderWithBlankState {
		b.RenderBlank()
	}

	return b
}

func getBlankState() state {
	now := time.Now()
	return state{
		startTime:   now,
		lastShown:   now,
		counterTime: now,
	}
}

// New returns a new ProgressBar
// with the specified maximum
func (m *Manage) New(max int) IProgressBar {
	return m.NewOptions(max)
}

// New64 returns a new ProgressBar
// with the specified maximum
func (m *Manage) New64(max int64) IProgressBar {
	return m.NewOptions64(max)
}

// RenderBlank renders the current bar state, you can use this to render a 0% state
func (p *ProgressBar) RenderBlank() error {
	return p.render()
}

// Reset will reset the clock that is used
// to calculate current time and the time left.
func (p *ProgressBar) Reset() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.state = getBlankState()
}

// Finish will fill the bar to full
func (p *ProgressBar) Finish() error {
	if p == nil {
		return fmt.Errorf("progressbar is nil")
	}
	p.lock.Lock()
	p.state.currentNum = p.config.max
	p.lock.Unlock()
	return p.Add(0)
}

// Add will add the specified amount to the progressbar
func (p *ProgressBar) Add(num int) error {
	return p.Add64(int64(num))
}

// Add64 will add the specified amount to the progressbar
func (p *ProgressBar) Add64(num int64) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.config.max == 0 {
		return errors.New("max must be greater than 0")
	}
	p.state.currentNum += num

	// reset the countdown timer every second to take rolling average
	p.state.counterNumSinceLast += num
	if time.Since(p.state.counterTime).Seconds() > 0.5 && time.Since(p.state.counterTime).Seconds() > 0 {
		p.state.counterLastTenRates = append(p.state.counterLastTenRates, float64(p.state.counterNumSinceLast)/time.Since(p.state.counterTime).Seconds())
		if len(p.state.counterLastTenRates) > 10 {
			p.state.counterLastTenRates = p.state.counterLastTenRates[1:]
		}
		p.state.counterTime = time.Now()
		p.state.counterNumSinceLast = 0
	}

	percent := float64(p.state.currentNum) / float64(p.config.max)
	p.state.currentSaucerSize = int(percent * float64(p.config.width))
	p.state.currentPercent = int(percent * 100)
	updateBar := p.state.currentPercent != p.state.lastPercent && p.state.currentPercent > 0

	p.state.currentBytes = float64(percent) * float64(p.config.maxBytes)
	p.state.lastPercent = p.state.currentPercent
	if p.state.currentNum > p.config.max {
		return errors.New("current number exceeds max")
	}

	// always update if show bytes/second or its/second
	if updateBar || p.config.showIterationsPerSecond || p.config.maxBytes > 0 {
		return p.render()
	}

	return nil
}

// Clear erases the progress bar from the current line
func (p *ProgressBar) Clear() error {
	return clearProgressBar(p.config, p.state)
}

// Describe will change the description shown before the progress, which
// can be changed on the fly (as for a slow running process).
func (p *ProgressBar) Describe(description string) {
	p.config.description = description
}

// render renders the progress bar, updating the maximum
// rendered line width. this function is not thread-safe,
// so it must be called with an acquired lock.
func (p *ProgressBar) render() error {
	// make sure that the rendering is not happening too quickly
	// but always show if the currentNum reaches the max
	if time.Since(p.state.lastShown).Nanoseconds() < p.config.throttleDuration.Nanoseconds() &&
		p.state.currentNum < p.config.max {
		return nil
	}

	// first, clear the existing progress bar
	err := clearProgressBar(p.config, p.state)
	if err != nil {
		return err
	}

	// check if the progress bar is finished
	if !p.state.finished && p.state.currentNum >= p.config.max {
		p.state.finished = true
		if !p.config.clearOnFinish {
			renderProgressBar(p.config, p.state)
		}

		if p.config.onCompletion != nil {
			p.config.onCompletion()
		}
	}
	if p.state.finished {
		return nil
	}

	// then, re-render the current progress bar
	w, err := renderProgressBar(p.config, p.state)
	if err != nil {
		return err
	}

	if w > p.state.maxLineWidth {
		p.state.maxLineWidth = w
	}

	p.state.lastShown = time.Now()

	return nil
}

// State returns the current state
func (p *ProgressBar) State() IProgressState {
	p.lock.Lock()
	defer p.lock.Unlock()
	s := State{}
	s.CurrentPercent = float64(p.state.currentNum) / float64(p.config.max)
	s.CurrentBytes = p.state.currentBytes
	s.MaxBytes = p.config.maxBytes
	s.SecondsSince = time.Since(p.state.startTime).Seconds()
	if p.state.currentNum > 0 {
		s.SecondsLeft = s.SecondsSince / float64(p.state.currentNum) * (float64(p.config.max) - float64(p.state.currentNum))
	}
	s.KBsPerSecond = float64(p.state.currentBytes) / 1000.0 / s.SecondsSince
	return s
}

// regex matching ansi escape codes
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func renderProgressBar(c config, s state) (int, error) {
	averageRate := average(s.counterLastTenRates)
	if len(s.counterLastTenRates) == 0 {
		averageRate = float64(s.currentNum) / time.Since(s.startTime).Seconds()
	}

	var saucer string
	if s.currentSaucerSize > 0 {
		saucer = strings.Repeat(c.theme.Saucer, s.currentSaucerSize-1)
		saucerHead := c.theme.SaucerHead
		if saucerHead == "" || s.currentSaucerSize == c.width {
			// use the saucer for the saucer head if it hasn't been set
			// to preserve backwards compatibility
			saucerHead = c.theme.Saucer
		}
		saucer += saucerHead
	}

	// add on bytes string if max bytes option was set
	kbPerSecond := averageRate / float64(c.max) * float64(c.maxBytes) / 1000.0
	bytesString := ""
	if kbPerSecond > 1000.0 {
		bytesString = fmt.Sprintf("(%2.1f MB/s)", kbPerSecond/1000.0)
	} else if kbPerSecond > 0 {
		bytesString = fmt.Sprintf("(%2.1f kB/s)", kbPerSecond)
	}

	if c.showIterationsPerSecond && !c.showIterationsCount {
		// replace bytesString if used
		bytesString = fmt.Sprintf("(%2.0f it/s)", averageRate)
	} else if !c.showIterationsPerSecond && c.showIterationsCount {
		bytesString = fmt.Sprintf("(%d/%d)", s.currentNum, c.max)
	} else if c.showIterationsPerSecond && c.showIterationsCount {
		bytesString = fmt.Sprintf("(%d/%d, %2.0f it/s)", s.currentNum, c.max, averageRate)
	}

	str := fmt.Sprintf("\r%s%4d%% %s%s%s%s %s [%s/%s]",
		c.description,
		s.currentPercent,
		c.theme.BarStart,
		saucer,
		strings.Repeat(c.theme.SaucerPadding, c.width-s.currentSaucerSize),
		c.theme.BarEnd,
		bytesString,
		(time.Duration(time.Since(s.startTime).Seconds()) * time.Second).String(),
		(time.Duration((1/averageRate)*(float64(c.max)-float64(s.currentNum))) * time.Second).String(), // time left
	)

	//if c.colorCodes {
	//	// convert any color codes in the progress bar into the respective ANSI codes
	//	str = colorstring.Color(str)
	//}

	// the width of the string, if printed to the console
	// does not include the carriage return character
	cleanString := strings.Replace(str, "\r", "", -1)

	if c.colorCodes {
		// the ANSI codes for the colors do not take up space in the console output,
		// so they do not count towards the output string width
		cleanString = ansiRegex.ReplaceAllString(cleanString, "")
	}

	// get the amount of runes in the string instead of the
	// character count of the string, as some runes span multiple characters.
	// see https://stackoverflow.com/a/12668840/2733724
	stringWidth := len([]rune(cleanString))

	return stringWidth, writeString(c, str)
}

func clearProgressBar(c config, s state) error {
	// fill the current line with enough spaces
	// to overwrite the progress bar and jump
	// back to the beginning of the line
	str := fmt.Sprintf("\r%s\r", strings.Repeat(" ", s.maxLineWidth))
	return writeString(c, str)
}

func writeString(c config, str string) error {
	if _, err := io.WriteString(c.writer, str); err != nil {
		return err
	}

	if f, ok := c.writer.(*os.File); ok {
		// ignore any errors in Sync(), as stdout
		// can't be synced on some operating systems
		// like Debian 9 (Stretch)
		f.Sync()
	}

	return nil
}

// Reader is the progressbar io.Reader struct
type Reader struct {
	io.Reader
	bar *ProgressBar
}

// Read will read the data and add the number of bytes to the progressbar
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.bar.Add(n)
	return
}

// Close the reader when it implements io.Closer
func (r *Reader) Close() (err error) {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}

// Write implement io.Writer
func (p *ProgressBar) Write(b []byte) (n int, err error) {
	n = len(b)
	p.Add(n)
	return
}

// Read implement io.Reader
func (p *ProgressBar) Read(b []byte) (n int, err error) {
	n = len(b)
	p.Add(n)
	return
}

func average(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}
