package model

import (
	"fmt"
	"gitoday/service"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enescakir/emoji"
	"github.com/fogleman/ease"
	"github.com/pkg/errors"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

func newFetchModel() tea.Model {
	return fetchModel{0, false, 30, 0, false,
		0, 0, false, false, nil, make(chan error, 1), make(chan []*service.Repo, 1)}
}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}
func crawl(l int, crawlChannel chan []*service.Repo, errorChannel chan error) {
	slog.Debug("crawl start", slog.String("language", string(codeLanguage[l])))
	res, err := service.Crawl(codeLanguage[l])
	if err != nil {
		slog.Error("crawl error", slog.String("language", string(codeLanguage[l])),
			slog.String("original error:", fmt.Sprintf("%T %V", errors.Cause(err), errors.Cause(err))),
			slog.String("stack", fmt.Sprintf("%+v", err)))
		errorChannel <- err
		return
	}
	crawlChannel <- res
	slog.Debug("crawl success", slog.String("language", string(codeLanguage[l])))
}

type fetchModel struct {
	choice       int
	chosen       bool
	ticks        int
	frames       int
	crawling     bool
	resultCount  int
	progress     float64
	loaded       bool
	quitting     bool
	error        error
	errorChannel chan error
	crawlChannel chan []*service.Repo
}

func (m fetchModel) Init() tea.Cmd {
	return tick()
}
func (m fetchModel) TearDown() (tea.Model, tea.Cmd) {
	m.quitting = true
	close(m.crawlChannel)
	return m, tea.Quit
}

// Main update function.
func (m fetchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m.TearDown()
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.chosen {
		return updatechoices(msg, m)
	}
	if !m.crawling {
		go crawl(m.choice, m.crawlChannel, m.errorChannel)
		m.crawling = true
	}

	return updatechosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m fetchModel) View() string {
	var s string
	if m.quitting {
		return "\n  See you later!\n\n"
	}
	if !m.chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updatechoices(msg tea.Msg, m fetchModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.choice++
			if m.choice > len(codeLanguage)-1 {
				m.choice = len(codeLanguage) - 1
			}
		case "k", "up":
			m.choice--
			if m.choice < 0 {
				m.choice = 0
			}
		case "enter":
			m.chosen = true
			return m, frame()
		}

	case tickMsg:
		if m.ticks == 0 {
			m.quitting = true
			return m, tea.Quit
		}
		m.ticks--
		return m, tick()
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updatechosen(msg tea.Msg, m fetchModel) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case frameMsg:
		if !m.loaded {
			m.frames++
			if m.frames > 80 {
				m.frames = 0
			}
			var res []*service.Repo

			select {
			case res = <-m.crawlChannel:
				m.progress = 1
			case m.error = <-m.errorChannel:
				m.progress = 0
				m.ticks = 10
				return m, tick()
			default:
				m.progress = ease.OutBounce(float64(m.frames) / float64(100))
			}
			if m.progress >= 1 {
				m.progress = 1
				m.loaded = true
				m.resultCount = len(res)
				return m, EventCrawlDone(res)
			}
			return m, frame()
		}

	case tickMsg:
		if m.ticks == 0 {
			m.quitting = true
			return m, tea.Quit
		}
		m.ticks--
		return m, tick()
	}

	return m, nil
}

// Sub-views

// The first view, where you're choosing a task
func choicesView(m fetchModel) string {
	c := m.choice

	tpl := "Which language you want to pick up?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")
	var choices string
	for i, v := range codeLanguage {
		choices += checkbox(string(v), i == c) + "\n"
	}

	return fmt.Sprintf(tpl, choices, ticksStyle.Render(strconv.Itoa(m.ticks)))
}

// The second view, after a task has been chosen
func chosenView(m fetchModel) string {
	var msg string
	label := fmt.Sprintf("%v Crawling most excited %s porject about today in github", emoji.Crocodile, codeLanguage[m.choice])
	if m.loaded {
		label = fmt.Sprintf("Prefetch %d %s projects success,waiting for navigate or press [ENTER]", m.resultCount, codeLanguage[m.choice])
	}
	if m.error != nil {
		label = fmt.Sprintf("Error: %s. \nExiting in %s seconds...", m.error.Error(), ticksStyle.Render(strconv.Itoa(m.ticks)))
	}

	return msg + "\n\n" + label + "\n" + progressbar(m.progress) + "%"
}

func progressbar(percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += ramp[i].Render(progressFullChar)
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}
