package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateRenameModal(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.RenameInput, cmd = m.RenameInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			newName := strings.TrimSpace(m.RenameInput.Value())
			if newName != "" {
				m.RenameMap[m.RenameCurrent] = newName
			}
			m.RenameIndex++
			if m.RenameIndex >= len(m.RenameRepoKeys) {
				m.buildPendingOps()
				m.Screen = ScreenConfirm
				return m, nil
			}
			m.RenameCurrent = m.RenameRepoKeys[m.RenameIndex]
			for _, r := range m.Repos {
				if repoKey(r) == m.RenameCurrent {
					m.RenameInput.SetValue(r.Name)
					break
				}
			}
			return m, nil
		case "esc":
			m.RenameIndex++
			if m.RenameIndex >= len(m.RenameRepoKeys) {
				if len(m.RenameMap) == 0 {
					m.Screen = ScreenSelect
					return m, nil
				}
				m.buildPendingOps()
				m.Screen = ScreenConfirm
				return m, nil
			}
			m.RenameCurrent = m.RenameRepoKeys[m.RenameIndex]
			for _, r := range m.Repos {
				if repoKey(r) == m.RenameCurrent {
					m.RenameInput.SetValue(r.Name)
					break
				}
			}
			return m, nil
		}
	}

	return m, cmd
}

func (m Model) viewRenameModal() string {
	content := fmt.Sprintf(
		"Rename repository: %s\n\nNew name:\n%s\n\nEnter to confirm, Esc to skip",
		m.RenameCurrent,
		m.RenameInput.View(),
	)
	box := ModalBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, box)
}