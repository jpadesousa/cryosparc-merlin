package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	skipStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

// Progress bar padding
const (
	padding  = 2
	maxWidth = 80
)

// =============================================================================
// Model
// =============================================================================

type stepStatus struct {
	message string
	status  string
}

type model struct {
	spinner  spinner.Model
	progress progress.Model

	steps     []Step       // ordered list of steps to execute
	history   []stepStatus // it keeps the history of all step types
	current   int          // index of the step currently being executed
	completed []string     // messages for completed steps
	skipped   []string     // messages for skippeds steps
	err       error        // first error encountered, if any

	form *huh.Form // form to prompt the user
}

// =============================================================================
// Step
// =============================================================================
// Step represents a single unit of work in the progress UI
type Step struct {
	Message          string
	CompletedMessage string
	SkipMessage      string

	Condition func() (bool, error)
	Skip      func() bool

	ShowProgressBar bool

	Action func(update func(float64)) error
	Exec   func() *exec.Cmd

	Prompt func() *huh.Form
}

// stepResult is sent back to the update loop after a step finishes
type stepResult struct {
	err error
}

type stepSkipped struct {
	message string
	err     error
}

// =============================================================================
// Commands and Messages
// =============================================================================

// progressMsg is a type used when there is a progress output
type progressMsg float64

type startPrompt struct {
	form *huh.Form
}

var program *tea.Program

// runStep creates a Bubble Tea command that executes a step asynchronously
func runStep(step Step) tea.Cmd {
	return func() tea.Msg {
		shouldRun := false

		if step.Condition != nil {
			var err error
			shouldRun, err = step.Condition()
			if err != nil {
				return stepResult{err: err}
			}
		}

		if step.Prompt != nil && shouldRun {
			return startPrompt{form: step.Prompt()}
		}

		return executeStep(step)()
	}
}

func executeStep(step Step) tea.Cmd {
	return func() tea.Msg {
		if step.Skip != nil && step.Skip() {
			return stepSkipped{
				message: step.SkipMessage,
			}
		}

		if step.Exec != nil {
			return tea.ExecProcess(
				step.Exec(),
				func(err error) tea.Msg {
					return stepResult{err: err}
				},
			)()
		}

		err := step.Action(func(percent float64) {
			program.Send(progressMsg(percent))
		})

		return stepResult{err: err}
	}
}

// =============================================================================
// Initialization
// =============================================================================

// Init starts the spinner animation and immediately begins executing
// the first step.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.progress.Init(),
		runStep(m.steps[0]),
	)
}

// =============================================================================
// Update
// =============================================================================

// Update processes incoming messages and advances the program state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Ctrl+c exits
	if key, ok := msg.(tea.KeyPressMsg); ok {
		if key.Mod == tea.ModCtrl && key.Code == 'c' {
			return m, tea.Quit
		}
	}

	// Print form if exists
	if m.form != nil {
		form, cmd := m.form.Update(msg)
		m.form = form.(*huh.Form)

		if m.form.State == huh.StateCompleted {
			m.form = nil
			return m, tea.Batch(
				executeStep(m.steps[m.current]),
				m.spinner.Tick,
			)
		}

		return m, cmd
	}

	switch msg := msg.(type) {

	// If startPrompt is requested
	case startPrompt:
		m.form = msg.form
		return m, m.form.Init()

	// Skip step if condition is met
	case stepSkipped:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}

		m.skipped = append(
			m.skipped,
			m.steps[m.current].SkipMessage,
		)

		m.history = append(m.history, stepStatus{
			status:  "skipped",
			message: m.steps[m.current].SkipMessage,
		})

		m.current++

		if m.current == len(m.steps) {
			return m, tea.Quit
		}

		return m, tea.Batch(
			runStep(m.steps[m.current]),
			m.spinner.Tick)

	// Initial window size and resizing if the terminal windows changes
	case tea.WindowSizeMsg:
		m.progress.SetWidth(msg.Width - padding*2 - 4)
		if m.progress.Width() > maxWidth {
			m.progress.SetWidth(maxWidth)
		}
		return m, nil

	// Returns a float message (progress bar)
	case progressMsg:
		cmd := m.progress.SetPercent(float64(msg))
		return m, cmd

	// Progress animation step
	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	// Spinner animation tick
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// A step has finished executing
	case stepResult:
		if msg.err != nil {
			m.err = msg.err

			return m, tea.Sequence(
				func() tea.Msg { return nil },
				tea.Quit,
			)
		}

		// Record the completed step so it can be rendered with a checkmark.
		if m.steps[m.current].CompletedMessage != "" {
			m.completed = append(
				m.completed,
				m.steps[m.current].CompletedMessage,
			)

			m.history = append(m.history, stepStatus{
				status:  "completed",
				message: m.steps[m.current].CompletedMessage,
			})
		}

		// Advance to the next step.
		m.current++

		// If every step has completed, exit the program.
		if m.current == len(m.steps) {
			return m, tea.Quit
		}

		// Otherwise, start the next step.
		return m, tea.Batch(
			runStep(m.steps[m.current]),
			m.spinner.Tick)
	}

	return m, nil
}

// =============================================================================
// View
// =============================================================================

// View renders the current UI based on the model state.
func (m model) View() tea.View {

	// If there is a form, just show the form
	if m.form != nil {
		return tea.NewView(
			m.form.View(),
		)
	}

	var b strings.Builder

	for _, step := range m.history {
		switch step.status {
		case "completed":
			fmt.Fprintf(
				&b,
				"%s %s\n",
				successStyle.Render("✔"),
				step.message,
			)

		case "skipped":
			fmt.Fprintf(
				&b,
				"%s %s\n",
				skipStyle.Render("↷"),
				step.message,
			)
		}
	}

	// If execution failed, render the error and stop.
	if m.err != nil {
		fmt.Fprintf(
			&b,
			"%s %v\n",
			errorStyle.Render("✘"),
			m.err,
		)
		return tea.NewView(b.String())
	}

	// Show the currently running step with the animated spinner.
	if m.current < len(m.steps) && m.steps[m.current].Message != "" {
		fmt.Fprintf(
			&b,
			"%s %s\n",
			m.spinner.View(),
			infoStyle.Render(m.steps[m.current].Message),
		)

		if m.steps[m.current].ShowProgressBar {
			fmt.Fprintf(
				&b,
				"%s\n",
				m.progress.View(),
			)
		}
	}

	return tea.NewView(b.String())
}

// =============================================================================
// Run Program
// =============================================================================
func Run(steps []Step) error {
	s := spinner.New()
	s.Spinner = spinner.Dot

	p := progress.New(
		progress.WithDefaultBlend(),
	)

	m := model{
		spinner:  s,
		progress: p,
		steps:    steps,
	}

	program = tea.NewProgram(m)

	_, err := program.Run()

	return err
}
