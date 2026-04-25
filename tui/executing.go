package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type execDoneMsg struct {
	result ExecResult
}

func (m Model) execNext() tea.Cmd {
	if m.ExecIndex >= len(m.PendingOps) {
		return nil
	}
	op := m.PendingOps[m.ExecIndex]
	return func() tea.Msg {
		var err error
		if op.Op == OpDelete {
			err = deleteRepo(op.Repo.Owner, op.Repo.Name, m.Token)
		} else {
			err = renameRepo(op.Repo.Owner, op.Repo.Name, op.NewName, m.Token)
		}
		result := ExecResult{
			Repo:    repoKey(op.Repo),
			Op:      op.Op,
			NewName: op.NewName,
			Success: err == nil,
		}
		if err != nil {
			result.Error = err.Error()
		}
		return execDoneMsg{result: result}
	}
}

func (m Model) updateExecuting(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case execDoneMsg:
		m.ExecResults = append(m.ExecResults, msg.result)
		m.ExecIndex++
		if m.ExecIndex >= len(m.PendingOps) {
			m.ExecDone = true
			m.Screen = ScreenDone
		} else {
			cmd := m.execNext()
			m.refreshViewport()
			return m, cmd
		}
		m.refreshViewport()
		return m, nil
	}
	var cmd tea.Cmd
	if m.ViewportReady {
		m.Viewport, cmd = m.Viewport.Update(msg)
	}
	return m, cmd
}

func (m *Model) refreshViewport() {
	if !m.ViewportReady {
		return
	}
	m.Viewport.SetContent(m.formatExecResults())
}

func (m Model) formatExecResults() string {
	var lines []string
	tabs := InactiveTab.Render("Select") + " " + InactiveTab.Render("Confirm") + " " + ActiveTab.Render("Executing")
	lines = append(lines, tabs, "")

	for _, r := range m.ExecResults {
		if r.Op == OpRename {
			if r.Success {
				lines = append(lines, fmt.Sprintf("  %s Renamed %s → %s", SuccessMarker.Render("✓"), r.Repo, r.NewName))
			} else {
				lines = append(lines, fmt.Sprintf("  %s Failed to rename %s → %s: %s", ErrorMarker.Render("✗"), r.Repo, r.NewName, r.Error))
			}
		} else {
			if r.Success {
				lines = append(lines, fmt.Sprintf("  %s Deleted %s", SuccessMarker.Render("✓"), r.Repo))
			} else {
				lines = append(lines, fmt.Sprintf("  %s Failed to delete %s: %s", ErrorMarker.Render("✗"), r.Repo, r.Error))
			}
		}
	}

	// Show in-progress operation
	if !m.ExecDone && m.ExecIndex < len(m.PendingOps) {
		op := m.PendingOps[m.ExecIndex]
		if op.Op == OpRename {
			lines = append(lines, fmt.Sprintf("  ... Renaming %s → %s", repoKey(op.Repo), op.NewName))
		} else {
			lines = append(lines, fmt.Sprintf("  ... Deleting %s", repoKey(op.Repo)))
		}
	}

	if m.ExecDone {
		deleted := 0
		renamed := 0
		failed := 0
		for _, r := range m.ExecResults {
			if r.Success {
				if r.Op == OpDelete {
					deleted++
				} else {
					renamed++
				}
			} else {
				failed++
			}
		}
		lines = append(lines, "")
		summary := fmt.Sprintf("%d deleted, %d renamed, %d failed. Press q to quit.", deleted, renamed, failed)
		lines = append(lines, SummaryLine.Render(summary))
	}

	return strings.Join(lines, "\n") + "\n"
}

func (m Model) viewExecuting() string {
	if !m.ViewportReady {
		return ""
	}
	m.refreshViewport()
	return m.Viewport.View()
}

func newViewport(w, h int) viewport.Model {
	v := viewport.New(w, h)
	return v
}