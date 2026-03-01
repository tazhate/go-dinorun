package scenes

import (
	"fmt"
	"strings"

	"github.com/ahmad-alkadri/go-dinorun/internal/app/game"
	"github.com/ahmad-alkadri/go-dinorun/internal/app/sprites"
)

func RenderGame(
	MaxX, MaxY, spriteDinoY, groundSpeed *int,
	dino *sprites.SpriteDino,
	ground *sprites.SpriteGround,
	cactuses *sprites.SpriteCactuses,
	pteranodons *sprites.SpritePteranodons,
	scores *game.GameScores,
	exitChan chan bool,
	progressBar string,
) {
	scene := make([]string, *MaxY)
	for i := range scene {
		scene[i] = strings.Repeat(" ", *MaxX)
	}

	printSprite(5, *spriteDinoY, dino.Render(), scene)
	printSprite(0, *MaxY-1, ground.Render(*groundSpeed), scene)

	for _, cactus := range cactuses.Group {
		printSprite(cactus.Xoffset, *MaxY-1, cactus.Render(), scene)
	}
	for _, ptera := range pteranodons.Group {
		printSprite(ptera.Xoffset, *MaxY-1, ptera.Render(), scene)
	}

	var output strings.Builder
	output.WriteString("\u001B[?1049h")
	output.WriteString("\u001B[?25l")
	output.WriteString("\u001B[2J")
	output.WriteString("\u001B[H")
	if progressBar != "" {
		output.WriteString(progressBar + "\n")
	}
	output.WriteString(fmt.Sprintf("Score: %d\n", scores.Print()))
	for _, line := range scene {
		output.WriteString(line + "\n")
	}
	fmt.Print(output.String())
}

func AreClashing(
	MaxY, spriteDinoY *int,
	dino *sprites.SpriteDino,
	cactuses *sprites.SpriteCactuses,
	pteranodons *sprites.SpritePteranodons,
) bool {
	dinoCells := extractSpriteCells(5, *spriteDinoY, dino.Render())
	var cactusCells [][2]int
	var pteraCells [][2]int
	for _, cactus := range cactuses.Group {
		cactusCells = extractSpriteCells(cactus.Xoffset, *MaxY-1, cactus.Render())
		if shareChild(dinoCells, cactusCells) {
			return true
		}
	}
	for _, ptera := range pteranodons.Group {
		pteraCells = extractSpriteCells(ptera.Xoffset, *MaxY-1, ptera.Render())
		if shareChild(dinoCells, pteraCells) {
			return true
		}
	}
	return false
}

// RenderFinalFrame renders the last frame with optional progress bar.
func RenderFinalFrame(scene []string, score int, progressBar string) {
	var output strings.Builder
	output.WriteString("\u001B[2J")
	output.WriteString("\u001B[H")
	output.WriteString("\u001B[?25h")
	if progressBar != "" {
		output.WriteString(progressBar + "\n")
	}
	output.WriteString(fmt.Sprintf("Score: %d\n", score))
	for _, line := range scene {
		output.WriteString(line + "\n")
	}
	fmt.Print(output.String())
}

// RenderFinalScene creates the final scene with all game elements in their last positions
func RenderFinalScene(maxX, maxY, spriteDinoY, groundSpeed int,
	dino *sprites.SpriteDino,
	ground *sprites.SpriteGround,
	cactuses *sprites.SpriteCactuses,
	pteranodons *sprites.SpritePteranodons) []string {

	finalScene := make([]string, maxY)
	for i := range finalScene {
		finalScene[i] = strings.Repeat(" ", maxX)
	}

	printSprite(5, spriteDinoY, dino.Render(), finalScene)
	printSprite(0, maxY-1, ground.Render(groundSpeed), finalScene)

	for _, cactus := range cactuses.Group {
		printSprite(cactus.Xoffset, maxY-1, cactus.Render(), finalScene)
	}
	for _, ptera := range pteranodons.Group {
		printSprite(ptera.Xoffset, maxY-1, ptera.Render(), finalScene)
	}

	return finalScene
}
