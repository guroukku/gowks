package	gowks

import(
	"fmt"
	"testing"
)

var intPowTests = []struct {
	in_a int
	in_b int
	out int
}{
	{2, 1, 2},
	{2, 6, 64},
}
func TestIntPowValid(t *testing.T) {
	for _, ipt := range intPowTests {
		t.Run(fmt.Sprintf("%d, %d", ipt.in_a, ipt.in_b), func(t *testing.T) {
			i := intPow(ipt.in_a, ipt.in_b)
			if i != ipt.out {
				t.Errorf("got %d, expected %d", i, ipt.out)
			}
		})
	}
}

var distTests = []struct {
	in_a coord
	in_b coord
	out float64
}{
	{coord{1,1}, coord{1,1}, 0.0},
	{coord{3,0}, coord{0,4}, 5.0},
}
func TestDistance(t *testing.T) {
	for _, dt := range distTests {
		t.Run(fmt.Sprintf("%d, %d", dt.in_a, dt.in_b), func(t *testing.T) {
			d := distance(dt.in_a, dt.in_b)
			if d != dt.out {
				t.Errorf("got %f, expected %f", d, dt.out)
			}
		})
	}
}

var inLineTests = []struct {
	in_a coord
	in_b coord
	in_c coord
	out bool
}{
	{coord{0,0}, coord{5,0}, coord{3,0}, true},
	{coord{0,0}, coord{3,0}, coord{5,0}, false},
	{coord{-5,-5}, coord{5,5}, coord{0,0}, true},
	{coord{-3,3}, coord{3,1}, coord{0,2}, true},
}
func TestInLine(t *testing.T) {
	for _, ilt := range inLineTests {
		t.Run(fmt.Sprintf("%d, %d, %d", ilt.in_a, ilt.in_b, ilt.in_c), func(t *testing.T) {
			il := inLine(ilt.in_a, ilt.in_b, ilt.in_c)
			if il != ilt.out {
				t.Errorf("got %t, expected %t", il, ilt.out)
			}
		})
	}
}

var inBoundsTests = []struct {
	in_a coord
	in_b coord
	in_c coord
	out bool
}{
	{coord{0,0}, coord{5,5}, coord{1,2}, true},
	{coord{0,0}, coord{5,5}, coord{-1,3}, false},
	{coord{-5,-5}, coord{5,5}, coord{0,0}, true},
}
func TestInBounds(t * testing.T) {
	for _, ibt := range inBoundsTests {
		t.Run(fmt.Sprintf("%d, %d, %d", ibt.in_a, ibt.in_b, ibt.in_c), func(t *testing.T) {
			ib := inBounds(ibt.in_a, ibt.in_b, ibt.in_c)
			if ib != ibt.out {
				t.Errorf("got %t, expected %t", ib, ibt.out)
			}
		})
	}
}
