package model

import (
	"gitoday/global"
	"gitoday/service"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgChosenLanguage struct {
	Name global.LanguageType
}

type MsgQuitRepoView struct {
}
type MsgQuitLanguageView struct {
}
type MsgRestart struct {
}
type MsgCrawlDone struct {
	Data []*service.Repo
}
type MsgTriggerAI struct {
	Data *repoItem
}
type MsgAIFinish struct {
}

func EventAIFinish() tea.Cmd {
	return func() tea.Msg {
		return MsgAIFinish{}
	}
}

func EventRestart() tea.Cmd {
	return func() tea.Msg {
		return MsgRestart{}
	}
}
func EventCrawlDone(data []*service.Repo) tea.Cmd {
	return func() tea.Msg {
		return MsgCrawlDone{Data: data}
	}
}

func EventQuitRepoView() tea.Cmd {
	return func() tea.Msg {
		return MsgQuitRepoView{}
	}
}
