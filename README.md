<p align="center">
  <img alt="Termhnal Logo" src="https://github.com/user-attachments/assets/12071c2e-0aa0-453d-88e4-cc5517ddc01a">
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

Requires Go 1.20 or higher.

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
- <kbd>k</kbd> <kbd>Up</kbd> up
- <kbd>j</kbd> <kbd>Down</kbd> down
- <kbd>h</kbd> <kbd>Left</kbd> <kbd>PageUp</kbd> previous page
- <kbd>l</kbd> <kbd>Right</kbd> <kbd>PageDown</kbd> next page
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>/</kbd> search
- <kbd>q</kbd> <kbd>Esc</kbd> quit

### :book: Story View

- <kbd>k</kbd> <kbd>Up</kbd> scroll up
- <kbd>j</kbd> <kbd>Down</kbd> scroll down
- <kbd>h</kbd> <kbd>Left</kbd> <kbd>PageUp</kbd> previous page
- <kbd>l</kbd> <kbd>Right</kbd> <kbd>PageDown</kbd> next page
- <kbd>g</kbd> <kbd>Home</kbd> go to start
- <kbd>Shift+g</kbd> <kbd>End</kbd> go to end
- <kbd>q</kbd> <kbd>Esc</kbd> back to list
