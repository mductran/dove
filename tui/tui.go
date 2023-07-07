package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type App struct{}

func (app App) Init() tea.Cmd {
	return nil
}

func (app App) View() string {
	return ""
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return app, tea.Quit
		}
	}
	return app, nil
}

func RunTUI() {
	app := App{}
	p := tea.NewProgram(&app, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		panic(err)
	}
}
