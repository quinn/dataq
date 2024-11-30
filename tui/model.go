package tui

import (
	"context"
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.quinn.io/dataq/queue"
	"go.quinn.io/dataq/worker"
)

type state int

const (
	stateMenu state = iota
	stateStatus
	stateWorker
	stateStep
)

type Model struct {
	state       state
	menuCursor  int
	queue       queue.Queue
	worker      *worker.Worker
	status      []*queue.Task
	lastResult  *worker.TaskResult
	lastUpdated time.Time
	err         error
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7B2CBF")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
)

func NewModel(q queue.Queue, w *worker.Worker) Model {
	return Model{
		queue:       q,
		worker:      w,
		lastUpdated: time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.state != stateMenu {
				m.state = stateMenu
				return m, nil
			}
		}

		if m.state == stateMenu {
			switch msg.String() {
			case "up", "k":
				if m.menuCursor > 0 {
					m.menuCursor--
				}
			case "down", "j":
				if m.menuCursor < 2 {
					m.menuCursor++
				}
			case "enter":
				switch m.menuCursor {
				case 0:
					m.state = stateStatus
					return m, m.updateStatus
				case 1:
					m.state = stateWorker
					// Start the worker in a goroutine
					ctx, _ := context.WithCancel(context.Background())
					go func() {
						if err := m.worker.Start(ctx); err != nil && err != context.Canceled {
							log.Printf("Worker error: %v", err)
						}
					}()
					return m, nil
				case 2:
					m.state = stateStep
					return m, m.processNextTask
				}
			}
		} else if m.state == stateStep {
			switch msg.String() {
			case "enter", "space":
				return m, m.processNextTask
			}
		}
	case statusMsg:
		m.status = msg.tasks
		m.err = msg.err
		m.lastUpdated = time.Now()
		return m, m.updateStatus
	case taskResultMsg:
		m.lastResult = msg.result
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateMenu:
		return m.viewMenu()
	case stateStatus:
		return m.viewStatus()
	case stateWorker:
		return m.viewWorker()
	case stateStep:
		return m.viewStep()
	default:
		return "Unknown state"
	}
}

func (m Model) viewMenu() string {
	s := titleStyle.Render("DataQ") + "\n\n"
	s += "What would you like to do?\n\n"

	choices := []string{"View Status", "Start Worker", "Step Through Tasks"}

	for i, choice := range choices {
		cursor := " "
		if m.menuCursor == i {
			cursor = ">"
			choice = selectedStyle.Render(choice)
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n(press q to quit)\n"
	return s
}

func (m Model) viewStep() string {
	s := titleStyle.Render("Step Through Tasks") + "\n\n"

	if m.lastResult == nil {
		s += "Press ENTER to process the next task\n"
		s += "Press ESC to return to menu\n"
		return s
	}

	s += fmt.Sprintf("Task ID: %s\n", m.lastResult.Task.ID)
	s += fmt.Sprintf("Plugin: %s\n", m.lastResult.Task.PluginID)
	s += fmt.Sprintf("Status: %s\n", m.lastResult.Task.Status)

	if m.lastResult.Error != nil {
		s += errorStyle.Render(fmt.Sprintf("Error: %v\n", m.lastResult.Error))
	} else {
		s += successStyle.Render("Task completed successfully\n")
	}

	s += "\nPress ENTER to process next task\n"
	s += "Press ESC to return to menu\n"
	return s
}

func (m Model) viewStatus() string {
	s := titleStyle.Render("Queue Status") + "\n\n"
	s += fmt.Sprintf("Last updated: %s\n\n", m.lastUpdated.Format(time.RFC3339))

	if m.err != nil {
		return s + errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if len(m.status) == 0 {
		return s + "No tasks in queue"
	}

	for _, task := range m.status {
		s += fmt.Sprintf("Task ID: %s\n", task.ID)
		s += fmt.Sprintf("Plugin: %s\n", task.PluginID)
		s += fmt.Sprintf("Status: %s\n", task.Status)
		if task.Error != "" {
			s += fmt.Sprintf("Error: %s\n", task.Error)
		}
		s += "\n"
	}

	s += "\n(press esc to go back, q to quit)\n"
	return s
}

func (m Model) viewWorker() string {
	s := titleStyle.Render("Worker Running") + "\n\n"
	s += "Processing tasks...\n"
	s += "\n(press q to quit)\n"
	return s
}

type statusMsg struct {
	tasks []*queue.Task
	err   error
}

type taskResultMsg struct {
	result *worker.TaskResult
	err    error
}

func (m Model) updateStatus() tea.Msg {
	tasks, err := m.queue.List(context.Background(), "")
	return statusMsg{tasks: tasks, err: err}
}

func (m Model) processNextTask() tea.Msg {
	result, err := m.worker.ProcessSingleTask(context.Background())
	return taskResultMsg{result: result, err: err}
}
