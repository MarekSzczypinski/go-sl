package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
)

type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(40*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type model struct {
	screenWidth  int
	screenHeight int
	artXPos      int
	artWidth     int
	asciiArt     string
}

func initModel() model {
	return model{
		screenWidth:  0,
		screenHeight: 0,
		artXPos:      0,
		artWidth:     52,
		asciiArt: ` _   _      _ _        __        __         _     _
| | | | ___| | | ___   \ \      / /__  _ __| | __| |
| |_| |/ _ \ | |/ _ \   \ \ /\ / / _ \| '__| |/ _\ |
|  _  |  __/ | | (_) |   \ V  V / (_) | |  | | (_| |
|_| |_|\___|_|_|\___/     \_/\_/ \___/|_|  |_|\__,_|`,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
		m.artXPos = msg.Width
	case TickMsg:
		m.artXPos--
		if m.artXPos < -m.artWidth {
			return m, tea.Quit
		}
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	lines := strings.Split(m.asciiArt, "\n")
	var sb strings.Builder

	for _, line := range lines {
		spaces := ""
		if m.artXPos >= 0 {
			spaces = strings.Repeat(" ", m.artXPos)
			sb.WriteString(spaces + line + "\n")
		} else if m.artXPos < 0 && m.artXPos > -m.artWidth {
			sb.WriteString(line[m.artXPos*(-1):] + "\n")
		}
	}

	// lipgloss style for centering
	style := lipgloss.NewStyle().Height(m.screenHeight)
	centeredContent := style.AlignVertical(lipgloss.Center).Render(sb.String())
	return centeredContent
}

func main() {
	cmd := &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "sl",
		Usage:                  "cure your bad habit of mistyping",
		Description: `sl is a highly advanced animation program for curing your bad habit of mistyping.
This app is a Go port of the original sl program written in C by Toyoda Masashi.`,
		Version: "1.0.0",
		Authors: []any{
			mail.Address{Name: "Marek Szczypiński", Address: "markacy@gmail.com"},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "a",
				Usage: "An accident is occurring. People cry for help.",
			},
			&cli.BoolFlag{
				Name:  "l",
				Usage: "Little version",
			},
			&cli.BoolFlag{
				Name:  "F",
				Usage: "It flies like the galaxy express 999.",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("a") {
				fmt.Println("accident")
			}
			if cmd.Bool("l") {
				fmt.Println("little")
			}
			if cmd.Bool("F") {
				fmt.Println("flies")
			}
			p := tea.NewProgram(initModel(), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
