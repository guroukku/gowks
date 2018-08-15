package gowks

import(
	//"fmt"
	//"testing"
)

/*
var choosePointTests = []struct {
	in_a coord
	in_b int
	out bool
}{
	{coord{0,0},0,true},
	{coord{0,0},1,true},
	{coord{0,0},2,true},
	{coord{0,0},3,true},
}
func TestChoosePoint(t *testing.T) {
	for _, cpt := range choosePointTests {
		t.Run(fmt.Sprintf("%d, %d", cpt.in_a, cpt.in_b), func(t *testing.T) {
			c := choosePoint(cpt.in_a, cpt.in_b)
			switch cpt.in_b {
			case NORTH:
				if !(c.Y > cpt.in_a.Y) {
					t.Errorf("got %d, expected y > %d", c, cpt.in_a)
				}
			case SOUTH:
				if !(c.Y < cpt.in_a.Y) {
					t.Errorf("got %d, expected y < %d", c, cpt.in_a)
				}
			case EAST:
				if !(c.X < cpt.in_a.X) {
					t.Errorf("got %d, expected x < %d", c, cpt.in_a)
				}
			case WEST:
				if !(c.X > cpt.in_a.X) {
					t.Errorf("got %d, expected x > %d", c, cpt.in_a)
				}
			}
		})
	}
}
*/
