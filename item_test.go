package main

import (
	"sync"
	"testing"
)

// TestCommentConcurrency exercises the real concurrency shape: bubbletea runs
// Update/Render on one goroutine while each Cmd runs on its own goroutine and
// calls AddComment. Run with -race.
func TestCommentConcurrency(t *testing.T) {
	story := NewStory(0)
	for i := range 32 {
		story.Item.Kids = append(story.Item.Kids, i+1)
	}

	p := NewPaneView()

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			for j := range 32 {
				comment := NewComment(i*32 + j)
				comment.ID = i*32 + j
				comment.By = "user"
				comment.Text = "hello"
				story.AddComment(comment)
			}
		})
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 500 {
			_, _ = p.Update(ViewMsg[*Story]{Value: story})
		}
	}()

	wg.Wait()
	<-done
}

func TestAddCommentDedupe(t *testing.T) {
	item := &Item{}

	first := NewComment(0)
	first.ID = 42
	second := NewComment(1)
	second.ID = 42

	item.AddComment(first)
	item.AddComment(second)
	item.AddComment(nil)

	if len(item.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(item.Comments))
	}
	if item.Comments[0] != first {
		t.Fatal("expected first comment to be kept")
	}
	if !item.HasComment(42) {
		t.Fatal("expected HasComment(42) to be true")
	}
	if item.HasComment(43) {
		t.Fatal("expected HasComment(43) to be false")
	}
}
