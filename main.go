package main

import (
	bbt "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list   *WindowList
	view   *WindowView
	active Window
}

func NewModel() *Model {
	model := Model{
		list: NewWindowList(),
		view: NewWindowView(),
	}

	model.active = model.list
	return &model
}

func (m *Model) Init() bbt.Cmd {
	return bbt.Sequence(
		Activate("list"),
		List("top"),
	)
}

func (m *Model) Update(msg bbt.Msg) (bbt.Model, bbt.Cmd) {
	switch msg := msg.(type) {
	case bbt.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// mask off ctrl+c
			return m, nil
		case "ctrl+d":
			return m, bbt.Quit
		}
	case ActivateMsg:
		switch msg {
		case "list":
			m.active = m.list
		case "view":
			m.active = m.view
		}
	case bbt.WindowSizeMsg:
		var cmds []bbt.Cmd
		for _, window := range []Window{m.list, m.view} {
			_, cmd := window.Update(msg)
			cmds = append(cmds, cmd)
		}

		return m, bbt.Batch(cmds...)
	}

	var cmd bbt.Cmd
	m.active, cmd = m.active.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.active.View()
}

func main() {
	if _, err := bbt.NewProgram(NewModel(), bbt.WithAltScreen()).Run(); err != nil {
		panic(err)
	}
}
