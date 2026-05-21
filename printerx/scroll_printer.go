package printerx

import (
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/term"
)

type ScrollPrinter struct {
	window  int
	buffer  []string
	mu      sync.Mutex
	isTTY   bool
	out     io.Writer
	lastLen int
}

func NewScrollPrinter(out io.Writer, window int) *ScrollPrinter {
	isTTY := false
	if f, ok := out.(*os.File); ok {
		isTTY = term.IsTerminal(int(f.Fd())) //nolint: gosec // for pass
	}

	return &ScrollPrinter{
		window: window,
		buffer: make([]string, 0, window),
		isTTY:  isTTY,
		out:    out,
	}
}

func (sp *ScrollPrinter) PrintLines(lines ...string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	for _, line := range lines {
		sp.buffer = append(sp.buffer, line)
		if len(sp.buffer) > sp.window {
			sp.buffer = sp.buffer[1:]
		}
	}

	if sp.isTTY {
		for range sp.lastLen {
			_, _ = fmt.Fprint(sp.out, "\033[F")
		}

		for _, l := range sp.buffer {
			_, _ = fmt.Fprintln(sp.out, l)
		}
		sp.lastLen = len(sp.buffer)
	} else {
		for _, l := range lines {
			_, _ = fmt.Fprintln(sp.out, l)
		}
	}
}
