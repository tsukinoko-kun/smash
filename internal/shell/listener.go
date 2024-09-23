package shell

import (
	"github.com/chzyer/readline"
)

type Listener struct {
	rl *readline.Instance
}

func (l *Listener) SetReadline(rl *readline.Instance) {
	l.rl = rl
}

func (l *Listener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	return line, pos, false
}
