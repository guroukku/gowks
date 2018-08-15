package gowks

import "github.com/nsf/termbox-go"

type keyboardEventType int

const (
	MOVE keyboardEventType = 1 + iota
	FIRE
	RESIZE
	RETRY
	END
)

type keyboardEvent struct {
	eventType	keyboardEventType
	key		termbox.Key
	ch		rune
}

func keyToDirection(k termbox.Key) direction {
	switch k {
	case termbox.KeyArrowUp:
		return NORTH
	case termbox.KeyArrowDown:
		return SOUTH
	case termbox.KeyArrowLeft:
		return WEST
	case termbox.KeyArrowRight:
		return EAST
	default:
		return 0
	}
}

func runeToDirection(r rune) direction {
	switch r {
	case '7':
		return NORTH | WEST
	case '8':
		return NORTH
	case '9':
		return NORTH | EAST
	case '6':
		return EAST
	case '3':
		return SOUTH | EAST
	case '2':
		return SOUTH
	case '1':
		return SOUTH | WEST
	case '4':
		 return WEST
	case 'q':
		return NORTH | WEST
	case 'w':
		return NORTH
	case 'e':
		return NORTH | EAST
	case 'd':
		return EAST
	case 'c':
		return SOUTH | EAST
	case 'x':
		return SOUTH
	case 'z':
		return SOUTH | WEST
	case 'a':
		return WEST
	default:
		return 0
	}
}

// listenToKeyboard() processes the output of termbox.PollEvent() to the provided event channel
func listenToKeyboard(evChan chan keyboardEvent) {
	termbox.SetInputMode(termbox.InputEsc)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if termbox.KeyArrowRight <= ev.Key && ev.Key <= termbox.KeyArrowUp {
				evChan <- keyboardEvent{eventType: MOVE, key: ev.Key}
			} else if ev.Key == termbox.KeyCtrlQ {
				evChan <- keyboardEvent{eventType: END, key: ev.Key}
			} else if ev.Key == termbox.KeyCtrlR {
				evChan <- keyboardEvent{eventType: RETRY, key: ev.Key}
			} else if ev.Key == 0 {
				if '1' <= ev.Ch && ev.Ch <= '9' {
					evChan <- keyboardEvent{eventType: MOVE, ch: ev.Ch}
				} else if ev.Ch == 'q' || ev.Ch == 'w' || ev.Ch == 'e' || ev.Ch == 'd' || ev.Ch == 'c' || ev.Ch == 'x' || ev.Ch == 'z' || ev.Ch == 'a'{
					evChan <- keyboardEvent{eventType: FIRE, ch: ev.Ch}
				}
			}
		case termbox.EventResize:
			evChan <- keyboardEvent{eventType: RESIZE, key: ev.Key}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
