package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.Screen = ScreenExecuting
			m.ExecIndex = 0
			m.ExecDone = false
			m.ExecResults = nil
			m.Viewport = newViewport(m.Width, m.Height-4)
			m.ViewportReady = true
			return m, m.execNext()
		case "esc":
			m.Screen = ScreenSelect
			return m, nil
		}
	}
	return m, nil
}

func (m Model) viewConfirm() string {
	tabs := InactiveTab.Render("Select") + " " + ActiveTab.Render("Confirm")
	summary := m.formatPendingOps()

	hasDelete := false
	for _, op := range m.PendingOps {
		if op.Op == OpDelete {
			hasDelete = true
			break
		}
	}

	var warn string
	if hasDelete {
		warn = "\n\n" + WarningBanner.Render("WARNING: Deletion is permanent and cannot be undone!")
	}

	return tabs + "\n\n" + summary + warn + "\n\n  Enter to execute, Esc to go back"
}