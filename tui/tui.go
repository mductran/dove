package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	itemView ItemModel
	listView ListView
}

func (app App) Init() tea.Cmd {
	return nil
}

func (app App) View() string {
	// add color
	// return app.itemView.View()
	return lipgloss.JoinHorizontal(lipgloss.Top, app.listView.View(), app.itemView.View())
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
	app.listView, _ = app.listView.Update(msg)
	app.itemView, _ = app.itemView.Update(ItemMsg{Item: app.listView.SelectItem()})

	return app, nil
}

func RunTUI() {

	p := tea.NewProgram(&App{listView: ListDir()}, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		panic(err)
	}
}
