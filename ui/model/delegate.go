package model

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enescakir/emoji"
	"github.com/mitchellh/go-wordwrap"
)

type AIStatus int

const (
	Ready AIStatus = iota
	InProgress
	Failed
	Success
)

type repoItem struct {
	Index     int      `json:"index"`
	Name      string   `json:"name"`
	Url       string   `json:"url"`
	Desc      string   `json:"desc"`
	Lang      string   `json:"lang"`
	Star      string   `json:"star"`
	Fork      string   `json:"fork"`
	TodayStar string   `json:"todayStar"`
	AIProcess AIStatus `json:"AIProcess"`
	AIAnswer  string   `json:"AIAnswer"`
}

func (r repoItem) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s", r.Name, r.Lang, r.Url, r.Star, r.Fork, r.TodayStar, r.Desc)
}

func (r repoItem) Title() string {
	return fmt.Sprintf("%v %s", emoji.LargeOrangeDiamond, r.Name)
}

func (r repoItem) Description() string {
	lang := fmt.Sprintf("%s%v", r.Lang, emoji.Laptop)
	star := fmt.Sprintf("%s%v", r.Star, emoji.Star)
	fork := fmt.Sprintf("%s%v", r.Fork, emoji.Wrench)
	starToday := fmt.Sprintf("%s%v", r.TodayStar, emoji.Fire)
	s := Trim(r.Desc, getRepoListWidth())
	des := wrapText(s, uint(getRepoListWidth()))
	return fmt.Sprintf("  %s  %s  %s  %s", lang, starToday, fork, star) + "\n" + des
}

func (r repoItem) FilterValue() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}
func newAppItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = true
	d.SetHeight(3)
	//set a cmd to show details
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		return nil

	}
	return d
}
func wrapText(text string, lineWidth uint) string {
	return wordwrap.WrapString(text, lineWidth)
}
