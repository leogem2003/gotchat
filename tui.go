package main;

import (
	"fmt"
	"strings"
	 tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)


const gap = "\n\n"


type (
	errMsg error
)


type model struct {
	viewport    viewport.Model
	input     	textarea.Model
	messages    []string
	senderStyle lipgloss.Style
	err         error
	in  				chan []byte
	out 				chan []byte
	name 				string
	scrolling 	bool
}

type incomingMessage string

func listen(c chan []byte) tea.Cmd {
	return func () tea.Msg {
		return incomingMessage(<- c)
	}
}

func initialModel(in chan[] byte, out chan[] byte) model {
	ta := initInputArea()
	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	m := model{
		input:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		in: 				 in,
		out: 				 out,
		name: 			 string(<-out),
		scrolling: 	 false,
	}

	return m
}

func (m *model) renderMessages() {
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, listen(m.out))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		vpCmd tea.Cmd
	)
	
	cmds := make([]tea.Cmd, 1)

	m.input, _ = m.input.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.input.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.renderMessages()
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.input.Value())
			return m, tea.Quit
		case tea.KeyCtrlV:
			m.scrolling = !m.scrolling 
			if (m.scrolling) {
				m.input.Blur()
			} else {
				m.input.Focus()
			}
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.input.Value())
			m.in <- []byte(m.input.Value())
			m.renderMessages()
			m.viewport.GotoBottom()
			m.input.Reset()
		case tea.KeyRunes:
			if m.scrolling {
				m.viewport, vpCmd = m.viewport.Update(msg)
				cmds = append(cmds, vpCmd)
			}
		}
	case incomingMessage:
		m.messages = append(m.messages,
			fmt.Sprintf("%s: %s", m.name, string(msg)),
		)
		m.renderMessages()
		m.viewport.GotoBottom()
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}
	cmds = append(cmds, listen(m.out))
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.input.View(),
	)
}

