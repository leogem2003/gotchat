package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/leogem2003/directchan"
)

// Possible calls:
// no args: TUI login
// 1 arg: file containing login info
// 4 args: login info
// login info are signaling ip, key, role (offer | answer), name
func main() {
	var p *tea.Program
	var e error
	switch len(os.Args) {
	case 1:
		p, e = bootstrap()
	case 2:
		p, e = fromFile(os.Args[1])
	case 5:
		p, e = fromArgs(os.Args[1:])
	default:
		log.Fatalf("Invalid argument number: expected 0, 1 or 4, found %d", len(os.Args)-1)
	}
	
	if e != nil {
		log.Fatal(e)
		return
	}

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func fromFile(file string) (*tea.Program, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	args := strings.Split(strings.TrimSpace(string(b)), " ")
	return fromArgs(args)
}

func fromArgs(args []string) (*tea.Program, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("Expected 4 arguments, got %d", len(args))
	}

	settings := new(connection.ConnectionSettings)
	settings.Signaling = args[0] 
	settings.Key = args[1] 
	settings.STUN = []string{"stun:stun.l.google.com:19302"}
	settings.BufferSize = 1
	switch args[2] {
	case "answer":
		settings.Operation = 1
	case "offer":
		settings.Operation = 0
	default:
		return nil, fmt.Errorf("Invalid operation: %s", args[2])
	}

	conn, err := connection.FromSettings(settings)
	if err != nil {
		if conn != nil {
			conn.CloseAll()
		}
		log.Fatalln(err)
	}

	conn.In <- []byte(args[3])
	p := tea.NewProgram(initialModel(conn.In, conn.Out))
	return p, nil
}

func bootstrap() (*tea.Program, error) {
	p := tea.NewProgram(initialLogin())
	return p, nil	
}
