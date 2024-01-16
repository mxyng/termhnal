# term-hn
Browse Hacker News from the Terminal. Built with [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea).

### TODO

- [ ] Text posts and comments
  - [ ] Word wrapping
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
go install github.com/mxyng/term-hn@main
```

## Key Maps

- <kbd>Ctrl+d</kbd> quit

### List View

- <kbd>1</kbd> Top stories
- <kbd>2</kbd> new stories
- <kbd>3</kbd> best stories
- <kbd>4</kbd> ask hn
- <kbd>5</kbd> show hn
- <kbd>6</kbd> jobs
- <kbd>k</kbd> <kbd>Up</kbd> up
- <kbd>j</kbd> <kbd>Down</kbd> down
- <kbd>l</kbd> <kbd>PageDown</kbd> next page
- <kbd>h</kbd> <kbd>PageUp</kbd> previous page
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>/</kbd> search
- <kbd>q</kbd> <kbd>Esc</kbd> quit
- <kbd>F5</kbd> refresh current list
- <kbd>?<kbd> help

### Story View

- <kbd>k</kbd> <kbd>Up</kbd> scroll up
- <kbd>j</kbd> <kbd>Down</kbd> scroll down
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>q</kbd> <kbd>Esc</kbd> back to list
