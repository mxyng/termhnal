package main

import (
	"cmp"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

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

	mu sync.RWMutex
}

func (i *Item) AddComment(c *Comment) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.Comments = append(i.Comments, c)
	slices.SortFunc(i.Comments, func(i, j *Comment) int {
		return cmp.Compare(i.Rank, j.Rank)
	})
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
		element    string
		attributes map[string]string
	}

	var fn func(state, *html.Node)
	fn = func(s state, n *html.Node) {
		switch n.Type {
		case html.TextNode:
			text := n.Data
			switch s.element {
			case "a":
				if val, ok := s.attributes["href"]; ok {
					text = fmt.Sprintf("(%s %s)", n.Data, val)
					if n.Data == val {
						text = val
					} else if trim := strings.TrimSuffix(n.Data, "..."); strings.HasPrefix(val, trim) {
						// HN truncates long links and appends "..."
						text = val
					}
				}
			}

			sb.WriteString(text)
		case html.ElementNode:
			switch n.Data {
			case "a":
				s.element = n.Data
				s.attributes = make(map[string]string)
				for _, attr := range n.Attr {
					// discard Namespace
					s.attributes[attr.Key] = attr.Val
				}
			case "p":
				sb.WriteString("\n\n")
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			fn(s, child)
		}
	}

	fn(state{}, root)
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
