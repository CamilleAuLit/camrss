package main

import (
	"github.com/mmcdole/gofeed"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Feed struct {
	title   string
	article *gofeed.Item
}

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("201")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct {
	index       int
	questions   []Question
	width       int
	height      int
	answerField textinput.Model
	styles      *Styles
	article     []Feed
}

type Question struct {
	question string
	answer   string
}

func NewQuestion(question string) Question {
	return Question{question: question}
}

func NewFeed(link string) (*Feed, error) {
	var a *gofeed.Item
	t := ""
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(link)
	if err != nil {
		return nil, err
	}
	for _, i := range feed.Items {
		t = feed.Title
		a = i
	}
	return &Feed{title: t, article: a}, nil
}

func New(questions []Question, articles []Feed) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Your answer here"
	answerField.Focus()

	return &model{
		questions:   questions,
		answerField: answerField,
		styles:      styles,
		article:     articles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := m.questions[m.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			current.answer = m.answerField.Value()
			m.answerField.SetValue("")
			m.Next()
			return m, nil
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var articleArticle string

	if len(m.article) > 0 {
		articleArticle = m.article[0].article.Content
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			articleArticle,
		),
	)
}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func main() {
	questions := []Question{
		NewQuestion("What is your name?"),
		NewQuestion("What is your favorite editor?"),
		NewQuestion("What is your wisdom?"),
	}

	feedURL := "https://hnrss.org/frontpage"
	myFeed, err := NewFeed(feedURL)
	if err != nil {
		log.Fatal(err)
	}

	articles := []Feed{*myFeed}

	m := New(questions, articles)
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
