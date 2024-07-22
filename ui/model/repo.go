package model

import (
	"encoding/json"
	"fmt"
	"gitoday/service"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enescakir/emoji"
	"github.com/pkg/errors"
)

type repoModel struct {
	repoList      list.Model
	repoDetail    viewport.Model
	keyMap        list.KeyMap
	repoListItems []*repoItem
	mapAiChannel  map[string]chan *service.ChatResponse
}

func (m repoModel) Init() tea.Cmd {
	return nil
}

func (m *repoModel) tearDown() {
	for _, v := range m.mapAiChannel {
		close(v)
	}
}

func (m *repoModel) updateSize() {
	m.repoList.SetHeight(getRepoListHeight())
	m.repoList.SetWidth(getRepoListWidth())
	m.repoDetail.Width = getRepoDetailWidth()
	m.repoDetail.Height = getRepoDetailHeight()
}
func (m repoModel) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("repo model update", slog.String("msg", fmt.Sprintf("%T %v", tmsg, tmsg)))
	switch msg := tmsg.(type) {
	case tea.WindowSizeMsg:
		setTerminalSize(msg.Width, msg.Height)
		m.updateSize()
		return m, nil
	case MsgQuitRepoView:
		m.tearDown()
		return m, EventRestart()
	case MsgAIFinish:
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.CursorUp):
			m.repoList.CursorUp()
			return show(&m)
		case key.Matches(msg, m.keyMap.CursorDown):
			m.repoList.CursorDown()
			return show(&m)
		}
		switch msg.String() {
		case "enter":
			selected := m.repoList.SelectedItem()
			var r repoItem
			err := json.Unmarshal([]byte(selected.FilterValue()), &r)
			if err != nil {
				slog.Error("json unmarshal error when press enter",
					slog.String("error", fmt.Sprintf("%T %v", errors.Cause(err), errors.Cause(err))))
				return m, nil
			}
			if r.AIProcess == Failed || r.AIProcess == Ready {
				r.AIProcess = InProgress
				go askAI(r.Url, m.mapAiChannel[r.Url])
				m.repoDetail.SetContent(getRepoDetailContent(r))
				return m, m.repoList.SetItem(m.repoList.Index(), r)
			}
			return m, nil
		case "ctrl+c", "esc", "q":
			// Exit the program
			return m, EventQuitRepoView()
		}
	}
	return m, nil
}

func (m repoModel) View() string {
	repoListView := m.repoList.View()
	detailView := m.repoDetail.View()
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		repoListStyle.Render(repoListView),
		detailView,
	)
	return lipgloss.JoinVertical(lipgloss.Left, content)
}

func newRepoModel(repos []*service.Repo) repoModel {
	jobItems := make([]list.Item, len(repos))
	r := makeRepoItem(repos)
	for i, repo := range r {
		jobItems[i] = repo
	}

	l := list.New(jobItems, newAppItemDelegate(), getRepoListWidth(), getRepoListHeight())

	l.Title = fmt.Sprintf("%v Top Repositories %v", emoji.Rocket, emoji.Rocket)
	mapAiChannel := map[string]chan *service.ChatResponse{}
	for _, r := range repos {
		mapAiChannel[r.Url] = make(chan *service.ChatResponse, 1)
	}
	if len(jobItems) > 0 {
		l.Select(0)
	}
	return repoModel{
		repoList:      l,
		repoListItems: r,
		repoDetail: viewport.Model{
			Width:  getRepoDetailWidth(),
			Height: getRepoDetailHeight(),
		},
		keyMap:       list.DefaultKeyMap(),
		mapAiChannel: mapAiChannel,
	}
}

func makeRepoItem(repo []*service.Repo) []*repoItem {
	items := make([]*repoItem, len(repo))
	for i, r := range repo {
		items[i] = &repoItem{
			Index:     i,
			Name:      r.Name,
			Url:       r.Url,
			Desc:      r.Desc,
			Lang:      r.Lang,
			Star:      r.Star,
			Fork:      r.Fork,
			TodayStar: r.TodayStar,
			AIProcess: Ready,
			AIAnswer:  "",
		}

	}
	return items
}
func show(m *repoModel) (tea.Model, tea.Cmd) {
	selected := m.repoList.SelectedItem()
	if selected != nil {
		var r repoItem
		_ = json.Unmarshal([]byte(selected.FilterValue()), &r)
		var aiAnswer []byte
		if r.AIProcess == InProgress {
			ai, err := getAIDetail(m.mapAiChannel, selected.FilterValue())
			if err != nil {
				slog.Error("get ai detail error,set AIProcess failed",
					slog.String("original error", fmt.Sprintf("%T %v", errors.Cause(err), errors.Cause(err))),
					slog.String("stack", fmt.Sprintf("%+v", err)))
				r.AIProcess = Failed
			} else if ai != nil {
				r.AIProcess = Success
				aiAnswer, _ = json.Marshal(ai)
				r.AIAnswer = string(aiAnswer)
			}
		}
		m.repoDetail.SetContent(getRepoDetailContent(r))
		return m, m.repoList.SetItem(m.repoList.Index(), r)
	}
	return m, nil
}
func getRepoDetailContent(r repoItem) string {
	title := fmt.Sprintf("%v Repository Inspiration %v", emoji.OncomingFist, emoji.OncomingFist)
	name := fmt.Sprintf("%v %s ", emoji.TwoHearts, r.Name)
	url := fmt.Sprintf("%v %s", emoji.Link, r.Url)
	des := fmt.Sprintf("%v %s", emoji.OpenBook, r.Desc)
	var aiAnswer string
	switch r.AIProcess {
	case InProgress:
		aiAnswer = fmt.Sprintf("%v AI is analyzing the project,please waiting...%v", emoji.Robot, emoji.TimerClock)
	case Failed:
		aiAnswer = fmt.Sprintf("%v AI is tired,please press [ENTER] to retry later.", emoji.TiredFace)
	case Success:
		aiAnswer = fmt.Sprintf("%v AI analyse finished %v\n\n%s", emoji.FastDownButton, emoji.FastDownButton, formatAI(r.AIAnswer))
	case Ready:
		aiAnswer = fmt.Sprintf("%v Press [ENTER] to unlock AI Power %v", emoji.Locked, emoji.Robot)
	default:
		aiAnswer = fmt.Sprintf("%v Press [ENTER] to unlock AI Power %v", emoji.Locked, emoji.Robot)
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s\n", title, name, url) + "\n" + des + "\n\n\n" + aiAnswer
}
func getAIDetail(responseChans map[string]chan *service.ChatResponse, data string) (*service.ChatResponse, error) {
	var r repoItem
	err := json.Unmarshal([]byte(data), &r)
	if err != nil {
		err = errors.Wrap(err, "json unmarshal error")
		return nil, err
	}
	responseChan := responseChans[r.Url]
	select {
	case res := <-responseChan:
		if res.Error != nil {
			return nil, res.Error
		} else {
			return res, nil
		}
	default:
		return nil, nil
	}
}
func formatAI(answer string) string {
	a := &service.ChatResponse{}
	json.Unmarshal([]byte(answer), a)

	why := fmt.Sprintf("%v %s\n", emoji.QuestionMark, "WHY")
	for _, v := range a.Why {
		why += wrapText(fmt.Sprintf("%v %s\n", emoji.RedCircle, v), uint(getRepoDetailWidth()-4))
	}
	how := fmt.Sprintf("%v %s:\n", emoji.Hammer, "HOW")
	for _, v := range a.How {
		how += wrapText(fmt.Sprintf("%v %s\n", emoji.BlueCircle, v), uint(getRepoDetailWidth()-4))
	}
	others := fmt.Sprintf("%v %s:\n", emoji.BarChart, "MORE")
	for _, v := range a.Other {
		others += wrapText(fmt.Sprintf("%v %s\n", emoji.GreenCircle, v), uint(getRepoDetailWidth()-4))
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s", why, how, others)
}
