package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Point struct {
	X int
	Y int
}

type Game struct {
	snake      []Point
	food       Point
	direction  Point
	gameOver   bool
	score      int
	width      int
	height     int
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().UnixNano())

	width, height := termbox.Size()
	game := &Game{
		width:  width,
		height: height,
		direction: Point{
			X: 1,
			Y: 0,
		},
	}
	game.initSnake()
	game.placeFood()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

gameLoop:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Key == termbox.KeyArrowUp && game.direction.Y == 0:
					game.direction = Point{X: 0, Y: -1}
				case ev.Key == termbox.KeyArrowDown && game.direction.Y == 0:
					game.direction = Point{X: 0, Y: 1}
				case ev.Key == termbox.KeyArrowLeft && game.direction.X == 0:
					game.direction = Point{X: -1, Y: 0}
				case ev.Key == termbox.KeyArrowRight && game.direction.X == 0:
					game.direction = Point{X: 1, Y: 0}
				case ev.Ch == 'q' || ev.Key == termbox.KeyEsc:
					break gameLoop
				}
			}
		case <-ticker.C:
			if !game.gameOver {
				game.move()
				game.checkCollision()
				game.draw()
			}
		}
	}
}

func (g *Game) initSnake() {
	centerX := g.width / 2
	centerY := g.height / 2
	g.snake = []Point{
		{X: centerX, Y: centerY},
		{X: centerX - 1, Y: centerY},
		{X: centerX - 2, Y: centerY},
	}
}

func (g *Game) placeFood() {
	g.food = Point{
		X: rand.Intn(g.width-2) + 1,
		Y: rand.Intn(g.height-2) + 1,
	}
	// Make sure food doesn't spawn on snake
	for _, p := range g.snake {
		if p.X == g.food.X && p.Y == g.food.Y {
			g.placeFood()
			return
		}
	}
}

func (g *Game) move() {
	head := g.snake[0]
	newHead := Point{
		X: head.X + g.direction.X,
		Y: head.Y + g.direction.Y,
	}

	// Wrap around screen edges
	if newHead.X <= 0 {
		newHead.X = g.width - 2
	} else if newHead.X >= g.width-1 {
		newHead.X = 1
	}
	if newHead.Y <= 0 {
		newHead.Y = g.height - 2
	} else if newHead.Y >= g.height-1 {
		newHead.Y = 1
	}

	g.snake = append([]Point{newHead}, g.snake...)

	// Check if snake ate food
	if newHead.X == g.food.X && newHead.Y == g.food.Y {
		g.score++
		g.placeFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func (g *Game) checkCollision() {
	head := g.snake[0]
	
	// Check collision with itself
	for _, p := range g.snake[1:] {
		if head.X == p.X && head.Y == p.Y {
			g.gameOver = true
			return
		}
	}
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw border
	for x := 0; x < g.width; x++ {
		termbox.SetCell(x, 0, '═', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, g.height-1, '═', termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < g.height; y++ {
		termbox.SetCell(0, y, '║', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(g.width-1, y, '║', termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.SetCell(0, 0, '╔', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, 0, '╗', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, g.height-1, '╚', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, g.height-1, '╝', termbox.ColorWhite, termbox.ColorDefault)

	// Draw snake
	for i, p := range g.snake {
		color := termbox.ColorGreen
		if i == 0 {
			color = termbox.ColorYellow // Head is different color
		}
		termbox.SetCell(p.X, p.Y, '■', color, termbox.ColorDefault)
	}

	// Draw food
	termbox.SetCell(g.food.X, g.food.Y, '●', termbox.ColorRed, termbox.ColorDefault)

	// Draw score
	scoreStr := fmt.Sprintf(" Score: %d ", g.score)
	for i, ch := range scoreStr {
		termbox.SetCell(i+2, g.height-1, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	// Game over message
	if g.gameOver {
		msg := " GAME OVER - Press 'q' to quit "
		startX := (g.width - len(msg)) / 2
		startY := g.height / 2
		for i, ch := range msg {
			termbox.SetCell(startX+i, startY, ch, termbox.ColorRed, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}