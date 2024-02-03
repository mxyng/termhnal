package main

import (
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
	active Pane
}

func NewWindowView() *WindowView {
	window := WindowView{
		header: NewPaneHeader(
			PaneHeaderItem{
				Name: "Back",
				Func: func() bbt.Cmd {
					return Activate("list")
				},
			},
		),
		view: NewPaneView(),
	}

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
	case bbt.WindowSizeMsg:
		for _, pane := range []Pane{w.header, w.view} {
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
	return sb.String()
}

type WindowList struct {
	header *PaneHeader
	list   *PaneList
	active Pane
}

func NewWindowList() *WindowList {
	var items []PaneHeaderItem
	for _, item := range []string{"Top", "New", "Best", "Ask", "Show", "Job"} {
		items = append(items, PaneHeaderItem{
			Name: item,
			Func: func() bbt.Cmd {
				return bbt.Sequence(
					Activate("list"),
					List("clear"),
					List(item),
				)
			},
		})
	}

	window := WindowList{
		header: NewPaneHeader(items...),
		list:   NewPaneList(),
	}

	window.active = window.list
	return &window
}

func (w *WindowList) Update(msg bbt.Msg) (Window, bbt.Cmd) {
	switch msg := msg.(type) {
	case ActivateMsg:
		if msg == "toggle" {
			switch w.active.(type) {
			case *PaneHeader:
				msg = "view"
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
		for _, pane := range []Pane{w.header, w.list} {
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
	return sb.String()
}
