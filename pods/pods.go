package pods

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type podlist *v1.PodList

type keyMap struct {
	Help key.Binding
	Quit key.Binding
}

type Model struct {
	help help.Model
	K8sclient *kubernetes.Clientset
	keys keyMap
	list []string
}

func (model Model) Init() tea.Cmd {
	model.list = make([]string, 1)
	return getPods(model)
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case podlist:
		for _, pod := range msg.Items {
			model.list = append(model.list, pod.Name)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
	}

	return model, nil
}

func (model Model) View() string {
	return strings.Join(model.list, ",")
}

func getPods(model Model) tea.Cmd {
	return func() tea.Msg {
		pods, err := model.K8sclient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		return podlist(pods)
	}
}
