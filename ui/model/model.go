package model

import (
	"context"
	"fmt"
	"gitoday/service"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
)

type StateView uint

const (
	fetchView StateView = iota + 1
	repoView
)

type MainModel struct {
	activeView    StateView
	languageModel tea.Model
	repoModel     tea.Model
	fetchView     tea.Model
}

// implement the mdoel interface
func (m MainModel) Init() tea.Cmd {
	return m.fetchView.Init()
}

// implement the mdoel interface
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MsgRestart:
		m := NewModel()
		return m, m.Init()
	case MsgCrawlDone:
		m.activeView = repoView
		m.repoModel = newRepoModel(msg.Data)
		return m, m.repoModel.Init()
	default:
		switch m.activeView {
		case repoView:
			model, cmd := m.repoModel.Update(msg)
			m.repoModel = model
			return m, cmd
		case fetchView:
			model, cmd := m.fetchView.Update(msg)
			m.fetchView = model
			return m, cmd
		}
	}

	return m, nil
}

// implement the mdoel interface
func (m MainModel) View() string {
	switch m.activeView {
	case repoView:
		return m.repoModel.View()
	case fetchView:
		return m.fetchView.View()
	}
	return ""

}

// NewModel returns a new model
func NewModel() MainModel {
	return MainModel{
		activeView: fetchView, //languageView,
		fetchView:  newFetchModel(),
	}
}

func askAI(repoUrl string, channel chan *service.ChatResponse) {
	slog.Debug("ask ai", slog.String("repoUrl", repoUrl))
	ai, err := service.Chat(context.Background(), repoUrl)
	if err != nil {
		slog.Error("ask ai error", slog.String("repoUrl", repoUrl),
			slog.String("original error", fmt.Sprintf("%T %V", errors.Cause(err), errors.Cause(err))),
			slog.String("stack", fmt.Sprintf("%+v", err)),
		)
		channel <- ai
		return
	}
	channel <- ai
	slog.Debug("ask ai success", slog.String("repoUrl", repoUrl))
}
