package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO(cmz): make this composable vs all this hard coded stuff
type applicationStateModel struct {
	choices []string
	cursor int
	selected map[int]struct{}
	spinner spinner.Model
	status int
	err error
}

func (m applicationStateModel) Init() tea.Cmd {
	return tea.Batch(
		doHttpRequest,
		m.spinner.Tick,
	)
}

func (m applicationStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO(cmz): Refactor to something else, all these nested switches are a code smell
	switch msg := msg.(type) {
	
	case statusMsg:
		m.status = int(msg)
		return m, nil
	
	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m applicationStateModel) View() string {
	// TODO (cmz): Refactor to Builder pattern
	// TODO (cmz): Create the concept of React-like components to allow enable multiple developers to work on this at one time
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := "Friday, November 10, 2023 - Tasks to Complete\n=================================================\n\n"
	s += fmt.Sprintf("%s %d %s!\n\n", m.spinner.View(), m.status, http.StatusText(m.status))

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func initApplicationStateModel() applicationStateModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	return applicationStateModel{
		spinner: s,	
		choices: []string{"(10:00 - 11:00) Work on Productivity CLI", "(11:00 - 12:30) Exercise", "(13:00 - 14:30) Add Automated Testing to P.A."},
		selected: make(map[int]struct{}),
		status: 102,
	}
}

// TODO(cmz): Move to "Repository" struct, as this fetches Entities (another struct to create)
func doHttpRequest() tea.Msg {
	time.Sleep(5 * time.Second)
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	}
	return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }