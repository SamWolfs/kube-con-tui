package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SamWolfs/kube-con-tui/pods"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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

type model struct {
	help      help.Model
	K8sclient *kubernetes.Clientset
	keys      keyMap
}

func initModel() model {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return model{
		help:      help.New(),
		K8sclient: clientset,
		keys:      keys,
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
			newModel := pods.Model{
				K8sclient: model.K8sclient,
			}
			return newModel, newModel.Init()
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
