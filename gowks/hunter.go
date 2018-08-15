package gowks

import (
	"errors"
)

const (
	NORTH direction = 1 << iota
	SOUTH
	WEST
	EAST
)

var DEFAULT_HUNTER_1 = [4]cell{{'0', hunter1Attrib, bgColor},
				{'0', hunter1Attrib, bgColor},
				{'<', hunter1Attrib, bgColor},
				{'>', hunter1Attrib, bgColor}}
var DEFAULT_HUNTER_2 = [4]cell{{'@', hunter2Attrib, bgColor},
				{'@', hunter2Attrib, bgColor},
				{'~', hunter2Attrib, bgColor},
				{'~', hunter2Attrib, bgColor}}
var DEFAULT_HUNTER_3 = [4]cell{{'!', hunter3Attrib, bgColor},
				{'!', hunter3Attrib, bgColor},
				{'>', hunter3Attrib, bgColor},
				{'<', hunter3Attrib, bgColor}}


var GOWK = cell{'☺', gowkAttrib, bgColor}
var GOWK_ARROW = [8]cell{{'↑', gowkAttrib, bgColor},
	{'↖', gowkAttrib, bgColor},
	{'↗', gowkAttrib, bgColor},
	{'↓', gowkAttrib, bgColor},
	{'↙', gowkAttrib, bgColor},
	{'↘', gowkAttrib, bgColor},
	{'←', gowkAttrib, bgColor},
	{'→', gowkAttrib, bgColor}}

type direction int

type hunter struct {
	name string
	id int
	score int
	loc coord
	gfx [4]cell
}

type gowk struct {
	loc coord
	dst coord
	dir direction
	loaded bool
	moving bool
	alive bool
	ghost bool
	cooldown int
}

type nest struct {
	loc coord
	frame int
	alive bool
}

type bullet struct {
	loc coord
	dir direction
	frame bool
	alive bool
	isGowk bool
	hunter *hunter
}

func (b *bullet) updatePosition() {
	b.loc = b.loc.sum(dir_coord[b.dir])
}

func newHunter(c coord, n string) *hunter {
	return &hunter{
		loc:	c,
		name:	n,
		gfx:	DEFAULT_HUNTER_1,
	}
}

func (h *hunter) die() error {
	return errors.New("Died")
}

func (h *hunter) move(e keyboardEvent) {
	if e.key != 0 {
		h.loc = h.loc.sum(dir_coord[keyToDirection(e.key)])
	} else {
		h.loc = h.loc.sum(dir_coord[runeToDirection(e.ch)])
	}
}

func (h *hunter) fire(d direction) *bullet {
	return &bullet{h.loc.sum(dir_coord[d]), d, true, true, false, h}
}

func (h *hunter) isOnPosition(c coord) bool {
	return false
}

func (h *hunter) addPoints(p int) {
	h.score += p
}

func (g *gowk) move() {
	g.loc = g.loc.sum(dir_coord[g.dir])
}

func (g *gowk) fire() *bullet {
	g.loaded = false
	g.cooldown = 100
	return &bullet{g.loc.sum(dir_coord[g.dir]), g.dir, true, true, true, nil}
}

func (g *gowk) updateDirection() {
	g.dir = deg_dir[int(g.loc.angle(g.dst)/45) % 8]
}

func (b *bullet) testCollision(w *node) bool {
	for _, c := range(w.c) {
		if b.loc.inLine(w.v, c.v) {
			return true
		} else {
			b.testCollision(c)
		}
	}
	return false
}

func (b *bullet) checkHit(c coord) bool {
	return b.loc.eq(c)
}

func (b *bullet) checkHitHunter(c coord) bool {
	d, e, f := coord{c.x+1, c.y}, coord{c.x, c.y-1}, coord{c.x+1, c.y-1}
	return b.loc.eq(c) || b.loc.eq(d) || b.loc.eq(e) || b.loc.eq(f)
}

