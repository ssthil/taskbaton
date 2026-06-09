package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

var colorEnabled = detectColor()

func detectColor() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
	ansiGray   = "\x1b[90m"
)

func paint(code, s string) string {
	if !colorEnabled {
		return s
	}
	return code + s + ansiReset
}

func bold(s string) string   { return paint(ansiBold, s) }
func dim(s string) string    { return paint(ansiDim, s) }
func green(s string) string  { return paint(ansiGreen, s) }
func yellow(s string) string { return paint(ansiYellow, s) }
func cyan(s string) string   { return paint(ansiCyan, s) }
func gray(s string) string   { return paint(ansiGray, s) }

const (
	glyphOK    = "✓"
	glyphArrow = "→"
	glyphDot   = "•"
	glyphWarn  = "▲"
)

func runeLen(s string) int { return utf8.RuneCountInString(s) }

func boxTop(w io.Writer, width int) {
	fmt.Fprintln(w, gray("╭"+strings.Repeat("─", width)+"╮"))
}

func boxBottom(w io.Writer, width int) {
	fmt.Fprintln(w, gray("╰"+strings.Repeat("─", width)+"╯"))
}

func boxRow(w io.Writer, width int, plain string, color func(string) string) {
	pad := width - runeLen(plain)
	if pad < 0 {
		pad = 0
	}
	s := plain + strings.Repeat(" ", pad)
	if color != nil {
		s = color(s)
	}
	fmt.Fprintf(w, "%s%s%s\n", gray("│"), s, gray("│"))
}

func success(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, "%s %s\n", green(glyphOK), fmt.Sprintf(format, a...))
}

func info(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, "%s %s\n", cyan(glyphArrow), fmt.Sprintf(format, a...))
}

func warn(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, "%s %s\n", yellow(glyphWarn), fmt.Sprintf(format, a...))
}

func note(w io.Writer, format string, a ...any) {
	fmt.Fprintln(w, gray("  "+fmt.Sprintf(format, a...)))
}

func step(w io.Writer, n int, cmd, desc string) {
	fmt.Fprintf(w, "  %s %s %s\n",
		cyan(fmt.Sprintf("%d.", n)), bold(fmt.Sprintf("%-22s", cmd)), gray(desc))
}

func tip(w io.Writer, cmd, desc string) {
	fmt.Fprintf(w, "  %s %s %s\n",
		gray(glyphDot), bold(fmt.Sprintf("%-22s", cmd)), gray(desc))
}
