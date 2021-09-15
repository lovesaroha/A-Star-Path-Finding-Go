package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

// Spot in a grid.
type spot struct {
	Gvalue   float64
	Hvalue   float64
	Fvalue   float64
	Ivalue   int
	Jvalue   int
	Current  bool
	Previous *spot
	Path     bool
	Block    bool
}

var rows = 50
var columns = 50
var grid [50][50]spot

func main() {
	var openSet = make([]*spot, 0, 25)
	var closedSet = make([]*spot, 0, 25)
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			var random = rand.Intn(100)
			var block bool
			if random > 70 && (i != 25 || j != 2) && (i != 39 || j != 40) {
				block = true
			}
			var newSpot = spot{Ivalue: i, Jvalue: j, Block: block}
			grid[i][j] = newSpot
		}
	}
	var start = &grid[25][2]
	var end = &grid[39][40]
	openSet = append(openSet, start)
	for len(openSet) > 0 {
		var winIndex int
		for k, v := range openSet {
			if v.Fvalue < openSet[winIndex].Fvalue {
				winIndex = k
			}
		}
		if openSet[winIndex] == end {
			fmt.Println("Done")
			showPath(start, end)
			return
		}
		var current = openSet[winIndex]
		current.Current = true
		showGrid(start, end, false, true)
		// Removing from openSet.
		openSet = append(openSet[0:winIndex], openSet[(winIndex+1):]...)
		// Adding to closedSet.
		closedSet = append(closedSet, current)

		for _, neighbor := range findNeighbors(current) {
			if contains(closedSet, neighbor) || neighbor.Block {
				continue
			}
			// Neighbor not in closedSet then.
			var Gvalue = current.Gvalue + 1
			// Not in openSet.
			if !contains(openSet, neighbor) {
				neighbor.Gvalue = Gvalue
				openSet = append(openSet, neighbor)
			} else if Gvalue < neighbor.Gvalue {
				neighbor.Gvalue = Gvalue
			} else {
				continue
			}
			// f(n) = g(n) + h(n).
			neighbor.Hvalue = distance(neighbor.Ivalue, neighbor.Jvalue, 39, 40)
			neighbor.Fvalue = neighbor.Gvalue + neighbor.Hvalue
			neighbor.Previous = current
		}
	}
	fmt.Println("No Path")
}

// Show grid data.
func showGrid(start, end *spot, showPathSpot, cls bool) {
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if grid[i][j].Current && grid[i][j] != *start && !showPathSpot {
				g := color.New(color.FgGreen, color.Bold)
				g.Printf(" x ")
			} else if grid[i][j] == *start {
				fmt.Printf(" S ")
			} else if grid[i][j] == *end {
				fmt.Printf(" E ")
			} else if grid[i][j].Block {
				r := color.New(color.FgRed, color.Bold)
				r.Printf(" o ")
			} else if grid[i][j].Path && showPathSpot {
				d := color.New(color.FgCyan, color.Bold)
				d.Printf(" x ")
			} else {
				fmt.Printf("   ")
			}
		}
		fmt.Println(" ")
	}

	if cls {
		osName := runtime.GOOS
		switch osName {
		case "windows":
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
			return
		case "linux":
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	}
}

// Find neighbors.
func findNeighbors(current *spot) []*spot {
	var neighbors []*spot
	if current.Ivalue > 0 {
		neighbors = append(neighbors, &grid[(current.Ivalue - 1)][current.Jvalue])
	}
	if current.Jvalue > 0 {
		neighbors = append(neighbors, &grid[current.Ivalue][(current.Jvalue-1)])
	}
	if current.Ivalue < (rows - 1) {
		neighbors = append(neighbors, &grid[(current.Ivalue + 1)][current.Jvalue])
	}
	if current.Jvalue < (columns - 1) {
		neighbors = append(neighbors, &grid[current.Ivalue][(current.Jvalue+1)])
	}
	if current.Ivalue > 0 && current.Jvalue > 0 {
		neighbors = append(neighbors, &grid[(current.Ivalue - 1)][(current.Jvalue-1)])
	}
	if current.Ivalue < (rows-1) && current.Jvalue < (columns-1) {
		neighbors = append(neighbors, &grid[(current.Ivalue + 1)][(current.Jvalue+1)])
	}
	if current.Jvalue > 0 && current.Ivalue < (rows-1) {
		neighbors = append(neighbors, &grid[(current.Ivalue + 1)][(current.Jvalue-1)])
	}
	if current.Ivalue > 0 && current.Jvalue < (columns-1) {
		neighbors = append(neighbors, &grid[(current.Ivalue - 1)][(current.Jvalue+1)])
	}
	return neighbors
}

// Check if set includes element.
func contains(set []*spot, element *spot) bool {
	for _, cs := range set {
		if cs == element {
			return true
		}
	}
	return false
}

// Distance between two points.
func distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt((float64(x2)-float64(x1))*(float64(x2)-float64(x1)) + (float64(y2)-float64(y1))*(float64(y2)-float64(y1)))
}

// Show path.
func showPath(start, end *spot) {
	var spot = end
	for spot.Previous != nil {
		spot.Path = true
		showGrid(start, end, true, true)
		spot = spot.Previous
	}
	showGrid(start, end, true, false)
}
