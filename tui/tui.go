package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"sort"
	"sync"
)

type Status struct {
	Name     string
	IP       string
	Location string
	Flag     string
	Results  string
}

type Model struct {
	Statuses map[string]Status
	table    *table.Table
	mu       *sync.Mutex
}

func NewModel(statuses map[string]Status, mu *sync.Mutex) *Model {
	var TableStyle = func(row, col int) lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("DNS Server", "DNS IP", "DNS Location", "DNS Query Result")

	return &Model{
		Statuses: statuses,
		table:    t,
		mu:       mu,
	}

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.String() {
		case "q":
			return m, tea.Quit
		}
	case string:
		switch msg {
		case "update":
			return m, nil
		case "quit":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders every time Update is called
func (m Model) View() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Extract keys and sort them
	keys := make([]string, 0, len(m.Statuses))
	for k := range m.Statuses {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate in sorted order
	m.table.ClearRows()
	for _, k := range keys {
		m.table.Row(m.Statuses[k].Name, m.Statuses[k].IP, m.Statuses[k].Flag+" "+m.Statuses[k].Location, m.Statuses[k].Results)
	}

	return m.table.String() + "\n"
}
