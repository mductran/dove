package tui

import (
	"bytes"

	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	Name         string
	Size         int64
	Dir          bool
	ModifiedTime string
}

type ItemModel struct {
	item Item
}

type ItemMsg struct {
	Item Item
}

func (model ItemModel) Init() tea.Cmd {
	return nil
}

func (model ItemModel) Update(msg tea.Msg) (ItemModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ItemMsg:
		model.item = msg.Item
	}
	return model, nil
}

func (model ItemModel) View() string {
	var buffer bytes.Buffer

	buffer.WriteString(model.item.Name)
	buffer.WriteString("\nFile")

	return buffer.String()
}
