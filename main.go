package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/tazhate/go-dinorun/internal/app/game"
	"github.com/tazhate/go-dinorun/internal/app/scenes"
	"github.com/tazhate/go-dinorun/internal/app/sprites"

	"github.com/eiannone/keyboard"
)

// runGame runs one round. Returns (finalScore, shouldRestart).
func runGame(immortal bool, gameSpeed time.Duration, groundSpeed, delayCactus, delayPteranodon, goal int, noEnemies bool) (int, bool) {
	var (
		MaxX, MaxY          int = 70, 18
		baseY               int = MaxY - 2
		spriteDinoY         int = baseY
		delayBetweenEnemies int = 20
		dino                sprites.SpriteDino
		cactuses            sprites.SpriteCactuses
		pteranodons         sprites.SpritePteranodons
		ground              sprites.SpriteGround
		spawnCactusTicker   *time.Ticker
		spawnPteraTicker    *time.Ticker
		frameDist           int = 100
		scores              game.GameScores
		mu                  sync.Mutex
	)

	termW := getTermWidth()
	if termW < MaxX {
		termW = MaxX
	}

	if noEnemies {
		spawnCactusTicker = time.NewTicker(24 * time.Hour)
		spawnPteraTicker = time.NewTicker(24 * time.Hour)
	} else {
		spawnCactusTicker = time.NewTicker(time.Duration(rand.Intn(1000)+delayCactus) * time.Millisecond)
		spawnPteraTicker = time.NewTicker(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)
	}

	defer func() {
		scores.Stop()
		spawnCactusTicker.Stop()
		spawnPteraTicker.Stop()
	}()

	dino.Init(30)
	ground.Init(&MaxX)
	scores.Init()

	// keyboard.GetKeys returns a channel — no goroutine races possible
	keyChan, err := keyboard.GetKeys(10)
	if err != nil {
		fmt.Println("keyboard error:", err)
		return 0, false
	}

	bar := func() string {
		if goal <= 0 {
			return ""
		}
		return renderProgressBar(scores.Print(), goal, termW)
	}

	// gameOver renders final frame, shows prompt, waits for Space or any key.
	gameOver := func() (int, bool) {
		finalScore := scores.Print()
		finalScene := scenes.RenderFinalScene(MaxX, MaxY, spriteDinoY, groundSpeed,
			&dino, &ground, &cactuses, &pteranodons)

		var pb string
		if goal > 0 {
			pb = renderProgressBar(finalScore, goal, termW)
		}
		scenes.RenderFinalFrame(finalScene, finalScore, pb)
		fmt.Print("\n\033[1;33m  ★ GAME OVER ★\033[0m   SPACE = новая игра  │  другая клавиша = выход\n")

		// Drain any buffered keys (e.g. last Space that killed dino)
	drain:
		for {
			select {
			case <-keyChan:
			default:
				break drain
			}
		}

		// Wait for user decision
		event, ok := <-keyChan
		if !ok || event.Err != nil {
			return finalScore, false
		}
		return finalScore, event.Key == keyboard.KeySpace || event.Rune == ' '
	}

	exitChan := make(chan bool, 1) // kept for RenderGame signature compat

	render := func() {
		scenes.RenderGame(&MaxX, &MaxY, &spriteDinoY, &groundSpeed,
			&dino, &ground, &cactuses, &pteranodons,
			&scores, exitChan, bar())
	}

	checkClash := func() bool {
		if immortal {
			return false
		}
		return scenes.AreClashing(&MaxY, &spriteDinoY, &dino, &cactuses, &pteranodons)
	}

	T := 12
	maxHeight := 12.0

	for {
		select {
		case event, ok := <-keyChan:
			if !ok || event.Err != nil {
				return gameOver()
			}
			switch {
			case event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC:
				return gameOver()
			case event.Key == keyboard.KeySpace || event.Rune == ' ':
				if spriteDinoY != baseY {
					continue
				}
				displacements := game.GetDisplacements(T, maxHeight)
				for i := 0; i <= 2*T; i++ {
					spriteDinoY -= displacements[i]
					render()
					if checkClash() {
						return gameOver()
					}
					// Check for exit key during jump
					select {
					case ev, ok := <-keyChan:
						if !ok || ev.Err != nil || ev.Key == keyboard.KeyEsc || ev.Key == keyboard.KeyCtrlC {
							return gameOver()
						}
					default:
					}
					cactuses.Update()
					pteranodons.Update()
					frameDist++
					time.Sleep(gameSpeed * time.Millisecond)
				}
			}

		case <-spawnCactusTicker.C:
			mu.Lock()
			if frameDist > delayBetweenEnemies {
				var c sprites.SpriteCactus
				c.Init(MaxX, groundSpeed)
				cactuses.Add(c)
				frameDist = 0
			}
			mu.Unlock()
			spawnCactusTicker.Reset(time.Duration(rand.Intn(1000)+delayCactus) * time.Millisecond)

		case <-spawnPteraTicker.C:
			mu.Lock()
			if frameDist > delayBetweenEnemies {
				var p sprites.SpritePteranodon
				p.Init(MaxX, groundSpeed, 30)
				pteranodons.Add(p)
				frameDist = 0
			}
			mu.Unlock()
			spawnPteraTicker.Reset(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)

		default:
			render()
			if checkClash() {
				return gameOver()
			}
			cactuses.Update()
			pteranodons.Update()
			frameDist++
			time.Sleep(gameSpeed * time.Millisecond)
		}
	}
}

func main() {
	immortal := flag.Bool("immortal", false, "god mode: no collision detection")
	speed := flag.Int("speed", 5, "game speed 1 (slow) .. 10 (insane)")
	noEnemies := flag.Bool("no-enemies", false, "disable all enemies")
	goal := flag.Int("goal", 1000, "score goal for progress bar (0 = disable bar)")
	flag.Parse()

	if *speed < 1 {
		*speed = 1
	}
	if *speed > 10 {
		*speed = 10
	}

	gameSpeed := time.Duration(44 - (*speed * 4))
	groundSpeed := 1 + (*speed / 4)
	delayCactus := 700 - (*speed * 50)
	delayPteranodon := 1400 - (*speed * 100)

	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return
	}

	var lastScore int
	for {
		score, restart := runGame(*immortal, gameSpeed, groundSpeed, delayCactus, delayPteranodon, *goal, *noEnemies)
		lastScore = score
		if !restart {
			break
		}
	}

	keyboard.Close()
	game.HandleGameOver(lastScore)
}
