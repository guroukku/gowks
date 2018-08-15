package gowks

import (
	"fmt"
//	"log"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	defaultColor = termbox.ColorDefault
	bgColor      = termbox.ColorDefault
	hunter1Attrib   = termbox.ColorWhite | termbox.AttrBold
	hunter2Attrib	= termbox.ColorRed | termbox.AttrBold
	hunter3Attrib	= termbox.ColorBlue | termbox.AttrBold
	wallAttrib = termbox.ColorBlue
	bulletAttribA = termbox.ColorBlue | termbox.AttrBold
	bulletAttribB = termbox.ColorYellow | termbox.AttrBold
	gowkAttrib = termbox.ColorGreen | termbox.AttrBold

	//box-drawing characters
	WE =	'═'
	NS =	'║'
	NW =	'╝'
	NE =	'╚'
	SW =	'╗'
	SE =	'╔'
	SWE =	'╦'
	NSE =	'╠'
	NSW =	'╣'
	NWE =	'╩'
	NSWE =	'╬'
)

var attrBold = []termbox.Attribute{
termbox.ColorWhite | termbox.AttrBold,
termbox.ColorYellow | termbox.AttrBold,
termbox.ColorRed | termbox.AttrBold,
termbox.ColorMagenta | termbox.AttrBold,
termbox.ColorBlue | termbox.AttrBold,
termbox.ColorCyan | termbox.AttrBold,
termbox.ColorGreen | termbox.AttrBold}

var (	dir_runes = map[direction]rune{
		WEST | EAST: WE,
		NORTH | SOUTH: NS,
		NORTH | WEST: NW,
		NORTH | EAST: NE,
		SOUTH | WEST: SW,
		SOUTH | EAST: SE,
		SOUTH | WEST | EAST: SWE,
		NORTH | SOUTH | EAST: NSE,
		NORTH | SOUTH | WEST: NSW,
		NORTH | WEST | EAST: NWE,
		NORTH | SOUTH | WEST | EAST: NSWE,
	}

	dir_coord = map[direction]coord{
		WEST: coord{-1,0},
		EAST: coord{1,0},
		NORTH: coord{0,1},
		SOUTH: coord{0,-1},
		NORTH | WEST: coord{-1, 1},
		NORTH | EAST: coord{1, 1},
		SOUTH | WEST: coord{-1, -1},
		SOUTH | EAST: coord{1, -1},
	}

	dir_idx = map[direction]int{
		NORTH: 0,
		NORTH | WEST: 1,
		NORTH | EAST: 2,
		SOUTH: 3,
		SOUTH | WEST: 4,
		SOUTH | EAST: 5,
		WEST: 6,
		EAST: 7,
	}

	dir_deg = map[direction]float64{
		NORTH: 270,
		NORTH | WEST: 225,
		NORTH | EAST: 315,
		SOUTH: 90,
		SOUTH | WEST: 135,
		SOUTH | EAST: 45,
		WEST: 180,
		EAST: 0,
	}

	deg_dir = map[int]direction{
		0: EAST,
		1: SOUTH | EAST,
		2: SOUTH,
		3: SOUTH | WEST,
		4: WEST,
		5: NORTH | WEST,
		6: NORTH,
		7: NORTH | EAST,
	}

	idx_dir = map[int]direction{
		0: NORTH,
		1: NORTH | WEST,
		2: NORTH | EAST,
		3: SOUTH,
		4: SOUTH | WEST,
		5: SOUTH | EAST,
		6: WEST,
		7: EAST,
	}
)

type term struct {
	off	coord
	size	coord
	loc	coord
	h	*hunter
	dirty	bool
}

type msg struct {
	loc coord
	fg termbox.Attribute
	bg termbox.Attribute
	str string
}

type cell struct {
	ch	rune
	fg	termbox.Attribute
	bg	termbox.Attribute
}

func (t *term) resize() {
	t.size.x, t.size.y = termbox.Size()
}

func (t *term) update() {
	t.loc = t.h.loc.sum(t.size.div(coord{2,2}).flipX())
}

func (t *term) setCell(x, y int, ch rune, fg, bg  termbox.Attribute) {
	x = x - t.loc.x
	y = t.loc.y - y
	c := coord{x, y}
	if c.inBounds(t.off, t.size){
		termbox.SetCell(x, y, ch, fg, bg)
	} else {
		return
	}
}

func (g *Game) render(t *term) error {
	termbox.Clear(defaultColor, defaultColor)
	t.resize()
	t.update()
	t.renderMaze(g.maze)
	t.renderTitle()
	t.renderHunter(t.h)
	t.renderScore(g.maze.hunter.score)
	t.renderQuitMessage()
	t.renderMessages(g.msgs)
	return termbox.Flush()
}

func (t *term) renderMessages(msgs []*msg) {
	for _, m := range(msgs) {
		tbprint(m.loc.x, m.loc.y, m.fg, m.bg, m.str)
	}
}

func (t *term) renderBullets(b_s []*bullet) {
	for _, b := range(b_s) {
		if b.alive {
			if b.isGowk {
				g := GOWK_ARROW[dir_idx[b.dir]]
				t.setCell(b.loc.x, b.loc.y, g.ch, g.fg, g.bg)
			} else if b.frame {
				t.setCell(b.loc.x, b.loc.y, '*', bulletAttribA, bgColor)
			} else {
				t.setCell(b.loc.x, b.loc.y, '◦', bulletAttribB, bgColor)
			}
			b.frame = !b.frame
		}
	}
}

func (t *term) renderHunter(h *hunter) {
	for i, r := range(h.gfx) {
		t.setCell(h.loc.x + i % 2, h.loc.y - (i/2), r.ch, r.fg, r.bg)
	}
}

func (t *term) renderGowk(g *gowk) {
	t.setCell(g.loc.x, g.loc.y, GOWK.ch, GOWK.fg, GOWK.bg)
	if g.loaded {
		bc := dir_coord[g.dir].sum(g.loc)
		t.setCell(bc.x, bc.y, GOWK_ARROW[dir_idx[g.dir]].ch, GOWK.fg, GOWK.bg)
	}
}

func (t *term) renderMaze(m *maze) {
	for _, w := range(m.walls) {
		t.renderWall(w)
	}
	t.renderBullets(m.bullets)
	for _, g := range(m.gowks) {
		t.renderGowk(g)
	}
	for _, n := range(m.nests) {
		t.renderNest(n)
	}
	for i := 1; i < t.size.y-1; i++ {
		termbox.SetCell(0, i, '│', defaultColor, bgColor)
		termbox.SetCell(t.size.x-1, i, '│', defaultColor, bgColor)
	}
	termbox.SetCell(0, 1, '┌', defaultColor, bgColor)
	termbox.SetCell(0, t.size.y-2, '└', defaultColor, bgColor)
	termbox.SetCell(t.size.x-1, 1, '┐', defaultColor, bgColor)
	termbox.SetCell(t.size.x-1, t.size.y-2, '┘', defaultColor, bgColor)

	fill(1, 1, t.size.x-2, 1, termbox.Cell{Ch: '─'})
	fill(1, t.size.y-2, t.size.x-2, 1, termbox.Cell{Ch: '─'})
}

func (w *node) directions() direction {
	var d direction
	for _, wn := range(w.c) {
		if wn.v.y == w.v.y {
			if wn.v.x > w.v.x {
				d = d | EAST
			} else {
				d = d | WEST
			}
		} else if wn.v.y > w.v.y {
			d = d | SOUTH
		} else {
			d = d | NORTH
		}
	}
	return d
}

func (t *term) renderNest(n *nest) {
	t.setCell(n.loc.x, n.loc.y, '┏', attrBold[n.frame], bgColor)
	t.setCell(n.loc.x+1, n.loc.y, '┓', attrBold[n.frame], bgColor)
	t.setCell(n.loc.x, n.loc.y-1, '┗', attrBold[n.frame], bgColor)
	t.setCell(n.loc.x+1, n.loc.y-1, '┛', attrBold[n.frame], bgColor)
	n.frame++
	if n.frame >= 6 {
		n.frame = 0
	}
}

func (t *term) renderWall(w *node) {
	var r rune
	for _, wn := range(w.c) {
		if w.v.y == wn.v.y {
			r = WE
		} else {
			r = NS
		}
		t.drawLine(w.v, wn.v, termbox.Cell{Ch: r, Fg: wallAttrib})
		t.renderWall(wn)
	}
	t.setCell(w.v.x, w.v.y, dir_runes[w.directions()], wallAttrib, bgColor)
}

func (t *term) renderScore(s int) {
	score := fmt.Sprintf("Score: %v", s)
	loc := fmt.Sprintf("h: %d, t: %d, t.size: %d", t.h.loc, t.loc, t.size)
	tbprint(0, t.size.y - 1, defaultColor, defaultColor, loc)
	tbprint(t.size.x - len(score), 0, defaultColor, defaultColor, score)
}

func (t *term) renderQuitMessage() {
	m := "CTRL-Q to quit"
	tbprint(t.size.x-17, t.size.y - 1, defaultColor, defaultColor, m)
}

func (t *term) renderTitle() {
	tbprint(t.size.x / 2, 0, defaultColor, defaultColor, "Gowks")
}

func (t *term) drawLine(start, end coord, cell termbox.Cell) {
	d := int(start.dist(end))
	if d < 0 {
		d = d * -1
	}
	p := start
	for l := 1; l <= d; l++ {
		p = start.onLine(end, int(l))
		t.setCell(p.x, p.y, cell.Ch, cell.Fg, cell.Bg)
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
