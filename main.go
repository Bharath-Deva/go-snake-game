package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
)

var gameSpeed = time.Second / 6 // 10 FPS
var fruitSpeed = time.Second * 10

type Point struct {
	X, Y int
}

type Fruit struct {
	point      Point
	lastUpdate time.Time
}

type Game struct {
	snake      []Point
	direction  Point
	lastUpdate time.Time
	fruit      Fruit
}

func AbsInt(n int, isHeight bool) int {
	if n >= 0 {
		return n
	}
	if isHeight {
		fmt.Println("Negative value detected, adjusting for height:", -n+(screenHeight/gridSize))
		return -n + (screenHeight / gridSize) - 1
	}
	return -n + (screenWidth / gridSize) - 1
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyUp) && g.direction.Y == 0 {
		g.direction = Point{X: 0, Y: -1} // Move up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.direction.Y == 0 {
		g.direction = Point{X: 0, Y: 1} // Move down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.direction.X == 0 {
		g.direction = Point{X: -1, Y: 0} // Move left
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && g.direction.X == 0 {
		g.direction = Point{X: 1, Y: 0} // Move right
	}

	g.updateFruit()

	if time.Since(g.lastUpdate) < gameSpeed {
		return nil // Limit the update rate to 10 FPS
	}

	g.updateSnake(&g.snake, g.direction)

	g.lastUpdate = time.Now()
	return nil
}

func (g *Game) updateFruit() {
	if time.Since(g.fruit.lastUpdate) > fruitSpeed {
		g.fruit.point = Point{
			X: rand.Intn(screenWidth / gridSize),
			Y: rand.Intn(screenHeight / gridSize),
		}
		g.fruit.lastUpdate = time.Now()
	}
}

func (g *Game) updateSnake(
	snake *[]Point,
	direction Point,
) {
	head := (*snake)[0]

	// Check for fruit collision

	newHead := Point{
		X: AbsInt((head.X+direction.X)%(screenWidth/gridSize), false),
		Y: AbsInt((head.Y+direction.Y)%(screenHeight/gridSize), true),
	}

	if g.fruit.point.X == head.X && g.fruit.point.Y == head.Y {
		// Grow the snake by adding a new segment at the head
		*snake = append([]Point{newHead}, *snake...)
		g.fruit.point = Point{-1, -1} // Reset fruit position
		return
	}
	*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)

}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.fruit.point.X >= 0 && g.fruit.point.Y >= 0 {
		vector.DrawFilledRect(screen, float32(g.fruit.point.X*gridSize), float32(g.fruit.point.Y*gridSize), float32(gridSize), float32(gridSize), color.RGBA{255, 0, 0, 255}, true)
	}

	for _, point := range g.snake {
		vector.DrawFilledRect(screen, float32(point.X*gridSize), float32(point.Y*gridSize), float32(gridSize), float32(gridSize), color.White, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	game := &Game{
		snake: []Point{
			{X: screenWidth / gridSize / 2, Y: screenHeight / gridSize / 2},
		},
		direction: Point{X: 1, Y: 0}, // Start moving to the right
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
