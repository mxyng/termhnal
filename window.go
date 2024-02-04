package main

import (
	"fmt"
	"strconv"
	"strings"

	bbt "github.com/charmbracelet/bubbletea"
)

type Window interface {
	Update(bbt.Msg) (Window, bbt.Cmd)
	View() string
}

type WindowView struct {
	header *PaneHeader
	view   *PaneView
	footer *PaneFooter
	active Pane
}

func NewWindowView() *WindowView {
	var window WindowView
	window.view = NewPaneView()
	window.header = NewPaneHeader(
		PaneHeaderItem{
			Name: "Back",
			Func: func() bbt.Cmd {
				return Activate("list")
			},
		},
	)

	window.footer = NewPaneFooter(
		func() string {
			return fmt.Sprintf("%3.f%%", window.view.model.ScrollPercent()*100)
		},
		func() string {
			return ""
		},
	)

	window.active = window.view
	return &window
}

func (w *WindowView) Update(msg bbt.Msg) (Window, bbt.Cmd) {
	switch msg := msg.(type) {
	case ActivateMsg:
		if msg == "toggle" {
			switch w.active.(type) {
			case *PaneHeader:
				msg = "view"
			case *PaneView:
				msg = "header"
			}
		}

		switch strings.ToLower(string(msg)) {
		case "header":
			w.active.Deactivate()
			w.active = w.header.Activate()
		case "view":
			w.active.Deactivate()
			w.active = w.view.Activate()
		}
	case bbt.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return w, Activate("list")
		}
	case bbt.WindowSizeMsg:
		for _, pane := range []Pane{w.header, w.footer, w.view} {
			pane.SetSize(msg.Width, msg.Height)
			width, height := pane.Size()
			msg.Width -= width
			msg.Height -= height
		}
	}

	var cmd bbt.Cmd
	w.active, cmd = w.active.Update(msg)
	return w, cmd
}

func (w *WindowView) View() string {
	var sb strings.Builder
	sb.WriteString(w.header.View())
	sb.WriteString(w.view.View())
	sb.WriteString(w.footer.View())
	return sb.String()
}

type WindowList struct {
	header *PaneHeader
	list   *PaneList
	footer *PaneFooter
	active Pane
}

func NewWindowList() *WindowList {
	var items []PaneHeaderItem
	values := []string{"Top", "New", "Best", "Ask", "Show", "Job"}
	for i := range values {
		value := values[i]
		items = append(items, PaneHeaderItem{
			Name: value,
			Func: func() bbt.Cmd {
				return bbt.Sequence(
					Activate("list"),
					List("clear"),
					List(value),
				)
			},
		})
	}

	var window WindowList
	window.header = NewPaneHeader(items...)
	window.list = NewPaneList()
	window.footer = NewPaneFooter(
		func() string {
			return fmt.Sprintf("%d of %d", window.list.model.Paginator.Page+1, window.list.model.Paginator.TotalPages)
		}, func() string {
			return ""
		},
	)

	window.active = window.list
	return &window
}

func (w *WindowList) Update(msg bbt.Msg) (Window, bbt.Cmd) {
	switch msg := msg.(type) {
	case ActivateMsg:
		if msg == "toggle" {
			switch w.active.(type) {
			case *PaneHeader:
				msg = "list"
			case *PaneList:
				msg = "header"
			}
		}

		switch strings.ToLower(string(msg)) {
		case "header":
			w.active.Deactivate()
			w.active = w.header.Activate()
		case "list":
			w.active.Deactivate()
			w.active = w.list.Activate()
		}
	case bbt.KeyMsg:
		switch msg.String() {
		case "1", "2", "3", "4", "5", "6":
			n, _ := strconv.Atoi(msg.String())
			return w, bbt.Sequence(
				Activate("header"),
				Header(n-1),
			)
		}
	case bbt.WindowSizeMsg:
		for _, pane := range []Pane{w.header, w.footer, w.list} {
			pane.SetSize(msg.Width, msg.Height)
			width, height := pane.Size()
			msg.Width -= width
			msg.Height -= height
		}
	}

	var cmd bbt.Cmd
	w.active, cmd = w.active.Update(msg)
	return w, cmd
}

func (w *WindowList) View() string {
	var sb strings.Builder
	sb.WriteString(w.header.View())
	sb.WriteString(w.list.View())
	sb.WriteString(w.footer.View())
	return sb.String()
}
