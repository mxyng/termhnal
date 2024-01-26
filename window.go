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
	view   *PaneView
	active Pane
}

func NewWindowView() *WindowView {
	window := WindowView{
		view: NewPaneView(),
	}

	window.active = window.view
	return &window
}

func (w *WindowView) Update(msg bbt.Msg) (Window, bbt.Cmd) {
	switch msg := msg.(type) {
	case bbt.WindowSizeMsg:
		for _, pane := range []Pane{w.view} {
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
	sb.WriteString(w.view.View())
	return sb.String()
}

type WindowList struct {
	header *PaneHeader
	list   *PaneList
	active Pane
}

func NewWindowList() *WindowList {
	window := WindowList{
		header: NewPaneHeader("Top", "New", "Best", "Ask", "Show", "Job"),
		list:   NewPaneList(),
	}

	window.active = window.list
	return &window
}

func (w *WindowList) Update(msg bbt.Msg) (Window, bbt.Cmd) {
	switch msg := msg.(type) {
	case FocusMsg:
		switch strings.ToLower(string(msg)) {
		case "header":
			w.active = w.header
		case "list":
			w.active = w.list
		}
	case bbt.KeyMsg:
		switch msg.String() {
		case "1", "2", "3", "4", "5", "6":
			n, _ := strconv.Atoi(msg.String())
			return w, bbt.Sequence(
				Focus("header"),
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
