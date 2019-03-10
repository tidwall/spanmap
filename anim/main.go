package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/tidwall/spanmap"
)

var speed float64 = 1

func main() {
	hideCursor()
	defer showCursor()
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		moveTo(9)
		showCursor()
		os.Exit(0)
	}()
	for {
		clearScreen()
		moveTo(1)
		min, max := uint64(11), uint64(26)
		var m spanmap.Map
		fmt.Printf("\033[1mVisualization of how SpanMap stores items.\033[0m")
		moveTo(5)
		fmt.Printf("\033[2m[ ] Set 16 items starting at index 11\033[0m\n")
		fmt.Printf("\033[2m[ ] Delete item at index 13\033[0m\n")
		fmt.Printf("\033[2m[ ] Set item at index 29, where 13 use to be\033[0m\n")
		moveTo(3)
		fmt.Printf("\033[2m[ -- ]\033[0m")
		moveTo(9)
		sleep(1)
		for i := min; i <= max; i++ {
			m.Set(i, i)
			moveTo(3)
			fmt.Printf("%s", colorItem(mapString(&m, false), i, "\033[32m\033[1m"))
			moveTo(5)
			fmt.Printf("\033[2m[ ]\033[0m\033[1m Set %d items starting at index %d\033[0m", max-min+1, min)
			moveTo(9)
			sleep(0.5)
		}
		moveTo(3)
		fmt.Printf("%s", mapString(&m, false))
		moveTo(5)
		fmt.Printf("[✓] Set 16 items starting at index 11")
		moveTo(6)
		fmt.Printf("\033[2m[ ]\033[0m\033[1m Delete item at index %d\033[0m", 13)
		moveTo(9)

		// delete item
		moveTo(3)
		fmt.Printf("%s", colorItem(mapString(&m, false), 13, "\033[31m\033[1m"))
		moveTo(9)
		sleep(2)
		m.Delete(13)
		moveTo(3)
		fmt.Printf("%s", colorItem(mapString(&m, false), 13, "\033[31m\033[1m"))

		moveTo(6)
		fmt.Printf("[✓] Delete item at index 13")
		moveTo(9)
		sleep(1)
		moveTo(7)
		fmt.Printf("\033[2m[ ]\033[0m\033[1m Set item at index 29, where 13 use to be\033[0m")

		// overset item
		m.Set(29, 29)
		moveTo(3)
		fmt.Printf("%s", colorItem(mapString(&m, false), 29, "\033[32m\033[1m"))
		moveTo(9)
		sleep(2)
		moveTo(3)
		fmt.Printf("%s", mapString(&m, false))
		moveTo(7)
		fmt.Printf("[✓] Set item at index 29, where 13 use to be")
		moveTo(9)
		sleep(4)
	}
}

func colorItem(line string, index uint64, color string) string {
	return strings.Replace(line,
		fmt.Sprintf(" %02d ", index),
		fmt.Sprintf(" %s%02d\033[0m ", color, index),
		-1,
	)
}

func sleep(seconds float64) {
	time.Sleep(time.Duration(seconds / speed * float64(time.Second)))
}

func hideCursor() {
	fmt.Printf("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func clearScreen() {
	fmt.Printf("\033[2J")
}

func moveTo(line int) {
	fmt.Printf("\033[%d;0H", line)
}

type mapItem struct {
	index uint64
	item  interface{}
}

// Map is a map that is optimized for data that span contiguous indexes.
type tMap struct {
	min, max uint64
	items    []mapItem
}

func mapString(om *spanmap.Map, showBounds bool) string {
	m := (*tMap)(unsafe.Pointer(om))
	var out []byte
	if showBounds {
		out = append(out, fmt.Sprintf("%02d-%02d ", m.min, m.max)...)
	}
	out = append(out, "[ "...)
	for i := 0; i < len(m.items); i++ {
		if m.items[i].item == nil {
			out = append(out, fmt.Sprintf("-- ")...)
		} else {
			out = append(out, fmt.Sprintf("%02d ", m.items[i].index)...)
		}
	}
	out = append(out, ']')
	return string(out)
}
