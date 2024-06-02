package progress

import (
	"bytes"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/constraints"
)

type Model struct {
	tasks []*TaskModel

	width int
}

func (m *Model) Init() tea.Cmd {
	return tickCmd()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		for i := range m.tasks {
			m.tasks[i].width = min(m.width, 80)
		}
		return m, nil

	case tickMsg:
		quit := true
		for i := range m.tasks {
			m.tasks[i].width = min(m.width, 80)
			if m.tasks[i].percent < 1.0 {
				quit = false
			}
		}
		if quit {
			return m, tea.Quit
		}

		return m, tickCmd()

	default:
		return m, nil
	}

}

func (m *Model) View() string {
	var buf bytes.Buffer
	for _, t := range m.tasks {
		buf.WriteString(t.View())
		buf.WriteByte('\n')
	}
	return buf.String()
}

type TaskModel struct {
	description   string
	current       int64
	total         int64
	unitFormatter func(i int64) string
	percent       float64

	progress  progress.Model
	startTime time.Time
	width     int

	isFinish   bool
	finishView string
}

type TaskModelOption func(*TaskModel)

func WithDescription(description string) TaskModelOption {
	return func(t *TaskModel) {
		t.description = description
	}
}

func WithUnitFormatter(f func(i int64) string) TaskModelOption {
	return func(t *TaskModel) {
		t.unitFormatter = f
	}
}

func WithTotal[T constraints.Integer](total T) TaskModelOption {
	return func(t *TaskModel) {
		t.total = int64(total)
	}
}

func WithProgress(progress progress.Model) TaskModelOption {
	return func(t *TaskModel) {
		t.progress = progress
	}
}

func NewTaskModel(description string, opts ...TaskModelOption) *TaskModel {
	t := &TaskModel{
		description: description,
		startTime:   time.Now(),
		unitFormatter: func(i int64) string {
			return fmt.Sprint(i)
		},
		progress: progress.New(
			func(p *progress.Model) {
				p.Empty = '━'
				p.Full = p.Empty
				// grey
				p.EmptyColor = "#737373"
			},
			// red
			progress.WithSolidFill("#F92672"),
		),
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func (p *TaskModel) View() string {
	if p.isFinish {
		return p.finishView
	}

	p.progress.Width = p.width
	if p.current > 0 && p.total > 0 {
		p.percent = float64(p.current) / float64(p.total)
	}
	if p.percent >= 1 {
		p.percent = 1
		p.isFinish = true
	}

	var buf bytes.Buffer
	buf.WriteString(TextStyle(p.description))

	// progress bar
	if p.percent >= 1 {
		// green
		p.progress.FullColor = "#729C1F"
		buf.WriteString(" " + p.progress.ViewAs(1.0))
	} else if p.percent > 0 {
		buf.WriteString(" " + p.progress.ViewAs(p.percent))
	}

	// calculate per second
	if p.current > 0 {
		perSecond := float64(p.current) / time.Since(p.startTime).Seconds()
		buf.WriteString(textStylef(" %s/s", p.unitFormatter(int64(perSecond))))
	}

	// current/total
	if p.total > 0 {
		buf.WriteString(textStylef(" %s/%s", p.unitFormatter(p.current), p.unitFormatter(p.total)))
	} else if p.current > 0 {
		buf.WriteString(textStylef(" %s/?", p.unitFormatter(p.current)))
	}

	// calculate the time remaining
	if p.isFinish {
		buf.WriteString(
			textStylef(" Cost: %v", time.Since(p.startTime).Truncate(10*time.Millisecond)),
		)
	} else if p.current > 0 && p.total > 0 {
		perCostTime := time.Since(p.startTime) / time.Duration(p.current)
		needTime := time.Duration(p.total-p.current) * perCostTime
		buf.WriteString(textStylef(" ETA: %v", needTime.Truncate(10*time.Millisecond)))
	} else if p.percent > 0 {
		per2int := int(p.percent * 100)
		perCostTime := time.Since(p.startTime) / time.Duration(per2int)
		needTime := time.Duration(100-per2int) * perCostTime
		buf.WriteString(textStylef(" ETA: %v", needTime.Truncate(10*time.Millisecond)))
	}

	if p.isFinish {
		p.finishView = buf.String()
		return p.finishView
	}

	return buf.String()
}
