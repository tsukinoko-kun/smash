# smash

super magical awesome shell

Do **not** use this as your default shell (/bin/sh) because it is not POSIX compliant. Like fish, it is designed to be a user-friendly interactive shell. It is not designed to be a system shell that runs scripts.

Environment variables are loaded from your system default shell.

smash uses ANSI escape codes to format the prompt. If you are using a terminal that does not support ANSI escape codes, you will not be able to use smash.

There are some bugs on Windows because I don't use Windows. Feel free to open PRs to fix them.

## Features

- Substitute environment variables and ~
- Tab completion
- Configurable and interactive prompt
- Command history
- Command aliases
- Context prompt (e.g. git branch)
- Context aware completions

### Builtins

- `calc` (simple calculator)
- `cd` (change directory)
- `echo` (print arguments to stdout)
- `exit` (exit shell session)
- `printf` (print formatted string to stdout)
- `time` (time command execution)
- `zu` (similar to [ajeetdsouza/zoxide](https://github.com/ajeetdsouza/zoxide), change directory based on history and frequency)

## Upcoming

- More configuration options
- More builtins

## Configuration

`~/.config/smash/config.toml`

- `ps1`: Prompt shown on user input
- `ps2`: Prompt shown on the executed commands
- `alias`: Map of command aliases where keys can be a string or list of strings
- `color`: Map of colors supporting all [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) colors
  - `completion_text`: text color of completion list
  - `completion_selected_bg`: background color of selected completion entry
- `on_start`: List of commands to run on shell start

## Install

```shell
brew install tsukinoko-kun/tap/smash
```

```shell
go install github.com/tsukinoko-kun/smash@latest
```
