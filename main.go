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

const (
	bigEngineBase = `"      ====        ________                ___________     "
"  _D _|  |_______/        \\__I_I_____===__|_________|    "
"   |(_)---  |   H\\________/ |   |        =|___ ___|      "
"   /     |  |   H  |  |     |   |         ||_| |_||       "
"  |      |  |   H  |__--------------------| [___] |       "
"  | ________|___H__/__|_____/[][]~\\______|       |       "
"  |/ |   |-----------I_____I [][] []  D   |=======|__     "`

	bigEngineWheels1 = `"__/ =| o |=-~~\\  /~~\\  /~~\\  /~~\\ ____Y___________|__ "
" |/-=|___|=    ||    ||    ||    |_____/~\\___/           "
"  \\_/      \\O=====O=====O=====O_/   \\__/               "`

	bigEngineWheels2 = `"__/ =| o |=-~~\\  /~~\\  /~~\\  /~~\\ ____Y___________|__ "
" |/-=|___|=O=====O=====O=====O   |_____/~\\___/           "
"  \\_/      \\__/  \\__/  \\__/       \\__/               "`

	bigEngineWheels3 = `"__/ =| o |=-O=====O=====O=====O /~~\\ ____Y___________|__ "
" |/-=|___|=    ||    ||    ||    |_____/~\\___/           "
"  \\_/      \\__/  \\__/  \\__/       \\__/               "`

	bigEngineWheels4 = `"__/ =| o |=-~O=====O=====O=====O/~~\\ ____Y___________|__ "
" |/-=|___|=    ||    ||    ||    |_____/~\\___/           "
"  \\_/      \\__/  \\__/  \\__/       \\__/               "`

	bigEngineWheels5 = `"__/ =| o |=-~~\\  /~~\\  /~~\\  /~~\\ ____Y___________|__ "
" |/-=|___|=   O=====O=====O=====O|_____/~\\___/           "
"  \\_/      \\__/  \\__/  \\__/       \\__/               "`

	bigEngineWheels6 = `"__/ =| o |=-~~\\  /~~\\  /~~\\  /~~\\ ____Y___________|__ "
" |/-=|___|=    ||    ||    ||    |_____/~\\___/           "
"  \\_/      \\_O=====O=====O=====O/   \\__/               "`
)

var (
	bigEngineWheels = []string{bigEngineWheels1, bigEngineWheels6, bigEngineWheels5, bigEngineWheels4, bigEngineWheels3, bigEngineWheels2}
)

type model struct {
	screenWidth  int
	screenHeight int
	artXPos      int
	artWidth     int
	engineWheels int
}

func initModel() model {
	return model{
		screenWidth:  0,
		screenHeight: 0,
		artXPos:      0,
		artWidth:     61,
		engineWheels: 0,
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
		m.engineWheels++
		if m.engineWheels > 5 {
			m.engineWheels = 0
		}
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	baseLines := strings.Split(bigEngineBase, "\n")
	wheelsLines := strings.Split(bigEngineWheels[m.engineWheels], "\n")
	var sb strings.Builder

	for _, line := range baseLines {
		spaces := ""
		if m.artXPos >= 0 {
			spaces = strings.Repeat(" ", m.artXPos)
			sb.WriteString(spaces + line + "\n")
		} else if m.artXPos < 0 && m.artXPos > -m.artWidth {
			sb.WriteString(line[m.artXPos*(-1):] + "\n")
		}
	}

	for _, line := range wheelsLines {
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
