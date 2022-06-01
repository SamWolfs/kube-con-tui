package pods

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {}

func initModel(model tea.Model) Model {
	return Model{}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
	}

	return model, nil
}

func (model Model) View() string {
	return "We're PODS now"
}
