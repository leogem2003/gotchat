package main

import (
	 tea "github.com/charmbracelet/bubbletea"
)

type errorModel struct {
	message string 
}

func initialError(msg string) errorModel {
	return errorModel {
		message: msg,
	}
}

func (e errorModel) Init() tea.Cmd {
	return nil
}


func (e errorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return e, tea.Quit
		}
	}
	return e, nil
}


func (e errorModel) View() string {
	return e.message
}
