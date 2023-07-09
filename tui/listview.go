package tui

import (
	"bytes"
	"io/ioutil"

	tea "github.com/charmbracelet/bubbletea"
)

type ItemType int8

type Item struct {
	Name string
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
			Name: file.Name(),
			Dir:  file.Mode().IsDir(),
			Size: file.Size(),
		}
		items = append(items, item)
	}
	return items
}

func ListDir() ItemList {
	return ItemList{
		items: ls("."),
	}
}

func (list ItemList) Init() tea.Cmd {
	return nil
}

func (list ItemList) Update() (ItemList, tea.Msg) {
	return list, nil
}

func (list ItemList) View() string {
	var buffer bytes.Buffer

	for _, item := range list.items {
		buffer.WriteString(item.Name + "\n")
	}

	return buffer.String()
}
