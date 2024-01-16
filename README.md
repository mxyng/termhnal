# termhnal
Browse Hacker News from the Terminal. Built with [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea).

### TODO

- [ ] Text posts and comments
  - [x] Word wrapping
  - [x] Convert HTML to Markdown
    - [ ] Handle HTML tags
- [x] Story list switching (new, ask, show, etc.)
  - [ ] Display non-active story lists
- [ ] Save stories to read later
- [ ] Footer and help in story view
- [ ] Replace viewport with custom model

## Install

Requires Go 1.20 or higher.

```shell
go install github.com/mxyng/termhnal@main
```

## Key Maps

- <kbd>Ctrl+d</kbd> quit

### List View

- <kbd>1</kbd> top
- <kbd>2</kbd> new
- <kbd>3</kbd> best
- <kbd>4</kbd> ask
- <kbd>5</kbd> show
- <kbd>6</kbd> jobs
- <kbd>k</kbd> <kbd>Up</kbd> up
- <kbd>j</kbd> <kbd>Down</kbd> down
- <kbd>h</kbd> <kbd>Left</kbd> <kbd>PageUp</kbd> previous page
- <kbd>l</kbd> <kbd>Right</kbd> <kbd>PageDown</kbd> next page
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>/</kbd> search
- <kbd>q</kbd> <kbd>Esc</kbd> quit
- <kbd>F5</kbd> refresh current list
- <kbd>?</kbd> help

### Story View

- <kbd>k</kbd> <kbd>Up</kbd> scroll up
- <kbd>j</kbd> <kbd>Down</kbd> scroll down
- <kbd>h</kbd> <kbd>Left</kbd> <kbd>PageUp</kbd> previous page
- <kbd>l</kbd> <kbd>Right</kbd> <kbd>PageDown</kbd> next page
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>q</kbd> <kbd>Esc</kbd> back to list
