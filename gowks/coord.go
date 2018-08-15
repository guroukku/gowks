package gowks

import (
	"errors"
	"log"
	"math"
	"math/rand"
)

type coord struct {
	x, y int
}

type line struct {
	a, b coord
}

type node struct {
	p *node
	v coord
	c []*node
}

// c.sum(b) returns a coord that is the sum of coord b with coord c
func (c coord) sum(b coord) coord {
	return coord{c.x + b.x,c.y + b.y}
}

func (c coord) div(b coord) coord {
	return coord{c.x / b.x, c.y / b.y}
}

func (c coord) mul(b coord) coord {
	return coord{c.x * b.x, c.y * b.y}
}

func (c coord) eq(b coord) bool {
	return b.x == c.x && b.y == c.y
}

// c.neg() returns the negative of coord c
func (c coord) neg() coord {
	return coord{-c.x, -c.y}
}

func (c coord) flipY() coord {
	return coord{c.x, -c.y}
}

func (c coord) flipX() coord {
	return coord{-c.x, c.y}
}

// c.dist(b) returns a float64 that is the distance between coord c and b
func (c coord) dist(b coord) float64 {
	return math.Sqrt(float64((b.x - c.x) * (b.x - c.x) + (b.y - c.y) * (b.y - c.y)))
}

func (c coord) angle(b coord) float64 {
	var dx, dy int
	dy = c.y - b.y
	dx = c.x - b.x
	return math.Atan2(float64(-dy), float64(dx)) * (180 / math.Pi) - 180 + 360
}

func (c coord) grad(b coord) float64 {
	return float64(b.y - c.y + 1)/float64(b.x - c.x + 1)
}

func (c coord) onLine(b coord, d int) coord {
	return coord{(1 - d) * c.x + b.x * d, (1 - d) * c.y + b.y * d}
}

// c.inLine(a, b) returns true if coord c is on the line of coords a and b
func (c coord) inLine(a, b coord) bool {
	return c.dist(a) + c.dist(b) == a.dist(b)
}

// c.inBounds(b) returns true if coord c is within bounds of (-b, b)
// c.inBounds(a, b) returns true if coord c is within bounds of (a, b)
func (c coord) inBounds(b ...coord) bool {
	var b1, b2 coord
	switch len(b) {
	case 1:
		b1, b2 = b[0].neg(), b[0]
	case 2:
		b1, b2 = b[0], b[1]
	default:
		return false
	}
	return (b1.x <= c.x && c.x <= b2.x || b2.x <= c.x && c.x <= b1.x) && (b1.y <= c.y && c.y <= b2.y || b2.y <= c.y && c.y <= b1.y)
}

func (c coord) choosePoint(d float64, b ...coord) (coord, error) {
	var b1, b2, p coord
	switch len(b) {
	case 1:
		b1, b2 = b[0].neg(), b[0]
	case 2:
		b1, b2 = b[0], b[1]
	default:
		return coord{0,0}, errors.New("Too many coords to define bounds.")
	}
	if c.inBounds(b1, b2) {
		var h float64
		p = c
		for c.sum(p).inBounds(b1, b2) == false || p == c {
			h = float64(rand.Intn(b2.x))
			p.x = int(math.Cos(math.Pi * d / 180) * h)
			p.y = int(math.Sin(math.Pi * d / 180) * h)
		}
		return c.sum(p), nil
	}  else {
		return coord{0,0}, errors.New("Coord is out of given bounds.")
	}

}

func (l line) throughBounds(a, b coord) bool {
	return (l.a.x <= a.x && a.x <= l.b.x || l.b.x <= a.x && a.x <= l.a.x ||
	l.a.y <= a.y && a.y <= l.b.y || l.b.y <= a.y && a.y <= l.a.y) &&
	(l.a.inBounds(a,b) || l.b.inBounds(a,b))
}

func (n *node) genRandCardinalTree(b coord) {
	var c []*node
	d := []float64{0,90,180,270}
	for i := rand.Intn(3); i <= len(d); i++ {
		a := rand.Intn(len(d))
		d[a] = d[len(d)-1]
		d = d[:len(d)-1]
	}
	for i := 0; i < len(d); i++ {
		cv, e := n.v.choosePoint(d[i], b)
		if e != nil {
			log.Print(e)
		}
		cn := node{
			p:	n,
			v:	cv,
		}
		if rand.Intn(2) > 0 {
			cn.genRandCardinalTree(b)
		}
		c = append(c, &cn)
	}
	//log.Print("Complete tree:")
	//for _, e := range c {
		//log.Print("\tParent:", e.p.v)
		//log.Print("\tValue:", e.v)
	//}
	n.c = c
}

// intPow returns (a^b) as an int
func intPow(a, b int) int {
	n := 1
	for i := 0; i < b; i++ {
		n *= a
	}
	return n
}

