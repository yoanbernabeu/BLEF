package viewer

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

type viewMode int

const (
	booksView viewMode = iota
	collectionsView
	statsView
	detailView
)

// Model represents the bubbletea model for the viewer
type Model struct {
	doc          *blef.BLEFDocument
	mode         viewMode
	selectedIdx  int
	filteredData []interface{}
	detailData   interface{}
	width        int
	height       int
	err          error
}

// NewModel creates a new viewer model
func NewModel(doc *blef.BLEFDocument) Model {
	model := Model{
		doc:         doc,
		mode:        booksView,
		selectedIdx: 0,
		width:       80,
		height:      24,
	}
	model.updateFilteredData()
	return model
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			m.mode = (m.mode + 1) % 3
			m.selectedIdx = 0
			m.updateFilteredData()
			return m, nil

		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case "down", "j":
			maxIdx := len(m.filteredData) - 1
			if m.selectedIdx < maxIdx {
				m.selectedIdx++
			}
			return m, nil

		case "enter":
			if m.mode == booksView && m.selectedIdx < len(m.filteredData) {
				m.detailData = m.filteredData[m.selectedIdx]
				m.mode = detailView
			}
			return m, nil

		case "esc":
			if m.mode == detailView {
				m.mode = booksView
				m.detailData = nil
			}
			return m, nil

		case "home", "g":
			m.selectedIdx = 0
			return m, nil

		case "end", "G":
			m.selectedIdx = len(m.filteredData) - 1
			return m, nil
		}
	}

	return m, nil
}

// View renders the view
func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	var content strings.Builder

	// Header
	content.WriteString(m.renderHeader())
	content.WriteString("\n\n")

	// Content based on mode
	switch m.mode {
	case booksView:
		content.WriteString(m.renderBooksList())
	case collectionsView:
		content.WriteString(m.renderCollectionsList())
	case statsView:
		content.WriteString(m.renderStats())
	case detailView:
		content.WriteString(m.renderDetail())
	}

	// Footer
	content.WriteString("\n\n")
	content.WriteString(m.renderFooter())

	return content.String()
}

// renderHeader renders the header with tabs
func (m Model) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		PaddingLeft(2)

	tabActiveStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("36")).
		Background(lipgloss.Color("235")).
		Padding(0, 2)

	tabInactiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 2)

	title := titleStyle.Render("📚 BLEF Viewer")

	booksTab := "Books"
	collectionsTab := "Collections"
	statsTab := "Stats"

	if m.mode == booksView {
		booksTab = tabActiveStyle.Render(booksTab)
	} else {
		booksTab = tabInactiveStyle.Render(booksTab)
	}

	if m.mode == collectionsView {
		collectionsTab = tabActiveStyle.Render(collectionsTab)
	} else {
		collectionsTab = tabInactiveStyle.Render(collectionsTab)
	}

	if m.mode == statsView {
		statsTab = tabActiveStyle.Render(statsTab)
	} else {
		statsTab = tabInactiveStyle.Render(statsTab)
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, booksTab, collectionsTab, statsTab)

	return title + "\n\n" + tabs
}

// renderFooter renders the footer with help
func (m Model) renderFooter() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingLeft(2)

	var help string
	if m.mode == detailView {
		help = "esc: back  q: quit"
	} else {
		help = "↑/↓: navigate  enter: details  tab: switch view  q: quit"
	}

	return helpStyle.Render(help)
}

// updateFilteredData updates the filtered data based on current mode
func (m *Model) updateFilteredData() {
	m.filteredData = make([]interface{}, 0)

	switch m.mode {
	case booksView:
		for i := range m.doc.Books {
			m.filteredData = append(m.filteredData, &m.doc.Books[i])
		}
	case collectionsView:
		for i := range m.doc.Collections {
			m.filteredData = append(m.filteredData, &m.doc.Collections[i])
		}
	}

	if m.selectedIdx >= len(m.filteredData) {
		m.selectedIdx = 0
	}
}

// renderBooksList renders the list of books
func (m Model) renderBooksList() string {
	if len(m.doc.Books) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("  No books found")
	}

	var content strings.Builder
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		PaddingLeft(2)

	// Show a window of books around the selected one
	start := m.selectedIdx - 5
	if start < 0 {
		start = 0
	}
	end := start + m.height - 10
	if end > len(m.doc.Books) {
		end = len(m.doc.Books)
	}

	for i := start; i < end; i++ {
		book := &m.doc.Books[i]
		status := m.getBookStatus(book.ID)
		statusEmoji := getStatusEmoji(status)

		line := fmt.Sprintf("%s %s", statusEmoji, truncate(book.Title, 60))
		if len(book.Authors) > 0 {
			line += fmt.Sprintf(" - %s", book.Authors[0].Name)
		}

		if i == m.selectedIdx {
			content.WriteString(selectedStyle.Render("> " + line))
		} else {
			content.WriteString(normalStyle.Render("  " + line))
		}
		content.WriteString("\n")
	}

	// Add scroll indicator
	if len(m.doc.Books) > (end - start) {
		content.WriteString(normalStyle.Render(fmt.Sprintf("\n  (%d/%d books)", m.selectedIdx+1, len(m.doc.Books))))
	}

	return content.String()
}

// renderCollectionsList renders the list of collections
func (m Model) renderCollectionsList() string {
	if len(m.doc.Collections) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("  No collections found")
	}

	var content strings.Builder
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		PaddingLeft(2)

	for i, coll := range m.doc.Collections {
		count := m.getCollectionBookCount(coll.ID)
		line := fmt.Sprintf("📚 %s (%d books)", coll.Name, count)

		if i == m.selectedIdx {
			content.WriteString(selectedStyle.Render("> " + line))
		} else {
			content.WriteString(normalStyle.Render("  " + line))
		}
		content.WriteString("\n")
	}

	return content.String()
}

// renderStats renders statistics
func (m Model) renderStats() string {
	var content strings.Builder
	style := lipgloss.NewStyle().PaddingLeft(2)

	content.WriteString(style.Render(fmt.Sprintf("Format: %s v%s\n", m.doc.Format, m.doc.Version)))
	content.WriteString(style.Render(fmt.Sprintf("Exported: %s\n\n", m.doc.ExportedAt.Format("2006-01-02 15:04:05"))))

	content.WriteString(style.Render(fmt.Sprintf("Total Books: %d\n", len(m.doc.Books))))
	content.WriteString(style.Render(fmt.Sprintf("Collections: %d\n", len(m.doc.Collections))))
	content.WriteString(style.Render(fmt.Sprintf("Entries: %d\n\n", len(m.doc.Entries))))

	// Status breakdown
	statusCount := make(map[string]int)
	for _, entry := range m.doc.Entries {
		statusCount[entry.UserData.Status]++
	}

	content.WriteString(style.Render("Reading Status:\n"))
	for status, count := range statusCount {
		emoji := getStatusEmoji(status)
		content.WriteString(style.Render(fmt.Sprintf("  %s %s: %d\n", emoji, status, count)))
	}

	return content.String()
}

// renderDetail renders book details
func (m Model) renderDetail() string {
	if m.detailData == nil {
		return ""
	}

	book, ok := m.detailData.(*blef.Book)
	if !ok {
		return ""
	}

	var content strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	valueStyle := lipgloss.NewStyle().PaddingLeft(2)

	content.WriteString(titleStyle.Render(book.Title))
	content.WriteString("\n\n")

	// Authors
	if len(book.Authors) > 0 {
		content.WriteString(labelStyle.Render("Authors:"))
		content.WriteString("\n")
		for _, author := range book.Authors {
			content.WriteString(valueStyle.Render(fmt.Sprintf("  %s", author.Name)))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	// Identifiers
	if book.Identifiers.ISBN13 != "" {
		content.WriteString(labelStyle.Render("ISBN-13: "))
		content.WriteString(book.Identifiers.ISBN13)
		content.WriteString("\n")
	}

	if book.Language != "" {
		content.WriteString(labelStyle.Render("Language: "))
		content.WriteString(book.Language)
		content.WriteString("\n")
	}

	// Get entry data
	entries := m.doc.GetEntriesForBook(book.ID)
	if len(entries) > 0 {
		entry := entries[0]
		content.WriteString("\n")
		content.WriteString(labelStyle.Render("Status: "))
		emoji := getStatusEmoji(entry.UserData.Status)
		content.WriteString(fmt.Sprintf("%s %s", emoji, entry.UserData.Status))
		content.WriteString("\n")

		if entry.UserData.Rating > 0 {
			content.WriteString(labelStyle.Render("Rating: "))
			content.WriteString(fmt.Sprintf("%.1f/5 ⭐", entry.UserData.Rating))
			content.WriteString("\n")
		}

		if entry.UserData.Review != "" {
			content.WriteString("\n")
			content.WriteString(labelStyle.Render("Review:"))
			content.WriteString("\n")
			content.WriteString(valueStyle.Render(wrapText(entry.UserData.Review, 70)))
			content.WriteString("\n")
		}
	}

	return content.String()
}

// Helper functions

func (m *Model) getBookStatus(bookID string) string {
	for _, entry := range m.doc.Entries {
		if entry.BookID == bookID {
			return entry.UserData.Status
		}
	}
	return ""
}

func (m *Model) getCollectionBookCount(collID string) int {
	count := 0
	for _, entry := range m.doc.Entries {
		for _, cid := range entry.CollectionIDs {
			if cid == collID {
				count++
				break
			}
		}
	}
	return count
}

func getStatusEmoji(status string) string {
	switch status {
	case "read":
		return "✅"
	case "reading":
		return "📖"
	case "to-read":
		return "📚"
	case "abandoned":
		return "❌"
	case "wishlist":
		return "⭐"
	default:
		return "📙"
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := "  " + words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = "  " + word
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}
