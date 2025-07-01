package main

import (
	"fmt"
	"gifhelper"
)

func main() {
	numRows := 250
	numCols := 250

	initialBoard := InitializeBoard(numRows, numCols)

	frac := 0.05 // tells us what percentage of interior cells to color with predators

	// how many predator rows and columns are there?
	predRows := frac * float64(numRows)
	predCols := frac * float64(numCols)

	midRow := numRows / 2
	midCol := numCols / 2

	// a little for loop to fill predators
	for r := midRow - int(predRows/2); r < midRow+int(predRows/2); r++ {
		for c := midCol - int(predCols/2); c < midCol+int(predCols/2); c++ {
			initialBoard[r][c][1] = 1.0
		}
	}

	// make prey concentration 1 at every cell
	for i := range initialBoard {
		for j := range initialBoard[i] {
			initialBoard[i][j][0] = 1.0
		}
	}

	// let's set some parameters too
	numGens := 20000 // number of iterations
	feedRate := 0.034
	killRate := 0.095

	preyDiffusionRate := 0.2 // prey are twice as fast at running
	predatorDiffusionRate := 0.1

	// let's declare kernel
	var kernel [3][3]float64
	kernel[0][0] = .05
	kernel[0][1] = .2
	kernel[0][2] = .05
	kernel[1][0] = .2
	kernel[1][1] = -1.0
	kernel[1][2] = .2
	kernel[2][0] = .05
	kernel[2][1] = .2
	kernel[2][2] = .05

	// let's simulate Gray-Scott!
	// result will be a collection of Boards corresponding to each generation.
	boards := SimulateGrayScott(initialBoard, numGens, feedRate, killRate, preyDiffusionRate, predatorDiffusionRate, kernel)

	fmt.Println("Done with simulation!")

	// we will draw what we have generated.
	fmt.Println("Drawing boards to file.")

	//for the visualization, we are only going to draw every nth board to be more efficient
	n := 100

	cellWidth := 1 // each cell is 1 pixel

	imageList := DrawBoards(boards, cellWidth, n)
	fmt.Println("Boards drawn! Now draw GIF.")

	outFile := "Gray-Scott"
	gifhelper.ImagesToGIF(imageList, outFile) // code is given
	fmt.Println("GIF drawn!")
}

//DrawBoard takes a Board objects as input along with a cellWidth and n parameter.
//It returns an image corresponding to drawing every nth board to a file,
//where each cell is cellWidth x cellWidth pixels.
func DrawBoard(b Board, cellWidth int) image.Image {
    // need to know how many pixels wide and tall to make our image

    height := len(b) * cellWidth
    width := len(b[0]) * cellWidth

    // think of a canvas as a PowerPoint slide that we draw on
    c := canvas.CreateNewCanvas(width, height)

    // canvas will start as black, so we should fill in colored squares
    for i := range b {
        for j := range b[i] {
            prey := b[i][j][0]
            predator := b[i][j][1]

            // we will color each cell according to a color map.
            val := predator / (predator + prey)
            colorMap := moreland.SmoothBlueRed() // blue-red polar spectrum

            //set min and max value of color map
            colorMap.SetMin(0)
            colorMap.SetMax(1)

            //find the color associated with the value predator / (predator + prey)
            color, err := colorMap.At(val)

            if err != nil {
                panic("Error converting color!")
            }

            // draw a rectangle in right place with this color
            c.SetFillColor(color)

            x := i * cellWidth
            y := j * cellWidth
            c.ClearRect(x, y, x+cellWidth, y+cellWidth)
            c.Fill()
        }
    }

    // canvas has an image field that we should return
    return c.GetImage()
}
