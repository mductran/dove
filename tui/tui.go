package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	items ItemList
}

func (app App) Init() tea.Cmd {
	return nil
}

func (app App) View() string {
	return app.items.View()
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return app, tea.Quit
		}
	}

	// hacky way to display current directory with update
	app.items, _ = app.items.Update()

	return app, nil
}

func RunTUI() {

	p := tea.NewProgram(&App{items: ListDir()}, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		panic(err)
	}
}
