package progress

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cockroachdb/errors"
)

var (
	// TickInterval is the interval at which the progress bar updates.
	TickInterval = 300 * time.Millisecond
	// TextStyle is the style of the text.
	TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

func textStylef(format string, a ...any) string {
	return TextStyle(fmt.Sprintf(format, a...))
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(TickInterval, func(time.Time) tea.Msg {
		return tickMsg(time.Now())
	})
}

type Progress struct {
	tasksMap map[string]int
	model    *Model
}

func New(tasks ...*TaskModel) *Progress {
	tasksMap := make(map[string]int)
	for i, task := range tasks {
		if _, ok := tasksMap[task.description]; ok {
			panic(fmt.Sprintf("task %s already exists", task.description))
		}
		tasksMap[task.description] = i
	}
	return &Progress{
		tasksMap: tasksMap,
		model: &Model{
			tasks: tasks,
		},
	}
}

func (p *Progress) Start() {
	go func() {
		if _, err := tea.NewProgram(
			p.model,
			tea.WithInput(nil),
			tea.WithoutSignalHandler(),
		).Run(); err != nil {
			fmt.Printf("error: %+v\n", errors.WithStack(err))
		}
	}()
}

func (p *Progress) Finish() {
	for _, task := range p.model.tasks {
		task.percent = 1.0
		if task.total > 0 {
			task.current = task.total
		}
	}
	time.Sleep(2 * TickInterval)
}

func (p *Progress) AddTask(task *TaskModel) {
	if _, ok := p.tasksMap[task.description]; ok {
		panic(fmt.Sprintf("task %s already exists", task.description))
	}
	p.tasksMap[task.description] = len(p.model.tasks)
	p.model.tasks = append(p.model.tasks, task)
}

func (p *Progress) findTask(description string) (*TaskModel, bool) {
	if i, ok := p.tasksMap[description]; ok {
		return p.model.tasks[i], true
	}
	return nil, false
}

func (p *Progress) Incr(amount int, description ...string) {
	if len(description) == 0 {
		p.model.tasks[0].current += int64(amount)
		return
	}

	for _, desc := range description {
		if task, ok := p.findTask(desc); ok {
			task.current += int64(amount)
		} else {
			panic(fmt.Sprintf("task %s not found", desc))
		}
	}

}

func (p *Progress) IncrPercent(percent float64, description ...string) {
	if len(description) == 0 {
		p.model.tasks[0].percent += percent
		return
	}

	for _, desc := range description {
		if task, ok := p.findTask(desc); ok {
			task.percent += percent
		} else {
			panic(fmt.Sprintf("task %s not found", desc))
		}
	}

}

func (p *Progress) SetCurrent(current int64, description ...string) {
	if len(description) == 0 {
		p.model.tasks[0].current = current
		return
	}

	for _, desc := range description {
		if task, ok := p.findTask(desc); ok {
			task.current = current
		} else {
			panic(fmt.Sprintf("task %s not found", desc))
		}
	}

}

func (p *Progress) SetTotal(total int64, description ...string) {
	if len(description) == 0 {
		p.model.tasks[0].total = total
		return
	}

	for _, desc := range description {
		if task, ok := p.findTask(desc); ok {
			task.total = total
		} else {
			panic(fmt.Sprintf("task %s not found", desc))
		}
	}

}

func (p *Progress) SetPercent(percent float64, description ...string) {
	if len(description) == 0 {
		p.model.tasks[0].percent = percent
		return
	}

	for _, desc := range description {
		if task, ok := p.findTask(desc); ok {
			task.percent = percent
		} else {
			panic(fmt.Sprintf("task %s not found", desc))
		}
	}

}
