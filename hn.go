package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/segmentio/textio"
	"golang.org/x/net/html"
)

// ref: https://github.com/HackerNews/API
type HN struct {
	baseURL *url.URL
}

func NewHN() *HN {
	baseURL, err := url.Parse("https://hacker-news.firebaseio.com/v0")
	if err != nil {
		panic(err)
	}

	return &HN{baseURL: baseURL}
}

func (h *HN) Top() ([]int, error) {
	return h.items("top")
}

func (h *HN) New() ([]int, error) {
	return h.items("new")
}

func (h *HN) Best() ([]int, error) {
	return h.items("best")
}

func (h *HN) Ask() ([]int, error) {
	return h.items("ask")
}

func (h *HN) Show() ([]int, error) {
	return h.items("show")
}

func (h *HN) Job() ([]int, error) {
	return h.items("job")
}

func (h *HN) items(kind string) ([]int, error) {
	requestURL := h.baseURL.JoinPath(fmt.Sprintf("/%sstories.json", kind))
	request, err := http.NewRequestWithContext(context.Background(), "GET", requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var stories []int
	if err := json.NewDecoder(response.Body).Decode(&stories); err != nil {
		return nil, err
	}

	return stories, nil
}

func (h *HN) item(id int, item any) error {
	requestURL := h.baseURL.JoinPath(fmt.Sprintf("/item/%d.json", id))
	request, err := http.NewRequestWithContext(context.Background(), "GET", requestURL.String(), nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(item)
}

func (h *HN) Story(rank, id int) (*Story, error) {
	var story Story
	if err := h.item(id, &story); err != nil {
		return nil, err
	}

	story.rank = rank
	return &story, nil
}

func (h *HN) Comment(rank, id int) (*Comment, error) {
	var comment Comment
	if err := h.item(id, &comment); err != nil {
		return nil, err
	}

	comment.rank = rank
	return &comment, nil
}

type Item struct {
	rank int

	By    string `json:"by"`
	ID    int    `json:"id"`
	Kids  []int  `json:"kids"`
	Time  int64  `json:"time"`
	Title string `json:"title"`
	Type  string `json:"type"`

	Comments []*Comment
}

type Story struct {
	*Item

	Descendants int    `json:"descendants"`
	Score       int    `json:"score"`
	Text        string `json:"text"`
	URL         string `json:"url"`
}

func (s Story) Rank() int {
	return s.rank + 1
}

func (s Story) FilterValue() string {
	return s.Title()
}

func (s Story) Title() string {
	link, err := url.Parse(s.URL)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%d. %s (%s)", s.Rank(), s.Item.Title, link.Host)
}

func (s Story) Description() string {
	var prefix string
	switch {
	case s.rank < 10:
		prefix = strings.Repeat(" ", 3)
	case s.rank < 100:
		prefix = strings.Repeat(" ", 4)
	default:
		prefix = strings.Repeat(" ", 5)
	}

	return fmt.Sprintf("%s%d points by %s %s | %d comments", prefix, s.Score, s.By, humanize(time.Unix(s.Time, 0)), s.Descendants)
}

func (s Story) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s\n", titleStyle.Render(s.Item.Title))
	fmt.Fprintln(&sb, descriptionStyle.Render(strings.TrimSpace(s.Description())))

	if s.URL != "" {
		fmt.Fprintln(&sb, linkStyle.Render(s.URL))
	} else if s.Text != "" {
		fmt.Fprintln(&sb)
		pw := textio.NewPrefixWriter(&sb, "> ")
		fmt.Fprintln(pw, htmlString(s.Text))
	}

	for _, comment := range s.Comments {
		sb.WriteString(comment.String())
	}

	fmt.Fprintln(&sb)

	return sb.String()
}

type Comment struct {
	*Item

	Parent int    `json:"parent"`
	Text   string `json:"text"`
}

func (c Comment) Rank() int {
	return c.rank
}

func (c Comment) FilterValue() string {
	return c.Text
}

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})

var descriptionStyle = titleStyle.Copy().
	Foreground(lipgloss.AdaptiveColor{Light: "#a49fa5", Dark: "#777777"})

var linkStyle = descriptionStyle.Copy().
	Italic(true).
	Underline(true)

func (c Comment) Title() string {
	return titleStyle.Render(fmt.Sprintf("%s %s", c.By, humanize(time.Unix(c.Time, 0))))
}

func (c Comment) Description() string {
	return htmlString(c.Text)
}

func (c Comment) String() string {
	var sb strings.Builder
	if c.Text != "" {
		fmt.Fprintln(&sb)

		pw := textio.NewPrefixWriter(&sb, "│ ")
		fmt.Fprintln(pw, c.Title())
		fmt.Fprintln(pw, descriptionStyle.Render(c.Description()))

		for _, comment := range c.Comments {
			fmt.Fprint(pw, comment.String())
		}
	}

	return sb.String()
}

func htmlString(t string) string {
	root, err := html.Parse(strings.NewReader(html.UnescapeString(t)))
	if err != nil {
		log.Fatal(err)
	}

	var sb strings.Builder

	var fn func(*html.Node)
	fn = func(node *html.Node) {
		switch node.Type {
		case html.TextNode:
			sb.WriteString(node.Data)
		case html.ElementNode:
			switch node.Data {
			case "p":
				sb.WriteString("\n\n")
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			fn(child)
		}
	}

	fn(root)
	return sb.String()
}

func humanize(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < time.Hour*24:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
}
