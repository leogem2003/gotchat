package main;


import (
	"fmt"
	 tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/leogem2003/directchan"
)


type loginModel struct {
	inputs 				[]textinput.Model
	focusedIndex  int
	err 					error
}

func operationValidator(input string) error {
	if input != "o" && input != "a" {
		return fmt.Errorf("Invalid input: expected o|a, got %s", input)
	}
	return fmt.Errorf("error");
}

const (
	ipIndex int = iota 
	keyIndex
	roleIndex
	nameIndex
)

func initialLogin() loginModel {
	inputs := make([]textinput.Model, 4)
	ipInput := textinput.New()
	ipInput.Focus()
	keyInput := textinput.New()
	role := textinput.New()
	role.Validate = operationValidator
	role.Placeholder = "[o]ffer | [a]nswer"

	name := textinput.New()

	inputs[ipIndex] = ipInput
	inputs[keyIndex] = keyInput
	inputs[roleIndex] = role
	inputs[nameIndex] = name

	return loginModel {
		inputs, 0, nil,
	}
}

func (l loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func createConnection(host string, key string, role string, name string) (tea.Model, tea.Cmd) {
	settings := connection.ConnectionSettings {
		Key: key,
		STUN: []string{"stun:stun.l.google.com:19302"},
		Signaling: host,
		BufferSize: 1,
	}

	if role == "o" {
		settings.Operation = 0
	} else {
		settings.Operation = 1
	}

	conn, err := connection.FromSettings(&settings)
	if err != nil {
		if conn != nil {
			conn.CloseAll()
		}
		return initialError(err.Error()), nil
	}

	go func () {conn.In <- []byte(name)} ()
	return initialModel(conn.In, conn.Out), nil
}

func (l loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(l.inputs))	

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			l.inputs[l.focusedIndex].Blur()
			l.focusedIndex++
			if l.focusedIndex < len(l.inputs) {
				l.inputs[l.focusedIndex].Focus()
			} else {
				return createConnection(
					l.inputs[ipIndex].Value(),
					l.inputs[keyIndex].Value(),
					l.inputs[roleIndex].Value(),
					l.inputs[nameIndex].Value(),
				)
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return l, tea.Quit
		}
	case error:
		l.err = msg
		return l, nil
	}

	for i := range l.inputs {
		l.inputs[i], cmds[i] = l.inputs[i].Update(msg)
	}

	return l, tea.Batch(cmds...)
}

func (l loginModel) View() string {
	errorMsg := ""
	if l.err != nil {
		errorMsg = l.err.Error()
	}

	return errorMsg + fmt.Sprintf(
`IP:port or server address
%s	
Key
%s
Operation
%s
Name
%s`,
		l.inputs[ipIndex].View(),
		l.inputs[keyIndex].View(),
		l.inputs[roleIndex].View(),
		l.inputs[nameIndex].View(),
	) + "\n"
}
