package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

const (
	EditorName = "CAFFEETERIA"
	Version    = "0.2.0"
)

type Editor struct {
	screen      tcell.Screen
	lines       []string
	cursorX     int
	cursorY     int
	fileName    string
	commandMode bool
	commandLine string
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	e := &Editor{
		screen: s,
		lines:  []string{""},
	}

	for {
		e.draw()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if e.commandMode {
				e.handleCommandKey(ev)
			} else {
				if ev.Key() == tcell.KeyCtrlX {
					s.Fini()
					os.Exit(0)
				}
				if ev.Key() == tcell.KeyCtrlS {
					e.saveFile(e.fileName)
					continue
				}
				if ev.Key() == tcell.KeyCtrlO {
					e.commandMode = true
					e.commandLine = "open "
					continue
				}
				if ev.Rune() == ':' {
					e.commandMode = true
					e.commandLine = ""
					continue
				}
				e.handleKey(ev)
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func (e *Editor) draw() {
	e.screen.Clear()
	w, h := e.screen.Size()

	headerStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)
	fname := e.fileName
	if fname == "" {
		fname = "[New File]"
	}
	headerText := fmt.Sprintf(" %s v%s | File: %s | :cmd, Ctrl+S: Save", EditorName, Version, fname)
	for i := 0; i < w; i++ {
		e.screen.SetContent(i, 0, ' ', nil, headerStyle)
	}
	e.puts(0, 0, headerStyle, headerText)

	for y, line := range e.lines {
		if y+1 < h-1 {
			e.puts(0, y+1, tcell.StyleDefault, line)
		}
	}

	if e.commandMode {
		footerStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
		for i := 0; i < w; i++ {
			e.screen.SetContent(i, h-1, ' ', nil, footerStyle)
		}
		e.puts(0, h-1, footerStyle, ":"+e.commandLine)
		e.screen.ShowCursor(runewidth.StringWidth(e.commandLine)+1, h-1)
	} else {
		e.screen.ShowCursor(e.cursorX, e.cursorY+1)
	}

	e.screen.Show()
}

func (e *Editor) puts(x, y int, style tcell.Style, str string) {
	for _, r := range str {
		sw := runewidth.RuneWidth(r)
		if sw == 0 {
			sw = 1
		}
		e.screen.SetContent(x, y, r, nil, style)
		x += sw
	}
}

func (e *Editor) handleKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEnter:
		line := e.lines[e.cursorY]
		afterCursor := line[e.cursorX:]
		e.lines[e.cursorY] = line[:e.cursorX]
		newLines := append(e.lines[:e.cursorY+1], append([]string{afterCursor}, e.lines[e.cursorY+1:]...)...)
		e.lines = newLines
		e.cursorY++
		e.cursorX = 0
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if e.cursorX > 0 {
			line := e.lines[e.cursorY]
			e.lines[e.cursorY] = line[:e.cursorX-1] + line[e.cursorX:]
			e.cursorX--
		} else if e.cursorY > 0 {
			prevLineLen := len(e.lines[e.cursorY-1])
			e.lines[e.cursorY-1] += e.lines[e.cursorY]
			e.lines = append(e.lines[:e.cursorY], e.lines[e.cursorY+1:]...)
			e.cursorY--
			e.cursorX = prevLineLen
		}
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			if e.cursorX > len(e.lines[e.cursorY]) {
				e.cursorX = len(e.lines[e.cursorY])
			}
		}
	case tcell.KeyDown:
		if e.cursorY < len(e.lines)-1 {
			e.cursorY++
			if e.cursorX > len(e.lines[e.cursorY]) {
				e.cursorX = len(e.lines[e.cursorY])
			}
		}
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		} else if e.cursorY > 0 {
			e.cursorY--
			e.cursorX = len(e.lines[e.cursorY])
		}
	case tcell.KeyRight:
		if e.cursorX < len(e.lines[e.cursorY]) {
			e.cursorX++
		} else if e.cursorY < len(e.lines)-1 {
			e.cursorY++
			e.cursorX = 0
		}
	case tcell.KeyRune:
		line := e.lines[e.cursorY]
		e.lines[e.cursorY] = line[:e.cursorX] + string(ev.Rune()) + line[e.cursorX:]
		e.cursorX++
	}
}

func (e *Editor) handleCommandKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEnter:
		e.executeCommand()
		e.commandMode = false
	case tcell.KeyEsc:
		e.commandMode = false
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(e.commandLine) > 0 {
			e.commandLine = e.commandLine[:len(e.commandLine)-1]
		} else {
			e.commandMode = false
		}
	case tcell.KeyRune:
		e.commandLine += string(ev.Rune())
	}
}

func (e *Editor) executeCommand() {
	parts := strings.Fields(e.commandLine)
	if len(parts) == 0 {
		return
	}
	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "open":
		if len(args) > 0 {
			content, err := ioutil.ReadFile(args[0])
			if err == nil {
				e.lines = strings.Split(string(content), "\n")
				e.fileName = args[0]
				e.cursorX, e.cursorY = 0, 0
			}
		}
	case "new":
		if len(args) > 0 {
			e.lines = []string{""}
			e.fileName = args[0]
			e.cursorX, e.cursorY = 0, 0
		}
	case "save":
		target := e.fileName
		if len(args) > 0 {
			target = args[0]
		}
		e.saveFile(target)
	case "file_txt":
		e.generateTree()
	}
}

func (e *Editor) saveFile(name string) {
	if name == "" {
		return
	}
	content := strings.Join(e.lines, "\n")
	ioutil.WriteFile(name, []byte(content), 0644)
	e.fileName = name
}

func (e *Editor) generateTree() {
	var tree []string
	tree = append(tree, "--- Project Tree ---")
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		depth := strings.Count(path, string(os.PathSeparator))
		indent := strings.Repeat("  ", depth)
		icon := "📄 "
		if info.IsDir() {
			icon = "📁 "
		}
		tree = append(tree, indent+icon+filepath.Base(path))
		return nil
	})
	e.lines = tree
	e.cursorX, e.cursorY = 0, 0
}
