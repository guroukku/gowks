package gowks

import (
	//"log"
	"math/rand"
)

type wall struct {
	root node
}

type maze struct {
	hunter *hunter
	bounds coord
	walls []*node
	bullets []*bullet
	gowks []*gowk
	nests []*nest
	pointsChan chan (int)
}

func newMaze(h *hunter, p chan (int), w []*node, b coord) *maze {
	m := &maze{
		hunter:		h,
		bounds:		b,
		walls:		w,
		pointsChan:	p,
		nests:		[]*nest{&nest{coord{-b.x,-b.y}, 0, true}, &nest{coord{b.x, b.y}, 0, true}, &nest{coord{b.x, -b.y}, 0, true}, &nest{coord{-b.x, b.y}, 0, true}},
	}

	return m
}

func (m *maze) addPoints(p int) {
	m.pointsChan <- p
}

func (m *maze) moveBullets() {
	for i, _ := range(m.bullets) {
		if m.bullets[i].alive {
			m.bullets[i].updatePosition()
			for ii, _ := range(m.walls) {
				if m.bullets[i].testCollision(m.walls[ii]) {
					m.bullets[i].alive = false
				}
			}
		}
	}
}

func (m *maze) spawnGowks() {
	for _, n := range(m.nests) {
		if n.alive {
			d := idx_dir[rand.Intn(7)]
			c, _ := n.loc.choosePoint(dir_deg[d], m.bounds)
			m.gowks = append(m.gowks, &gowk{n.loc.sum(dir_coord[d]), c, d, true, true, true, false, 0})
		}
	}
}
func (m *maze) updateAlive() {
	for i, j := 0, len(m.bullets); i < j; i++ {
		if !m.bullets[i].alive {
			m.bullets[i] = m.bullets[j-1]
			m.bullets[j-1] = nil
			m.bullets = m.bullets[:j-1]
			j--
		}
	}
	for i, j := 0, len(m.gowks); i < j; i++ {
		if !m.gowks[i].alive {
			m.gowks[i] = m.gowks[j-1]
			m.gowks[j-1] = nil
			m.gowks = m.gowks[:j-1]
			j--
		}
	}
	for i, j := 0, len(m.nests); i < j; i++ {
		if !m.nests[i].alive {
			m.nests[i] = m.nests[j-1]
			m.nests[j-1] = nil
			m.nests = m.nests[:j-1]
			j--
		}
	}
}


//for each wall, check line between wall and CNodes for coord
func (w *wall) isWall(c coord) bool {
	for _, wc := range w.root.c {
		if c.inLine(w.root.v, wc.v) {
			return true
		}
		//return wc.isWall(c)
	}
	return false
}

func (m *maze) isOccupied(c coord) bool {
	return m.hunter.isOnPosition(c)
}
