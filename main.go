package main;

import (
	 tea "github.com/charmbracelet/bubbletea"
	"log"
	// "os"
	// "github.com/leogem2003/directchan"
)

// func main() {
// 	settings := new(connection.ConnectionSettings)
// 	settings.Key = os.Args[2]
// 	settings.STUN = []string{"stun:stun.l.google.com:19302"}
// 	settings.Signaling = "ws://localhost:8080"
// 	settings.BufferSize = 1
// 	switch os.Args[1] {
// 	case "answer":
// 		settings.Operation = 1
// 	case "offer":
// 		settings.Operation = 0
// 	default:
// 		log.Println("Error: invalid operation", os.Args[1])
// 		return
// 	}
//
// 	conn, err := connection.FromSettings(settings)
// 	if err != nil {
// 		if conn != nil {
// 			conn.CloseAll()
// 		}
// 		log.Fatalln(err)
// 	}
//
// 	conn.In <- []byte("sergio")
// 	p := tea.NewProgram(initialModel(conn.In, conn.Out))
// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	p := tea.NewProgram(initialLogin())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
