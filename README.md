<p align="center">
  <img alt="GitHub License" src="https://img.shields.io/github/license/mxyng/termhnal">
  <img alt="GitHub Contributors" src="https://img.shields.io/github/contributors/mxyng/termhnal">
  <img alt="GitHub top language" src="https://img.shields.io/github/languages/top/mxyng/termhnal">
  <img alt="GitHub Actions Workflow Status" src="https://img.shields.io/github/actions/workflow/status/mxyng/termhnal/build.yaml">
  <img alt="GitHub Issues" src="https://img.shields.io/github/issues/mxyng/termhnal">
  <img alt="GitHub Pull Requests" src="https://img.shields.io/github/issues-pr/mxyng/termhnal">
  <img alt="Docker Image Version" src="https://img.shields.io/docker/v/mxyng/termhnal/latest">
  <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/mxyng/termhnal">
  <img alt="Docker Image Size" src="https://img.shields.io/docker/image-size/mxyng/termhnal">
  <img alt="GitHub Downloads" src="https://img.shields.io/github/downloads/mxyng/termhnal/total">
  <img alt="GitHub Stargazers" src="https://img.shields.io/github/stars/mxyng/termhnal">
</p>

Browse Hacker News in a Terminal. Built with [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea).

## 👷 Install

Requires Go 1.27 or higher.

```shell
go install github.com/mxyng/termhnal@main
```

### :whale: Container Image

```shell
docker run -it ghcr.io/mxyng/termhnal
```

```shell
docker run -it mxyng/termhnal
```

## :keyboard: Key Maps

- <kbd>Ctrl+d</kbd> quit

### :notebook: List View

- <kbd>1</kbd> top
- <kbd>2</kbd> new
- <kbd>3</kbd> best
- <kbd>4</kbd> ask
- <kbd>5</kbd> show
- <kbd>6</kbd> jobs
- <kbd>Enter</kbd> open story
- <kbd>k</kbd> <kbd>Up</kbd> up
- <kbd>j</kbd> <kbd>Down</kbd> down
- <kbd>h</kbd> <kbd>Left</kbd> <kbd>PageUp</kbd> previous page
- <kbd>l</kbd> <kbd>Right</kbd> <kbd>PageDown</kbd> next page
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>/</kbd> search
- <kbd>Esc</kbd> clear search / quit
- <kbd>Tab</kbd> focus header

### :book: Story View

- <kbd>k</kbd> <kbd>Up</kbd> scroll up
- <kbd>j</kbd> <kbd>Down</kbd> scroll down
- <kbd>u</kbd> half page up
- <kbd>d</kbd> half page down
- <kbd>b</kbd> <kbd>PageUp</kbd> previous page
- <kbd>f</kbd> <kbd>Space</kbd> <kbd>PageDown</kbd> next page
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>Tab</kbd> focus header
- <kbd>Esc</kbd> <kbd>Backspace</kbd> back to list

### :bookmark: Header

Focused with <kbd>Tab</kbd> or by scrolling past the top with <kbd>k</kbd>/<kbd>Up</kbd>.

- <kbd>h</kbd> <kbd>Left</kbd> previous
- <kbd>l</kbd> <kbd>Right</kbd> next
- <kbd>Enter</kbd> select
- <kbd>j</kbd> <kbd>Down</kbd> <kbd>Tab</kbd> back
