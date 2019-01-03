package main

import (
	"log"
	"math"
	"time"
)

/*
I'm pretty sure this is an NP problem. If I was smart I'd be able to prove it.

sprinkler
	cross cover pattern
	0 . 0
	. s .
	0 . 0

Q sprinkler
	. . .
	. Q .
	. . .

Irid Sprinkler
	. . . . .
	. . . . .
	. . I . .
	. . . . .
	. . . . .

scarecrow
	someone has broken down how crows work, not sure how to incoroprate
		https://www.reddit.com/r/StardewValley/comments/7bndnd/does_anyone_know_exactly_how_crows_work/
	8 tile radius
			. . . . . . . . .
	      . . . . . . . . . . .
	    . . . . . . . . . . . . .
	  . . . . . . . . . . . . . . .
	. . . . . . . . . . . . . . . . .
	. . . . . . . . . . . . . . . . .
	. . . . . . . . . . . . . . . . .
	. . . . . . . . . . . . . . . . .
	. . . . . . . . S . . . . . . . .
	. . . . . . . . . . . . . . . . .
	. . . . . . . . . . . . . . . . .
	. . . . . . . . . . . . . . . . .
	. . . . . . . . . . . . . . . . .
	  . . . . . . . . . . . . . . .
		. . . . . . . . . . . . .
	      . . . . . . . . . . .
		    . . . . . . . . .
*/

const FARM_WIDTH = 6
const FARM_HEIGHT = 6

var MINIMUM_CROPS = basicStrategy(FARM_WIDTH, FARM_HEIGHT)

const WORKERS = 8

const CHUNKING = 10

const FARM_SIZE = FARM_WIDTH * FARM_HEIGHT

// basicStrategy is a way for getting an easy estimate for the minimum number
// of crops that should be on the farm
func basicStrategy(width int, height int) int {
	return 28
	// greaterSide := width
	// lesserSide := height
	// if height > width {
	// 	greaterSide, lesserSide = height, width
	// }

	// var crops int
	// if lesserSide%3 == 0 {
	// 	crops = greaterSide * (lesserSide / 3) * 2
	// } else if lesserSide%3 == 1 {
	// 	crops = (greaterSide * int(math.Floor(float64(lesserSide)/3.0)) * 2) + greaterSide
	// } else if lesserSide%3 == 2 {
	// 	crops = (greaterSide * int(math.Ceil(float64(lesserSide)/3.0)) * 2)
	// }

	// return crops
}

func coordinatesFromIndex(width int, height int, index int) (int, int) {
	return index % width, int(math.Floor(float64(index) / float64(height)))
}

func indexFromCoordinates(width int, x int, y int) int {
	return (width * y) + x
}

func validLayout(width int, height int, layout []byte) bool {

	// Create a bigger size for the player to walk a circle around the plot
	biggerWidth := width + 2
	biggerHeight := height + 2
	biggerSize := biggerWidth * biggerHeight
	temp := make([]byte, biggerSize)

	cropsToWater := 0

	queue := make([]int, 1) // (biggerWidth*2)+(biggerHeight*2)-8
	queue[0] = 0
	walkingSpaces := 0

	// Building the new board...
	for i := 0; i < biggerSize; i++ {
		x, y := coordinatesFromIndex(biggerWidth, biggerHeight, i)
		if x == 0 || x == width+1 || y == 0 || y == height+1 {
			temp[i] = 'x'
		} else {
			temp[i] = layout[indexFromCoordinates(width, x-1, y-1)]
			if temp[i] == 'c' {
				cropsToWater++
			}

			if temp[i] == 'x' {
				walkingSpaces++
			}

			// Just go ahead and fill in everything as walkable
			if temp[i] == '.' {
				temp[i] = 'x'
			}
		}
	}

	if FARM_SIZE-walkingSpaces < MINIMUM_CROPS {
		return false
	}

	cropsWatered := 0

	visited := make([]bool, biggerSize)
	for i := 0; i < biggerSize; i++ {
		visited[i] = false
	}

	// Travel the board and try to mark every spot
	for len(queue) > 0 {
		if visited[queue[0]] == false {
			visited[queue[0]] = true
			x, y := coordinatesFromIndex(biggerWidth, biggerHeight, queue[0])

			// We have a left
			if x > 0 {
				otherIndex := indexFromCoordinates(biggerWidth, x-1, y)
				if temp[otherIndex] == 'c' {
					temp[otherIndex] = 'w'
					cropsWatered++
				} else if temp[otherIndex] == 'x' && visited[otherIndex] == false {
					queue = append(queue, otherIndex)
				}
			}

			// We have a right
			if x < biggerWidth-1 {
				otherIndex := indexFromCoordinates(biggerWidth, x+1, y)
				if temp[otherIndex] == 'c' {
					temp[otherIndex] = 'w'
					cropsWatered++
				} else if temp[otherIndex] == 'x' && visited[otherIndex] == false {
					queue = append(queue, otherIndex)
				}
			}

			// We have a bottom..?
			if y > 0 {
				otherIndex := indexFromCoordinates(biggerWidth, x, y-1)
				if temp[otherIndex] == 'c' {
					temp[otherIndex] = 'w'
					cropsWatered++
				} else if temp[otherIndex] == 'x' && visited[otherIndex] == false {
					queue = append(queue, otherIndex)
				}
			}

			// We have a top..?
			if y < biggerHeight-1 {
				otherIndex := indexFromCoordinates(biggerWidth, x, y+1)
				if temp[otherIndex] == 'c' {
					temp[otherIndex] = 'w'
					cropsWatered++
				} else if temp[otherIndex] == 'x' && visited[otherIndex] == false {
					queue = append(queue, otherIndex)
				}
			}

			// break as early as we can
			if cropsToWater == cropsWatered {
				return true
			}

		}
		queue = queue[1:]
	}
	return cropsToWater == cropsWatered
}

var basicOptions = []byte{'x', 'c'}

func expand(farm *Farm) []*Farm {

	for layoutIndex, layoutSelection := range farm.layout {

		if layoutSelection == '.' {

			options := append(basicOptions, farm.remainingResources.Options()...)
			expansion := make([]*Farm, len(options))
			added := 0

			for _, option := range options {
				newLayout := make([]byte, FARM_SIZE)
				copy(newLayout, farm.layout)
				newLayout[layoutIndex] = option
				if layoutIndex < (FARM_WIDTH*2)+2 || validLayout(FARM_WIDTH, FARM_HEIGHT, newLayout) {
					expansion[added] = NewFarm(farm.remainingResources, newLayout)
					added++
				}
			}
			return expansion[:added]
		}
	}

	return nil
}

func worker(id int, jobs <-chan []*Farm, results chan<- []*Farm) {
	for j := range jobs {
		r := make([]*Farm, 0)
		for _, f := range j {
			r = append(r, expand(f)...)
		}
		results <- r
	}
}

func main() {
	startingResources := NewResources(0, 0, 0, 0)
	farm := NewEmptyFarm(startingResources)

	bestFarm := farm

	iterations := 0

	start := time.Now()

	jobs := make(chan []*Farm, 10000000)
	jobs <- expand(farm)
	outstanding := 1

	results := make(chan []*Farm, 100)

	for w := 1; w <= WORKERS; w++ {
		go worker(w, jobs, results)
	}

	for outstanding > 0 {
		result := <-results
		iterations += len(result)
		outstanding--
		for _, j := range result {
			if j.Score() > bestFarm.Score() {
				bestFarm = j
			}
		}
		if len(result) > 0 {
			if len(result) > WORKERS*CHUNKING {
				jobsPerWorker := (len(result) / WORKERS) + 1
				for len(result) > 0 {
					if len(result) >= jobsPerWorker {
						jobs <- result[:jobsPerWorker]
						result = result[jobsPerWorker:]
					} else {
						jobs <- result
						result = result[len(result):]
					}
					outstanding++
				}
			} else {
				jobs <- result
				outstanding++
			}
		}
	}
	close(jobs)

	// for len(queue) > 0 {
	// 	if queue[0].Score() > bestFarm.Score() {
	// 		bestFarm = queue[0]
	// 	}
	// 	queue = append(queue, expand(queue[0])...)
	// 	iterations++
	// 	queue = queue[1:]
	// }

	elapsed := time.Since(start)
	log.Printf("Search took %s with %d workers; chunking %d", elapsed, WORKERS, CHUNKING)
	log.Printf("Explored %d solutions\n", iterations)
	log.Printf("Total Crops: %d\n", bestFarm.Score())
	log.Print(bestFarm.Render())
}
