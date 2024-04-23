package main

import (
	_ "errors"
	"fmt"
	"log"

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

	fieldWidth  = 416
	fieldHeight = 336

	// Assume square:
	frameSize = 100.0 // Original size of png region (square)
	// frameWidth  = 32.0
	// frameHeight = 32.0
)

type Field struct {
	cells [13 * 4][13 * 4]CellType
}

type Game struct {
	player *tank
	count  int
}

var (
	blockSize   float64
	cellSize    float64
	playerScale float64
	tankUp      *ebiten.Image
	tankDown    *ebiten.Image
	tankLeft    *ebiten.Image
	tankRight   *ebiten.Image
)

func init() {
	blockSize = fieldHeight / 13.0
    if (fieldHeight > fieldWidth) {
        blockSize = fieldWidth / 13.0
    }
	cellSize = blockSize / 4.0 // 1 Block == 16 cells (4x4)
	playerScale = blockSize / frameSize

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
}

type tank struct {
	x      float64
	y      float64
	hitbox [2][2]int
	dir    int8
	vx     float64
	vy     float64
	vx0    float64
	vy0    float64
	debug  string
}

func (t *tank) draw(screen *ebiten.Image) {
	s := tankUp
	switch {
	case t.dir == 1:
		s = tankRight
	case t.dir == 2:
		s = tankDown
	case t.dir == 3:
		s = tankLeft
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(playerScale, playerScale)
	// op.GeoM.Translate(t.x/cellSize, t.y/cellSize)
	// op.GeoM.Translate(float64(t.x)/blockSize, float64(t.y)/blockSize)
	op.GeoM.Translate(t.x, t.y)
	screen.DrawImage(s, op)
}

func (t *tank) update() {
	t.vx0 = t.vx
	t.vy0 = t.vy
	t.x += t.vx
	t.y += t.vy
	t.vx = 0
	t.vy = 0
}

func (g *Game) Update() error {

	if g.player == nil {
        // Este de abajo lo dejé afuera porque me agrandaba las dimensiones y todavía no entiendo qué está pasando:
		// g.player = &tank{x: fieldWidth*cellSize - blockSize, y: fieldHeight - blockSize}
		g.player = &tank{x: (fieldWidth - blockSize)/2, y: fieldHeight - blockSize}
	}
    if (ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)) {
        g.player.debug = "Mouse"
    }
    
	g.count++
	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyK)) && g.player.vy0 <= 0 && g.player.vx0 == 0 {
		g.player.vy = -cellSize/4
		g.player.dir = 0
		g.player.update()
		g.player.debug = "^"
		return nil
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyJ)) && g.player.vy0 >= 0 && g.player.vx0 == 0 {
		g.player.vy = cellSize/4
		g.player.dir = 2
		g.player.update()
		g.player.debug = "v"
		return nil
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyL)) && g.player.vx0 >= 0 && g.player.vy0 == 0 {
		g.player.vx = cellSize/4
		g.player.dir = 1
		g.player.update()
		g.player.debug = "->"
		return nil
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyH)) && g.player.vx0 <= 0 && g.player.vy0 == 0 {
		g.player.vx = -cellSize/4
		g.player.dir = 3
		g.player.update()
		g.player.debug = "<-"
		return nil
	}
	g.player.vx0 = 0
	g.player.vy0 = 0

    if (ebiten.IsKeyPressed(ebiten.KeyEscape)) {
        return ebiten.Termination
    }


	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.draw(screen)
	// Show the message.
    // Display the information with "X: xx, Y: xx" format
	msg := fmt.Sprintf("TPS: %0.2f\n%s\n(x,y)=(%.2f, %.2f)\n(vx0, vy0)=(%.2f, %.2f)\nblockSize: %.2f\ncellSize: %.2f", ebiten.ActualTPS(), g.player.debug, g.player.x, g.player.y, g.player.vx0, g.player.vy0, blockSize, cellSize)
    x, y := ebiten.CursorPosition()
    msg += fmt.Sprintf("\nX: %d, Y: %d", x, y)
    msg += fmt.Sprintf("\nScale: %.2f", playerScale)
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
