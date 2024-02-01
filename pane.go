package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	bbt "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Pane interface {
	Update(bbt.Msg) (Pane, bbt.Cmd)
	View() string

	Size() (width, height int)
	SetSize(width, height int)
}

type FocusMsg string

func Focus(name string) bbt.Cmd {
	return func() bbt.Msg {
		return FocusMsg(name)
	}
}

type ViewType interface {
	*Story | *Comment
}

type ViewMsg[T ViewType] struct {
	Value T
}

func View[T ViewType](t T) bbt.Cmd {
	return func() bbt.Msg {
		return ViewMsg[T]{
			Value: t,
		}
	}
}

type PaneView struct {
	*Story
	style    lipgloss.Style
	viewport viewport.Model

	content strings.Builder

	styleTitle       lipgloss.Style
	styleDescription lipgloss.Style
	styleComment     lipgloss.Style
}

func NewPaneView() *PaneView {
	return &PaneView{
		style:    lipgloss.NewStyle().Margin(1, 2),
		viewport: viewport.New(0, 0),
		styleTitle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}),
		styleDescription: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#a49fa5", Dark: "#777777"}),
		styleComment: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			Border(lipgloss.NormalBorder(), false).
			BorderLeft(true).
			PaddingLeft(1).
			MarginTop(1),
	}
}

func (p *PaneView) Update(msg bbt.Msg) (Pane, bbt.Cmd) {
	comments := func(parent *Item) []bbt.Cmd {
		hn := NewHN()
		var cmds []bbt.Cmd
		for i := range parent.Kids {
			i := i
			cmds = append(cmds, func() bbt.Msg {
				comment, err := hn.Comment(i, parent.Kids[i])
				if err != nil {
					return err
				}

				parent.Comments = append(parent.Comments, comment)
				sort.Slice(parent.Comments, func(i, j int) bool {
					return parent.Comments[i].Rank < parent.Comments[j].Rank
				})

				return ViewMsg[*Comment]{
					Value: comment,
				}
			})
		}

		return cmds
	}

	switch msg := msg.(type) {
	case ViewMsg[*Story]:
		p.Story = msg.Value
		p.Render()
		return p, bbt.Batch(comments(msg.Value.Item)...)
	case ViewMsg[*Comment]:
		p.Render()
		return p, bbt.Batch(comments(msg.Value.Item)...)
	case bbt.KeyMsg:
		switch msg.String() {
		case "g", "home":
			p.viewport.GotoTop()
		case "G", "end":
			p.viewport.GotoBottom()
		case "q", "esc":
			p.viewport.GotoTop()
			return p, Focus("list")
		}
	case bbt.WindowSizeMsg:
		p.Render()
	}

	var cmd bbt.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

func (p *PaneView) View() string {
	return p.style.Render(p.viewport.View())
}

func (p *PaneView) Render() {
	p.content.Reset()
	if s := p.Story; s != nil {
		title := strings.TrimPrefix(s.Title(), fmt.Sprintf("%d. ", s.Rank+1))
		fmt.Fprintln(&p.content, p.styleTitle.Render(title))

		description := strings.TrimSpace(s.Description())
		fmt.Fprintln(&p.content, p.styleDescription.Render(description))

		if s.URL != "" {
			fmt.Fprintln(&p.content, p.styleDescription.Copy().Underline(true).Italic(true).Render(s.URL))
		} else if s.Text != "" {
			fmt.Fprintln(&p.content)
			fmt.Fprintln(&p.content, p.styleDescription.Copy().Width(p.style.GetWidth()).Render(HTMLText(s.Text)))
		}

		h, _ := p.styleComment.GetFrameSize()

		var view func(lipgloss.Style, []*Comment) string
		view = func(style lipgloss.Style, comments []*Comment) string {
			var lines []string
			for _, comment := range comments {
				if comment.By != "" {
					var sb strings.Builder
					fmt.Fprintln(&sb, comment.Title())
					fmt.Fprint(&sb, comment.Description())

					if len(comment.Comments) > 0 {
						fmt.Fprintln(&sb)
						fmt.Fprint(&sb, view(style.Copy().Width(style.GetWidth()-h), comment.Comments))
					}

					lines = append(lines, style.Render(sb.String()))
				}
			}

			return strings.Join(lines, "\n")
		}

		fmt.Fprintln(&p.content, view(p.styleComment.Copy().Width(p.style.GetWidth()-h), s.Comments))
	}

	p.viewport.SetContent(p.content.String())
}

func (p *PaneView) Size() (width, height int) {
	h, v := p.style.GetFrameSize()
	return p.style.GetWidth() + h, p.style.GetHeight() + v
}

func (p *PaneView) SetSize(width, height int) {
	h, v := p.style.GetFrameSize()
	p.style = p.style.Width(width - h).Height(height - v)
	p.viewport.Width, p.viewport.Height = width-h, height-v
}

type ListType interface {
	string | *Story
}

type ListMsg[T ListType] struct {
	Value T
}

func List[T ListType](t T) bbt.Cmd {
	return func() bbt.Msg {
		return ListMsg[T]{
			Value: t,
		}
	}
}

type PaneList struct {
	model list.Model
	style lipgloss.Style
}

func NewPaneList() *PaneList {
	color := lipgloss.Color("#ff6600")
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(color).BorderLeftForeground(color)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy().Faint(true)

	model := list.New([]list.Item{}, delegate, 0, 0)
	model.SetShowHelp(false)
	model.SetShowStatusBar(false)
	model.SetShowTitle(false)
	return &PaneList{
		model: model,
		style: lipgloss.NewStyle().Margin(1, 2),
	}
}

func (p *PaneList) Update(msg bbt.Msg) (Pane, bbt.Cmd) {
	hn := NewHN()
	switch msg := msg.(type) {
	case ListMsg[string]:
		var fn func() ([]int, error)
		switch strings.ToLower(msg.Value) {
		case "top":
			fn = hn.Top
		case "new":
			fn = hn.New
		case "best":
			fn = hn.Best
		case "ask":
			fn = hn.Ask
		case "show":
			fn = hn.Show
		case "job":
			fn = hn.Job
		case "clear":
			p.model.Select(0)
			return p, p.model.SetItems([]list.Item{})
		default:
			return p, nil
		}

		ids, err := fn()
		if err != nil {
			return p, nil
		}

		var cmds []bbt.Cmd
		for i := range ids {
			i := i
			cmds = append(cmds, func() bbt.Msg {
				story, err := hn.Story(i, ids[i])
				if err != nil {
					return err
				}

				return ListMsg[*Story]{
					Value: story,
				}
			})
		}

		return p, bbt.Batch(cmds...)
	case ListMsg[*Story]:
		items := append(p.model.Items(), msg.Value)
		sort.Slice(items, func(i, j int) bool {
			return items[i].(*Story).Rank < items[j].(*Story).Rank
		})

		return p, p.model.SetItems(items)
	case bbt.KeyMsg:
		switch msg.String() {
		case "enter":
			story := p.model.SelectedItem().(*Story)
			return p, bbt.Sequence(
				Focus("view"),
				View(story),
			)
		case "k", "up":
			if p.model.Index() == 0 {
				return p, Focus("header")
			}
		case "tab":
			return p, Focus("header")
		}
	}

	var cmd bbt.Cmd
	p.model, cmd = p.model.Update(msg)
	return p, cmd
}

func (p *PaneList) View() string {
	return p.style.Render(p.model.View())
}

func (p *PaneList) Size() (width, height int) {
	h, v := p.style.GetFrameSize()
	return p.model.Width() + h, p.model.Height() + v
}

func (p *PaneList) SetSize(width, height int) {
	h, v := p.style.GetFrameSize()
	p.model.SetSize(width-h, height-v)
}

type HeaderMsg int

func Header(n int) bbt.Cmd {
	return func() bbt.Msg {
		return HeaderMsg(n)
	}
}

type PaneHeader struct {
	current int
	states  []lipgloss.Style
	style   lipgloss.Style
}

func NewPaneHeader(states ...string) *PaneHeader {
	pane := PaneHeader{
		style: lipgloss.NewStyle().Margin(1, 2, 0),
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Margin(0, 2, 1, 0)

	for _, state := range states {
		pane.states = append(pane.states, style.Copy().SetString(state))
	}

	return &pane
}

func (p *PaneHeader) Update(msg bbt.Msg) (Pane, bbt.Cmd) {
	switch msg := msg.(type) {
	case HeaderMsg:
		p.current = int(msg)
		return p, bbt.Sequence(
			Focus("list"),
			List("clear"),
			List(p.states[p.current].Value()),
		)
	case bbt.KeyMsg:
		switch msg.String() {
		case "enter":
			return p, Header(p.current)
		case "h", "left":
			p.current = mod(p.current-1, len(p.states))
		case "l", "right":
			p.current = mod(p.current+1, len(p.states))
		case "j", "down":
			return p, Focus("list")
		case "tab":
			return p, Focus("list")
		}
	}

	return p, nil
}

func (p *PaneHeader) View() string {
	var views []string
	for i := range p.states {
		state := p.states[i]
		if i == p.current {
			state = state.Copy().
				Foreground(lipgloss.Color("#ff6600")).
				Underline(true)
		}

		views = append(views, state.String())
	}

	return p.style.Render(lipgloss.JoinHorizontal(lipgloss.Top, views...))
}

func (p *PaneHeader) Size() (width, height int) {
	_, v := p.style.GetFrameSize()
	return 0, v + 1
}

func (p *PaneHeader) SetSize(width, height int) {
}

func mod(a, b int) int {
	return (a%b + b) % b
}
