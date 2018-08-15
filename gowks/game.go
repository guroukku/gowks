package gowks

import (
	"encoding/json"
	//"errors"
	//"log"
	"math/rand"
	"os"
	"time"
	"github.com/nsf/termbox-go"
)

var (
	pointsChan = make(chan int)
	keyboardEventsChan = make(chan keyboardEvent)
)

type Game struct {
	maze *maze
	score int
	isOver bool
	msgs []*msg
}

func initialHunter() *hunter {
	return newHunter(coord{x: 1, y: 1}, "Hunter")
}

func initialScore() int {
	return 0
}

func initialMaze() *maze {
	return newMaze(initialHunter(), pointsChan, initialWalls(coord{100,50}), coord{100,100})
}

// initialWalls returns a []*node containing a random number of root walls with generated nodes.
func initialWalls(b coord) []*node {
	var iw []*node
	for i := rand.Intn(5); i <= 6; i++ {
		var rw node
		rw.v = coord{rand.Intn(b.x * 2) - b.x, rand.Intn(b.y * 2) - b.y}
		//log.Print("New wall root:", rw.v)
		rw.genRandCardinalTree(b)
		iw = append(iw, &rw)
	}
	return iw
}

func (g *Game) end() {
	g.isOver = true
}

func (g *Game) tick() time.Duration {
	return time.Duration(1000/30)* time.Millisecond
}

func (g *Game) retry() {
	g.maze = initialMaze()
	g.score = initialScore()
	g.msgs = nil
	g.isOver = false
}

func (g *Game) addPoints(p int) {
	g.score += p
}

func (g *Game) saveMaze() error {
	return nil
}

// NewGame returns a pointer to a new game object.
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	return &Game{maze: initialMaze(), score: initialScore()}
}

// Start begins the mainloop.
func (g *Game) Start(){
	b, _ := json.Marshal(g.maze)

	os.Stdout.Write(b)

	if err := termbox.Init(); err != nil {
		panic(err)
	}

	defer termbox.Close()
	t := term{
		h: g.maze.hunter,
	}
	//t.resize()
	//term.x, term.y = termbox.Size()

	go listenToKeyboard(keyboardEventsChan)

	if err := g.render(&t); err != nil {
		panic(err)
	}

mainloop:
	for {
		select {
		case p := <-pointsChan:
			g.addPoints(p)
		case e := <-keyboardEventsChan:
			switch e.eventType {
			case MOVE:
				if !g.isOver{
					g.maze.hunter.move(e)
				}
			case FIRE:
				if !g.isOver {
					g.maze.bullets = append(g.maze.bullets, g.maze.hunter.fire(runeToDirection(e.ch)))
				}
			case RESIZE:
			case RETRY:
				g.retry()
				t.h = g.maze.hunter
			case END:
				break mainloop
			}
		default:
			if err := g.render(&t); err != nil {
				panic(err)
			}

			if !g.isOver {
				g.maze.moveBullets()

				if rand.Intn(100) > 98 {
					g.maze.spawnGowks()
				}

				for _, gwk := range(g.maze.gowks) {
					//if rand.Intn(10) > 8 {
						gwk.dst = g.maze.hunter.loc
						//g.maze.gowks[i].dst, _ = g.maze.gowks[i].loc.choosePoint(dir_deg[idx_dir[rand.Intn(8)]], g.maze.bounds)
					//}
					gwk.updateDirection()
					for _, b := range(g.maze.bullets) {
						if b.checkHit(gwk.loc) {
							b.alive = false
							gwk.alive = false
							if b.hunter != nil {
								b.hunter.addPoints(1)
							}
						}
						if b.checkHitHunter(g.maze.hunter.loc) {
							_ = g.maze.hunter.die()
							g.end()
							g.msgs = append(g.msgs, &msg{t.size.div(coord{2,2}).sum(coord{-3,0}), termbox.ColorRed | termbox.AttrBold, bgColor, "YOU LOSE!"})
							//break mainloop
						}
						for _, n := range(g.maze.nests) {
							if n.alive {
								if b.hunter != nil {
									if b.checkHitHunter(n.loc) {
										n.alive = false
										b.hunter.addPoints(10)
									}
								}
							}
						}
					}

					if rand.Intn(15) > 10 {
						gwk.moving = false
					}

					if gwk.loaded {
						if rand.Intn(50) > 45 {
							g.maze.bullets = append(g.maze.bullets, gwk.fire())
						}
					} else {
						gwk.cooldown--
						if gwk.cooldown == 0 {
							gwk.loaded = true
						}
					}

					if gwk.moving {
						gwk.move()
					} else if rand.Intn(15) > 12 {
						gwk.moving = true

					}
				}
			g.maze.updateAlive()
			if len(g.maze.nests) == 0 && len(g.maze.gowks) == 0 {
				g.isOver = true
				g.msgs = append(g.msgs, &msg{t.size.div(coord{2,2}).sum(coord{-3,0}), termbox.ColorYellow | termbox.AttrBold, bgColor, "YOU WIN!"})
			}
		}
		time.Sleep(g.tick())
	}}
}

