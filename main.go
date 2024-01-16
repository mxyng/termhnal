package main

import (
	"log"
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	bbt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateTop state = iota
	stateNew
	stateBest
	stateAsk
	stateShow
	stateJob
	stateStory
)

func (s state) String() string {
	switch s {
	case stateTop:
		return "Top"
	case stateNew:
		return "New"
	case stateBest:
		return "Best"
	case stateAsk:
		return "Ask"
	case stateShow:
		return "Show"
	case stateJob:
		return "Job"
	default:
		return ""
	}
}

type model struct {
	hn *HN

	current, previous state
	*Story

	list      list.Model
	listStyle lipgloss.Style

	viewport      viewport.Model
	viewportStyle lipgloss.Style
}

func (m model) Init() bbt.Cmd {
	return m.fetchStories()
}

func (m model) fetchStories() bbt.Cmd {
	var fn func() ([]int, error)
	switch m.current {
	case stateTop:
		fn = m.hn.Top
	case stateNew:
		fn = m.hn.New
	case stateBest:
		fn = m.hn.Best
	case stateAsk:
		fn = m.hn.Ask
	case stateShow:
		fn = m.hn.Show
	case stateJob:
		fn = m.hn.Job
	}

	ids, err := fn()
	if err != nil {
		return nil
	}

	var cmds []bbt.Cmd
	for i := range ids {
		i := i
		cmds = append(cmds, func() bbt.Msg {
			story, err := m.hn.Story(i, ids[i])
			if err != nil {
				return err
			}

			return story
		})
	}

	return bbt.Batch(cmds...)
}

func (m model) fetchComments(parent *Item) bbt.Cmd {
	var cmds []bbt.Cmd
	for i := range parent.Kids {
		i := i
		cmds = append(cmds, func() bbt.Msg {
			comment, err := m.hn.Comment(i, parent.Kids[i])
			if err != nil {
				return err
			}

			parent.Comments = append(parent.Comments, comment)
			return comment
		})
	}

	return bbt.Batch(cmds...)
}

func (m model) Update(msg bbt.Msg) (bbt.Model, bbt.Cmd) {
	switch msg := msg.(type) {
	case bbt.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// noop; mask ctrl+c so it doesn't exit
			return m, nil
		case "ctrl+d":
			return m, bbt.Quit
		case "f5":
			if m.current < stateStory {
				m.list.SetItems([]list.Item{})
				return m, m.fetchStories()
			}
		case "enter":
			if m.current < stateStory {
				m.previous = m.current
				m.current = stateStory

				m.Story = m.list.SelectedItem().(*Story)
				return m, m.fetchComments(m.Story.Item)
			}
		case "esc", "q":
			if m.current == stateStory {
				m.current = m.previous
				return m, nil
			}
		case "1", "2", "3", "4", "5", "6":
			if m.current < stateStory {
				m.previous = m.current
				s, _ := strconv.Atoi(msg.String())

				if int(m.current) != s-1 {
					m.current = state(s - 1)
					m.list.SetItems([]list.Item{})
					return m, m.fetchStories()
				}
			}
		case "g", "home":
			if m.current == stateStory {
				m.viewport.GotoTop()
			}
		case "G", "end":
			if m.current == stateStory {
				m.viewport.GotoBottom()
			}
		case "h", "left":
			if m.current == stateStory {
				m.viewport.ViewUp()
			}
		case "l", "right":
			if m.current == stateStory {
				m.viewport.ViewDown()
			}
		}
	case bbt.WindowSizeMsg:
		h, v := m.listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		h, v = m.viewportStyle.GetFrameSize()
		m.viewport.Width = msg.Width - h
		m.viewport.Height = msg.Height - v
	case *Story:
		items := m.list.Items()
		items = append(items, msg)

		sort.Slice(items, func(i, j int) bool {
			return items[i].(*Story).Rank() < items[j].(*Story).Rank()
		})

		return m, m.list.SetItems(items)
	case *Comment:
		m.viewport.SetContent(m.Story.Format(m.viewport.Width))
		return m, m.fetchComments(msg.Item)
	}

	var cmd bbt.Cmd
	switch m.current {
	case stateStory:
		m.viewport, cmd = m.viewport.Update(msg)
	default:
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	m.list.Title = m.current.String()

	var view string
	switch m.current {
	case stateStory:
		m.viewport.SetContent(m.Story.Format(m.viewport.Width))
		view = m.viewportStyle.Render(m.viewport.View())
	default:
		view = m.listStyle.Render(m.list.View())
	}

	return view
}

func main() {
	p := bbt.NewProgram(model{
		hn: NewHN(),

		list:      list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		listStyle: lipgloss.NewStyle().Margin(1, 2),

		viewport:      viewport.New(0, 0),
		viewportStyle: lipgloss.NewStyle().Margin(1, 2),
	})

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
