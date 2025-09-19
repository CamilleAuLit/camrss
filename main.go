package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"strings"

	//	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Feed struct {
	title   string
	article string
}

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("200")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct {
	width    int
	height   int
	styles   *Styles
	articles []*gofeed.Item
}

func GetArticles(link string) ([]*gofeed.Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(link)
	if err != nil {
		return nil, err
	}
	return feed.Items, nil
}

func New(articles []*gofeed.Item) *model {
	styles := DefaultStyles()
	return &model{
		styles:   styles,
		articles: articles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			return m, nil
		}
	}
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var b strings.Builder

	for i, article := range m.articles {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%d. %s\n", i+1, article.Title)))
		b.WriteString("\n")
		b.WriteString(article.Content)
		b.WriteString("\n\n")

	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			b.String(),
		),
	)
}

func main() {
	feedURL := "https://hnrss.org/frontpage"
	items, err := GetArticles(feedURL)
	if err != nil {
		log.Fatal(err)
	}

	m := New(items)
	f, err := tea.LogToFile("debug.log", "debug")

	if err != nil {
		log.Fatalf("err: %W", err)
	}
	defer f.Close()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("err: %W", err)
	}
}
