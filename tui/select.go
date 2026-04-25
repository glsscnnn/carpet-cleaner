package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			m.Operation = OpDelete
			return m, nil
		case "r":
			m.Operation = OpRename
			return m, nil
		case " ":
			if m.List.FilterState() == list.Filtering {
				break
			}
			i, ok := m.List.SelectedItem().(RepoItem)
			if ok {
				k := repoKey(i)
				m.Selected[k] = !m.Selected[k]
				for idx, r := range m.Repos {
					if repoKey(r) == k {
						m.Repos[idx].Selected = m.Selected[k]
						break
					}
				}
				m.updateListItems()
			}
			return m, nil
		case "enter":
			if m.List.FilterState() == list.Filtering {
				break
			}
			if len(m.Selected) > 0 {
				if m.Operation == OpRename {
					return m.startRenameModal()
				}
				return m.gotoConfirm()
			}
		case "tab":
			if len(m.Selected) > 0 {
				if m.Operation == OpRename {
					return m.startRenameModal()
				}
				return m.gotoConfirm()
			}
		case "esc":
			if m.List.FilterState() == list.Filtering {
				break
			}
			if len(m.Selected) == 0 {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) viewSelect() string {
	opStr := "delete"
	if m.Operation == OpRename {
		opStr = "rename"
	}
	status := StatusBar.Render(fmt.Sprintf("%d selected  |  operation: %s  |  d: delete  r: rename  /: filter  space: toggle  enter: proceed", len(m.Selected), opStr))
	tabs := ActiveTab.Render("Select") + " " + InactiveTab.Render("Confirm")
	return tabs + "\n" + m.List.View() + "\n" + status
}

func (m *Model) updateListItems() {
	items := make([]list.Item, len(m.Repos))
	for i, r := range m.Repos {
		items[i] = r
	}
	m.List.SetItems(items)
}

func (m Model) gotoConfirm() (tea.Model, tea.Cmd) {
	m.buildPendingOps()
	m.Screen = ScreenConfirm
	return m, nil
}

func (m Model) startRenameModal() (tea.Model, tea.Cmd) {
	sel := selectedRepos(m)
	m.RenameRepoKeys = make([]string, len(sel))
	for i, r := range sel {
		m.RenameRepoKeys[i] = repoKey(r)
	}
	m.RenameIndex = 0
	m.RenameCurrent = m.RenameRepoKeys[0]
	for _, r := range m.Repos {
		if repoKey(r) == m.RenameCurrent {
			m.RenameInput.SetValue(r.Name)
			break
		}
	}
	cmd := m.RenameInput.Focus()
	m.Screen = ScreenRenameModal
	return m, cmd
}

func (m Model) getRepoByKey(k string) RepoItem {
	for _, r := range m.Repos {
		if repoKey(r) == k {
			return r
		}
	}
	return RepoItem{}
}

func (m *Model) buildPendingOps() {
	m.PendingOps = nil
	for _, r := range m.Repos {
		if !m.Selected[repoKey(r)] {
			continue
		}
		k := repoKey(r)
		if m.Operation == OpRename {
			newName := r.Name
			if n, ok := m.RenameMap[k]; ok {
				newName = n
			}
			m.PendingOps = append(m.PendingOps, PendingOp{Repo: r, NewName: newName, Op: OpRename})
		} else {
			m.PendingOps = append(m.PendingOps, PendingOp{Repo: r, Op: OpDelete})
		}
	}
}

func (m Model) formatPendingOps() string {
	var lines []string
	for _, op := range m.PendingOps {
		if op.Op == OpDelete {
			lines = append(lines, fmt.Sprintf("  Delete: %s", repoKey(op.Repo)))
		} else {
			lines = append(lines, fmt.Sprintf("  Rename: %s → %s", repoKey(op.Repo), op.NewName))
		}
	}
	return strings.Join(lines, "\n")
}