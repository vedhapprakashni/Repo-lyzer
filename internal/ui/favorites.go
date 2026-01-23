package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FavoritesModel struct {
	favorites *Favorites
	cursor    int
}

func NewFavoritesModel() FavoritesModel {
	return FavoritesModel{}
}

func (m FavoritesModel) Update(msg tea.Msg) (FavoritesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.favorites != nil && m.cursor < len(m.favorites.Items)-1 {
				m.cursor++
			}
		case "enter":
			// Analyze selected favorite
			if m.favorites != nil && len(m.favorites.Items) > 0 {
				repoName := m.favorites.Items[m.cursor].RepoName
				m.favorites.UpdateUsage(repoName)
				m.favorites.Save()
				return m, func() tea.Msg { return AnalyzeRepoMsg{repoName: repoName} }
			}
		case "d":
			// Remove from favorites
			if m.favorites != nil && len(m.favorites.Items) > 0 {
				m.favorites.Remove(m.favorites.Items[m.cursor].RepoName)
				m.favorites.Save()
				if m.cursor >= len(m.favorites.Items) && m.cursor > 0 {
					m.cursor--
				}
			}
		case "a":
			// Add new favorite (go to input)
			return m, func() tea.Msg { return SwitchToInputMsg{} }
		case "q", "esc":
			return m, func() tea.Msg { return BackToMenuMsg{} }
		}
	}
	return m, nil
}

func (m FavoritesModel) View(width, height int) string {
	header := TitleStyle.Render("⭐ Favorite Repositories")

	if m.favorites == nil || len(m.favorites.Items) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render("No favorites yet!\n\nAnalyze a repository and press 'b' to bookmark it."),
			SubtleStyle.Render("a: add new • q/ESC: back to menu"),
		)

		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	// Build favorites list
	var lines []string
	lines = append(lines, fmt.Sprintf("%-35s │ %-10s │ %s", "Repository", "Uses", "Last Used"))
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("─").Repeat(65))

	for i, fav := range m.favorites.Items {
		prefix := "  "
		if i == m.cursor {
			prefix = "▶ "
		}
		line := fmt.Sprintf("%s%-33s │ %-10d │ %s",
			prefix,
			fav.RepoName,
			fav.UseCount,
			fav.LastUsed.Format("2006-01-02"),
		)
		if i == m.cursor {
			lines = append(lines, SelectedStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	tableBox := BoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	footer := SubtleStyle.Render("↑↓: navigate • Enter: analyze • d: remove • a: add new • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
		footer,
	)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m *FavoritesModel) SetFavorites(favorites *Favorites) {
	m.favorites = favorites
}

func (m *FavoritesModel) GetCursor() int {
	return m.cursor
}

func (m *FavoritesModel) SetCursor(cursor int) {
	m.cursor = cursor
}
