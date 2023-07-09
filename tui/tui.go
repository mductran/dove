package tui

import (
	"io/ioutil"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct{}

type ItemType int8

type Item struct {
	Path string
	Size int64
	Dir  bool
}

type ItemList struct {
	items []Item
}

func ls(path string) []Item {
	items := []Item{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		item := Item{
			Path: file.Name(),
			Dir:  file.Mode().IsDir(),
			Size: file.Size(),
		}
		items = append(items, item)
	}
	return items
}

func list() ItemList {
	return ItemList{
		items: ls("."),
	}
}

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
