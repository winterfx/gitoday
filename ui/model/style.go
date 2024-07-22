package model

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/charmbracelet/lipgloss"
	tsize "github.com/kopoli/go-terminal-size"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	listColor           = "#fe8019"
	listPaneBorderColor = "#3c3836"
	inactivePaneColor   = "#928374"
)

var (
	terminalWidth  atomic.Int32
	terminalHeight atomic.Int32
)

func getRepoListWidth() int {
	return int(terminalWidth.Load()) / 3
}
func getRepoListHeight() int {
	return int(terminalHeight.Load()) * 5 / 6
}
func getRepoDetailWidth() int {
	return int(terminalWidth.Load()) * 2 / 3
}
func getRepoDetailHeight() int {
	return int(terminalHeight.Load()) * 5 / 6
}
func setTerminalSize(w, h int) {
	terminalWidth.Store(int32(w))
	terminalHeight.Store(int32(h))
}
func init() {
	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}
	setTerminalSize(s.Width, s.Height)
}

var (
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	baseStyle     = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#282828"))

	baseListStyle = lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingLeft(1).PaddingBottom(1)

	repoListStyle = baseListStyle.
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color(listPaneBorderColor))

	msgValueVPStyle = baseListStyle.Width(150).PaddingLeft(3)

	modeStyle = baseStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26"))

	repoDetailTitleStyle = baseStyle.
				Bold(true).
				Background(lipgloss.Color(inactivePaneColor)).
				Align(lipgloss.Left)
)

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[*] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
func RightPadTrim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
