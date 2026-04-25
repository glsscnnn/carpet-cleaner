package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Operation int

const (
	OpDelete Operation = iota
	OpRename
)

type Screen int

const (
	ScreenLoading Screen = iota
	ScreenSelect
	ScreenRenameModal
	ScreenConfirm
	ScreenExecuting
	ScreenDone
)

type RepoItem struct {
	Name    string
	Owner   string
	Desc    string
	Private bool
	Selected    bool
}

func (r RepoItem) Title() string {
	if r.Selected {
		return fmt.Sprintf("%s %s/%s", CheckboxChecked.Render("[x]"), r.Owner, r.Name)
	}
	return fmt.Sprintf("%s %s/%s", CheckboxUnchecked.Render("[ ]"), r.Owner, r.Name)
}

func (r RepoItem) Description() string {
	d := r.Desc
	if d == "" {
		d = "(No description)"
	}
	var badge string
	if r.Private {
		badge = PrivateBadge.Render("Private")
	} else {
		badge = PublicBadge.Render("Public")
	}
	return fmt.Sprintf("[%s] %s", badge, d)
}

func (r RepoItem) FilterValue() string { return r.Owner + "/" + r.Name }

type PendingOp struct {
	Repo    RepoItem
	NewName string
	Op      Operation
}

type ExecResult struct {
	Repo    string
	Op      Operation
	NewName string
	Success bool
	Error   string
}

type Model struct {
	Screen Screen

	Owner string
	Token string

	Repos []RepoItem

	Selected  map[string]bool
	Operation Operation

	Spinner   spinner.Model
	LoadErr   string

	List      list.Model
	ListReady bool

	RenameIndex    int
	RenameRepoKeys []string
	RenameInput    textinput.Model
	RenameCurrent  string
	RenameMap      map[string]string

	PendingOps []PendingOp

	Viewport       viewport.Model
	ViewportReady  bool
	ExecResults    []ExecResult
	ExecIndex      int
	ExecDone       bool

	Width  int
	Height int
}

func NewModel(owner, token string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SuccessMarker

	ti := textinput.New()
	ti.Placeholder = "New repo name"
	ti.CharLimit = 100
	ti.Width = 40

	return Model{
		Screen:      ScreenLoading,
		Owner:       owner,
		Token:       token,
		Selected:    make(map[string]bool),
		RenameMap:   make(map[string]string),
		Operation:   OpDelete,
		Spinner:     s,
		RenameInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, fetchRepos(m.Owner, m.Token))
}

type fetchedReposMsg struct {
	repos []RepoItem
	err   error
}

func fetchRepos(owner, token string) tea.Cmd {
	return func() tea.Msg {
		repos, err := fetchAllRepos(owner, token)
		return fetchedReposMsg{repos: repos, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if m.ListReady {
			m.List.SetSize(msg.Width, msg.Height-4)
		}
		if m.ViewportReady {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - 4
		}
		return m, nil

	case fetchedReposMsg:
		if msg.err != nil {
			m.LoadErr = msg.err.Error()
			return m, nil
		}
		m.Repos = msg.repos
		m.Screen = ScreenSelect
		cmd := m.initList()
		return m, cmd
	}

	switch m.Screen {
	case ScreenLoading:
		return m.updateLoading(msg)
	case ScreenSelect:
		return m.updateSelect(msg)
	case ScreenRenameModal:
		return m.updateRenameModal(msg)
	case ScreenConfirm:
		return m.updateConfirm(msg)
	case ScreenExecuting:
		return m.updateExecuting(msg)
	case ScreenDone:
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "q" {
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	switch m.Screen {
	case ScreenLoading:
		return m.viewLoading()
	case ScreenSelect:
		return m.viewSelect()
	case ScreenRenameModal:
		return m.viewSelect() + "\n" + m.viewRenameModal()
	case ScreenConfirm:
		return m.viewConfirm()
	case ScreenExecuting, ScreenDone:
		return m.viewExecuting()
	}
	return ""
}

func (m *Model) initList() tea.Cmd {
	items := make([]list.Item, len(m.Repos))
	for i, r := range m.Repos {
		items[i] = r
	}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#FFFFFF")).
		BorderForeground(lipgloss.Color("#7D56F4"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#CCCCCC")).
		BorderForeground(lipgloss.Color("#7D56F4"))
	m.List = list.New(items, delegate, m.Width, m.Height-4)
	m.List.Title = "Repositories"
	m.List.SetShowStatusBar(false)
	m.List.SetShowHelp(false)
	m.List.FilterInput.Prompt = "/"
	m.ListReady = true
	return nil
}

func (m Model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.LoadErr != "" {
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "q" {
			return m, tea.Quit
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

func (m Model) viewLoading() string {
	if m.LoadErr != "" {
		return fmt.Sprintf("\n  Error: %s\n\n  Press q to quit.", m.LoadErr)
	}
	return fmt.Sprintf("\n  %s Fetching repositories...", m.Spinner.View())
}

func selectedRepos(m Model) []RepoItem {
	var result []RepoItem
	for _, r := range m.Repos {
		if m.Selected[r.Owner+"/"+r.Name] {
			result = append(result, r)
		}
	}
	return result
}

func repoKey(r RepoItem) string {
	return r.Owner + "/" + r.Name
}