package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	story := NewStory(rank)
	if err := h.item(id, story); err != nil {
		return nil, err
	}

	return story, nil
}

func (h *HN) Comment(rank, id int) (*Comment, error) {
	comment := NewComment(rank)
	if err := h.item(id, &comment); err != nil {
		return nil, err
	}

	return comment, nil
}
