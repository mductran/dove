package tui

import (
	"bytes"
	"io/ioutil"

	tea "github.com/charmbracelet/bubbletea"
)

type ListView struct {
	items        []Item
	currentIndex int
}

func ls(path string) []Item {
	items := []Item{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		item := Item{
			Name:         file.Name(),
			Dir:          file.Mode().IsDir(),
			Size:         file.Size(),
			ModifiedTime: file.ModTime().String(),
		}
		items = append(items, item)
	}
	return items
}

func ListDir() ListView {
	return ListView{
		items: ls("."),
	}
}

func (list ListView) Init() tea.Cmd {
	return nil
}

func (list ListView) Update(msg tea.Msg) (ListView, tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			list.currentIndex -= 1

		case "down":
			list.currentIndex += 1
		}

		// wrap around
		if list.currentIndex < 0 {
			list.currentIndex = len(list.items) - 1
		}
		if list.currentIndex >= len(list.items) {
			list.currentIndex = 0
		}
	}

	return list, nil
}

func (list ListView) View() string {
	var buffer bytes.Buffer

	for index, item := range list.items {

		if index == list.currentIndex {
			buffer.WriteString(" > ")
		} else {
			buffer.WriteString("   ")
		}

		buffer.WriteString(item.Name + "\n")
	}

	return buffer.String()
}

func (list ListView) SelectItem() Item {
	return list.items[list.currentIndex]
}
