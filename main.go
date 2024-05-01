package main

import (
	"bufio"
	_ "errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	// "runtime/internal/math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type CellType int

const (
	CellTypeNone CellType = iota
	CellBrick             // 1
	CellSteel             // 2
	CellLeaves            // 3
	CellPower             // 4
	CellTypeMax  = CellPower
)

const (
	screenWidth  = 548
	screenHeight = 336

	// Assume square:
	frameSize = 100.0 // Original size of png region (square)
	// frameWidth  = 32.0
	// frameHeight = 32.0
)

type Game struct {
	player *tank
	count  int
}

var (
	fieldWidth  = float64(screenWidth - 100)
	fieldHeight = float64(screenHeight)
	blockSize   float64
	cellSize    float64
	tileScale   float64
	tankUp      *ebiten.Image
	tankDown    *ebiten.Image
	tankLeft    *ebiten.Image
	tankRight   *ebiten.Image
	fieldImage  *ebiten.Image
	brickImage  *ebiten.Image
	stageNX     = 13
	Stage1      = make([][]uint8, stageNX)
)

func init() {
	// field must be square:
	if fieldWidth > fieldHeight {
		fieldWidth = fieldHeight
	}
	if fieldWidth < fieldHeight {
		fieldWidth = fieldHeight
	}
	blockSize = fieldWidth / float64(stageNX)
	cellSize = blockSize / 4.0 // 1 Block == 16 cells (4x4)
	tileScale = blockSize / frameSize

	// https://www.deviantart.com/leetzero/art/Tanks-339084110
	// CC Atribution Non Comercial Share Alike 3.0 CC-By-NC-SA
	var err error
	tankUp, _, err = ebitenutil.NewImageFromFile("tank_up.png")
	if err != nil {
		log.Fatal(err)
	}
	tankDown, _, err = ebitenutil.NewImageFromFile("tank_down.png")
	if err != nil {
		log.Fatal(err)
	}
	tankLeft, _, err = ebitenutil.NewImageFromFile("tank_left.png")
	if err != nil {
		log.Fatal(err)
	}
	tankRight, _, err = ebitenutil.NewImageFromFile("tank_right.png")
	if err != nil {
		log.Fatal(err)
	}

	fieldImage, _, err = ebitenutil.NewImageFromFile("field.png")
	if err != nil {
		log.Fatal(err)
	}
	brickImage, _, err = ebitenutil.NewImageFromFile("brickblock.png")
	if err != nil {
		log.Fatal(err)
	}
}

type tank struct {
	x      float64
	y      float64
	hitbox [2][2]float64
	dir    int8
	vx     float64
	vy     float64
	vx0    float64
	vy0    float64
	debug  string
    speed  float64
}

func (t *tank) draw(screen *ebiten.Image) {
	s := tankUp
    xdiff, ydiff := 0.0, 0.0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(tileScale, tileScale)
	switch {
	case t.dir == 1:
		op.GeoM.Rotate(math.Pi / 2)
        xdiff = blockSize
		// s = tankRight
	case t.dir == 2:
		op.GeoM.Rotate(math.Pi)
        ydiff = blockSize
        xdiff = blockSize
		// s = tankDown
	case t.dir == 3:
		op.GeoM.Rotate(-math.Pi / 2)
        ydiff = blockSize
		// s = tankLeft
	}
	// op.GeoM.Translate(t.x/cellSize, t.y/cellSize)
	// op.GeoM.Translate(float64(t.x)/blockSize, float64(t.y)/blockSize)
	op.GeoM.Translate(t.x + xdiff, t.y + ydiff)
	screen.DrawImage(s, op)
}

func (t *tank) update() {
	t.vx0 = t.vx
	t.vy0 = t.vy
	t.x += t.vx
	t.y += t.vy
	if t.x > fieldHeight-blockSize {
		t.x = fieldHeight - blockSize
	}
	if t.x < 0 {
		t.x = 0
	}
	if t.y > fieldHeight-blockSize {
		t.y = fieldHeight - blockSize
	}
	if t.y < 0 {
		t.y = 0
	}
	t.vx = 0
	t.vy = 0
}

func init() {
	for i := range Stage1 {
		Stage1[i] = make([]uint8, stageNX)
	}
	file, err := os.Open("Stage1.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	j := 0
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		s := scanner.Text()
		v := strings.Fields(s)
		// fmt.Println(v)
		for i := range v {
			vj, err := strconv.Atoi(v[i])
			if err != nil {
				fmt.Println(err)
			}
			Stage1[i][j] = uint8(vj)
		}
		j++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	/* for i := range Stage1 {
		fmt.Println(Stage1[i])
	} */
}

func (g *Game) drawTiles(screen *ebiten.Image) {
	/* var (
	       nx = fieldWidth / blockSize
	       ny = fieldHeight / blockSize
	   )
	   fmt.Println(nx, ny) */
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < stageNX; i++ {
		for j := 0; j < stageNX; j++ {
			if Stage1[i][j] == 1 {
				op.GeoM.Scale(tileScale, tileScale)
				op.GeoM.Translate(blockSize*float64(i), blockSize*float64(j))
				screen.DrawImage(brickImage, op)
				op.GeoM.Reset()
			}
		}
	}
	// screen.DrawImage(brickImage, op)
}

func (g *Game) Update() error {

	if g.player == nil {
		// Este de abajo lo dejé afuera porque me agrandaba las dimensiones y todavía no entiendo qué está pasando:
		// g.player = &tank{x: fieldWidth*cellSize - blockSize, y: fieldHeight - blockSize}
		g.player = &tank{x: blockSize * 4, y: fieldHeight - blockSize, speed: cellSize / 8.0}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.player.debug = "Mouse"
	}

	g.count++
	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyK)) && g.player.vy0 <= 0 && g.player.vx0 == 0 {
		g.player.vy = -g.player.speed
		g.player.dir = 0
		g.player.update()
		g.player.debug = "^\n|"
		return nil
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyJ)) && g.player.vy0 >= 0 && g.player.vx0 == 0 {
		g.player.vy = g.player.speed
		g.player.dir = 2
		g.player.update()
		g.player.debug = "|\nv"
		return nil
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyL)) && g.player.vx0 >= 0 && g.player.vy0 == 0 {
		g.player.vx = g.player.speed
		g.player.dir = 1
		g.player.update()
		g.player.debug = "->\n"
		return nil
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyH)) && g.player.vx0 <= 0 && g.player.vy0 == 0 {
		g.player.vx = -g.player.speed
		g.player.dir = 3
		g.player.update()
		g.player.debug = "<-\n"
		return nil
	}
	g.player.vx0 = 0
	g.player.vy0 = 0

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(fieldWidth/336, fieldHeight/336)
	screen.DrawImage(fieldImage, op)
	g.drawTiles(screen)
	// op = &ebiten.DrawImageOptions{}
	g.player.draw(screen)
	x, y := ebiten.CursorPosition()
	msg := fmt.Sprintf("TPS: %0.2f\n%s", ebiten.ActualFPS(), g.player.debug)
	msg += fmt.Sprintf("\n(x,y)=(%.2f, %.2f)", g.player.x, g.player.y)
	msg += fmt.Sprintf("\n(vx0, vy0)=(%.2f, %.2f)", g.player.vx0, g.player.vy0)
	msg += fmt.Sprintf("\nblockSize: %.2f\ncellSize: %.2f", blockSize, cellSize)
	msg += fmt.Sprintf("\nfieldWidth: %.2f\nfieldHeight: %.2f", fieldWidth, fieldHeight)
	msg += fmt.Sprintf("\nX: %d, Y: %d", x, y)
	msg += fmt.Sprintf("\nScale: %.2f", tileScale)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// No es la cantidad de píxeles sino es algo con la escala...?
	// En este ejemplo hay como 4300 de ancho y 2400 de alto ponele. En qué unidades está???
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Tanks")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal()
	}
}
