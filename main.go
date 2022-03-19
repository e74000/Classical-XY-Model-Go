package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

type Game struct {
	pixels []byte

	spins []float64

	tmp float64
	ext float64
	ins float64

	x int
	y int
	s int
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.x, g.y
}

func mod(d, m int) int {
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

func fmod(d, m float64) float64 {
	res := math.Mod(d, m)
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

func torus(mx, my, x, y int) (int, int) {
	return mod(x, mx), mod(y, my)
}

func (g *Game) Update(_ *ebiten.Image) error {
	spinsNext := make([]float64, g.x*g.y)
	copy(spinsNext, g.spins)

	for i := 0; i < g.y*g.x; i++ {
		px := rand.Int() % g.x
		py := rand.Int() % g.y

		f := 0.0

		for _, j := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			tx, ty := torus(g.x, g.y, px+j[0], py+j[1])
			f += math.Sin(g.spins[py*g.x+px] - g.spins[ty*g.x+tx])
		}

		f *= g.ins
		f += g.tmp * (rand.Float64()*2 - 1)
		f += g.ext * math.Sin(g.spins[py*g.x+px])

		spinsNext[py*g.x+px] = fmod(g.spins[py*g.x+px]-f, 2*math.Pi)
	}

	copy(g.spins, spinsNext)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.tmp += 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.tmp > 0 {
		g.tmp -= 0.01
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.ins += 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.ins > 0 {
		g.ins -= 0.01
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, 4*g.x*g.y)

		for i := 0; i < g.x*g.y; i++ {
			g.pixels[4*i+3] = 0xff
		}
	}

	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			i := y*g.x + x

			c := colorful.Hsv(g.spins[i]/(math.Pi)*180, 1, 1)

			g.pixels[4*i+0] = uint8(c.R * 255)
			g.pixels[4*i+1] = uint8(c.G * 255)
			g.pixels[4*i+2] = uint8(c.B * 255)
		}
	}

	err := screen.ReplacePixels(g.pixels)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		tmp float64
		ext float64
		ins float64
		xs  int
		ys  int
		ss  int
		fs  bool
	)

	flag.Float64Var(&tmp, "t", 0, "The temperature of the simulation")
	flag.Float64Var(&ext, "e", 0, "The external field value")
	flag.Float64Var(&ins, "i", 0.4, "The interaction strength")
	flag.IntVar(&xs, "x", 1920/8, "X resolution")
	flag.IntVar(&ys, "y", 1080/8, "Y resolution")
	flag.IntVar(&ss, "s", 4, "Scale (only works if not fullscreen)")
	flag.BoolVar(&fs, "f", false, "Fullscreen")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	g := &Game{
		x: xs,
		y: ys,
		s: ss,

		tmp: tmp,
		ext: ext,
		ins: ins,
	}

	g.spins = make([]float64, g.x*g.y)

	for i := 0; i < g.x*g.y; i++ {
		g.spins[i] = rand.Float64() * 2 * math.Pi
	}

	ebiten.SetWindowTitle("Classical XY Model")
	ebiten.SetWindowSize(g.x*g.s, g.y*g.s)
	ebiten.SetFullscreen(fs)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
