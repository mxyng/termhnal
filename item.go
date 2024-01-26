package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/net/html"
)

type Item struct {
	Rank int

	By    string `json:"by"`
	ID    int    `json:"id"`
	Kids  []int  `json:"kids"`
	Time  int64  `json:"time"`
	Title string `json:"title"`
	Type  string `json:"type"`

	Comments []*Comment
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

func HTMLText(t string) string {
	root, err := html.Parse(strings.NewReader(html.UnescapeString(t)))
	if err != nil {
		panic(err)
	}

	var sb strings.Builder

	type state struct {
		attributes []html.Attribute
	}

	var fn func(*state, *html.Node)
	fn = func(s *state, n *html.Node) {
		switch n.Type {
		case html.TextNode:
			text := n.Data
			for _, attribute := range s.attributes {
				switch attribute.Key {
				case "href":
					text = fmt.Sprintf("(%s %s)", n.Data, attribute.Val)
					if n.Data == attribute.Val {
						text = fmt.Sprintf("%s", attribute.Val)
					} else if trim := strings.TrimSuffix(n.Data, "..."); strings.HasPrefix(attribute.Val, trim) {
						// HN truncates long links and appends "..."
						text = fmt.Sprintf("%s", attribute.Val)
					}

					continue
				}
			}

			sb.WriteString(text)
		case html.ElementNode:
			switch n.Data {
			case "a":
				s.attributes = append(s.attributes, n.Attr...)
			case "p":
				sb.WriteString("\n\n")
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			fn(s, child)
		}

		s.attributes = nil
	}

	fn(&state{}, root)
	return sb.String()
}

type Story struct {
	*Item

	Descendants int    `json:"descendants"`
	Score       int    `json:"score"`
	Text        string `json:"text"`
	URL         string `json:"url"`
}

func NewStory(Rank int) *Story {
	return &Story{
		Item: &Item{
			Rank: Rank,
		},
	}
}

func (s Story) FilterValue() string {
	return s.Title()
}

func (s Story) Title() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d. %s", s.Rank+1, s.Item.Title)

	if s.URL != "" {
		link, err := url.Parse(s.URL)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(&sb, " (%s)", link.Host)
	}

	return sb.String()
}

func (s Story) Description() string {
	var prefix string
	switch {
	case s.Rank < 10:
		prefix = strings.Repeat(" ", 3)
	case s.Rank < 100:
		prefix = strings.Repeat(" ", 4)
	default:
		prefix = strings.Repeat(" ", 5)
	}

	return fmt.Sprintf("%s%d points by %s %s | %d comments", prefix, s.Score, s.By, humanize(time.Unix(s.Time, 0)), s.Descendants)
}

type Comment struct {
	*Item

	Parent int    `json:"parent"`
	Text   string `json:"text"`
}

func NewComment(rank int) *Comment {
	return &Comment{
		Item: &Item{
			Rank: rank,
		},
	}
}

func (c Comment) Title() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6600"))
	return fmt.Sprintf(
		"%s %s",
		style.Render(c.By),
		style.Copy().Faint(true).Render(humanize(time.Unix(c.Time, 0))),
	)
}

func (c Comment) Description() string {
	return strings.TrimSpace(HTMLText(c.Text))
}
