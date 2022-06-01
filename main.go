package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/SamWolfs/kube-con-tui/pods"
)

type keyMap struct {
	Help key.Binding
	Pods key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Pods, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Pods},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Pods: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "pods"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct{
	keys keyMap
	help help.Model
}

func initModel() model {
	return model{
		keys: keys,
		help: help.New(),
	}
}

func (model model) Init() tea.Cmd {
	return nil
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, model.keys.Help):
			model.help.ShowAll = !model.help.ShowAll
		case key.Matches(msg, model.keys.Pods):
			return pods.Model{}, nil
		case key.Matches(msg, model.keys.Quit):
			return model, tea.Quit
		}
	}

	return model, nil
}

func (model model) View() string {
	s := "Hello World"
	helpView := model.help.View(model.keys)
	return s + "\n" + helpView
}

func main() {
	program := tea.NewProgram(initModel())
	if err := program.Start(); err != nil {
		fmt.Printf("Error during startup: %v", err)
		os.Exit(1)
	}
}
