// author: Jacky Boen

package main

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

type blockMap [4][4]byte
type block struct {
	m     blockMap
	color sdl.Color
}

var (
	winWidth    int16 = 800
	winHeight   int16 = 600
	blockSize   int16 = 20
	colorGreen        = sdl.Color{G: 200, A: 255}
	colorRed          = sdl.Color{R: 200, A: 255}
	colorBlue         = sdl.Color{B: 200, A: 255}
	colorOrange       = sdl.Color{R: 180, G: 100, A: 255}
	colorTeal         = sdl.Color{G: 180, B: 200, A: 255}
	colorYellow       = sdl.Color{R: 240, G: 230, A: 255}
	colorPurple       = sdl.Color{R: 240, B: 230, A: 255}
	colorWhite        = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	blockI            = block{
		m: blockMap{
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
		},
		color: colorTeal,
	}
	blockJ = block{
		m: blockMap{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
		},
		color: colorBlue,
	}
	blockL = block{
		m: blockMap{
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
		},
		color: colorOrange,
	}
	blockO = block{
		m: blockMap{
			{0, 0, 0, 0},
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		},
		color: colorYellow,
	}
	blockS = block{
		m: blockMap{
			{0, 1, 1, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		color: colorGreen,
	}
	blockT = block{
		m: blockMap{
			{1, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		color: colorPurple,
	}
	blockZ = block{
		m: blockMap{
			{1, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		color: colorRed,
	}
)

func rotateBlock(b block, rot int) block {
	result := block{
		color: b.color,
	}
	var r1, c1 int
	for r := range b.m {
		for c := range b.m[r] {
			switch rot {
			case 1:
				r1 = c
				c1 = r
			case 2:
				r1 = 3 - r
				c1 = c
			case 3:
				r1 = 3 - c
				c1 = r
			default:
				r1 = r
				c1 = c
			}
			result.m[r][c] = b.m[r1][c1]
		}
	}
	return result
}

func dim(c uint8, pct int) uint8 {
	r := (100*int(c) - int(c)*(pct)) / 100
	if r > 255 {
		r = 255
	}
	if r < 0 {
		r = 0
	}
	return uint8(r)
}

func dimColor(c sdl.Color, pct int) sdl.Color {
	return sdl.Color{R: dim(c.R, pct), G: dim(c.G, pct), B: dim(c.B, pct), A: c.A}
}

func drawCell(r *sdl.Renderer, x, y, w int16, c sdl.Color) {
	m := w / 10

	gfx.FilledPolygonColor(r,
		[]int16{x + m, x + m, x + w - m, x + w - m},
		[]int16{y + m, y + w - m, y + w - m, y + m},
		c)
	gfx.FilledPolygonColor(r,
		[]int16{x, x, x + m, x + m, x + w - m, x + w},
		[]int16{y, y + w, y + w - m, y + m, y + m, y},
		dimColor(c, 10))
	gfx.FilledPolygonColor(r,
		[]int16{x, x + w, x + w, x + w - m, x + w - m, x + m},
		[]int16{y + w, y + w, y, y + m, y + w - m, y + w - m},
		dimColor(c, -10))
}

func drawBlock(r *sdl.Renderer, x, y int16, b block) {
	for row, columns := range b.m {
		for column, v := range columns {
			if v != 0 {
				drawCell(r, x+int16(column)*blockSize, y+int16(row)*blockSize, blockSize, b.color)
			} else {
				x1 := x + int16(column)*blockSize
				y1 := y + int16(row)*blockSize
				gfx.PolygonColor(r,
					[]int16{x1, x1 + blockSize, x1 + blockSize, x1},
					[]int16{y1 + blockSize, y1 + blockSize, y1, y1},
					colorWhite,
				)
			}
		}
	}
}

func fpsleep(start time.Time) {
	delay := 16*time.Millisecond - time.Now().Sub(start)
	if delay < 0 {
		delay = 0
	}
	sdl.Delay(uint32(delay.Milliseconds()))
}

func tetris(r *sdl.Renderer) {
	running := true
	var (
		idx          = 0
		rot          = 0
		x, y   int16 = winWidth / 2, winHeight / 2
		sX, sY int16 = 2, 2
		blocks [4][]block
	)
	blocks[0] = []block{blockI, blockJ, blockL, blockO, blockS, blockT, blockZ}
	for rot := 1; rot <= 3; rot++ {
		blocks[rot] = make([]block, len(blocks[0]))
		for i := range blocks[rot] {
			blocks[rot][i] = rotateBlock(blocks[0][i], rot)
		}
	}
	for running {
		startT := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if ev.Type == sdl.KEYUP {
					if ev.Keysym.Sym == sdl.K_SPACE {
						idx = (idx + 1) % len(blocks[0])
					}
					if ev.Keysym.Sym == sdl.K_LEFT {
						rot--
					}
					if ev.Keysym.Sym == sdl.K_RIGHT {
						rot++
					}
					if rot > 3 {
						rot = 0
					}
					if rot < 0 {
						rot = 3
					}
				}
				log.Printf("rotation %d", rot)
			}
		}

		r.SetDrawColor(0, 0, 0, 255)
		r.Clear()
		drawBlock(r, x, y, blocks[rot][idx])
		r.Present()
		x += sX
		y += sY
		if x >= winWidth-blockSize*4 || x <= 0 {
			sX = -sX
		}
		if y >= winHeight-blockSize*4 || y <= 0 {
			sY = -sY
		}
		fpsleep(startT)
	}
}

func main() {
	window, err := sdl.CreateWindow("tetris", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("creating window: %s", err)
	}

	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("creating renderer: %s", err)
	}
	defer renderer.Destroy()

	tetris(renderer)
}
