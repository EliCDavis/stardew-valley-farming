# Stardew Valley Farming

Attempting to optimize crop placement in stardew valley.

I'm pretty sure this is an np problem.

This is an imense amount of different configurations that need to be checked. Moving this to cuda would be advantageous, atleast for the validation of a farming configuration. I'm running into memory issues, so something like c/++ is probably required. Any tips for figuring out a configuration is going to be sub optimal before exploring it further would be extremely advantagous.

# Problem Description

We have a player on a square grid. The player can move up, down, left, and right on the grid. They can not move in a diagonal. The player can also only water crops above, below, left, and right of them. They can not water in a diagonal. The player has a plot of land for farming represented as n by m matrix, where n and m are integers and greater than 0 and n can but not have to be equal to m. The goal is to maximize the number of crops on a nxm plot of land so that they can still be watered by the player. If a crop is not watered by the player it dies. The player has a plot of land for farming located in a plot of land they own. They can walk around the perimeter of farm land and enter their field at any point they'd like. This means a 5x5 grid is actually a 7x7 with the outer layer having no crops. At any cell in the grid it can be marked as either an `x`, which means it is free for the player to walk on, or a `c` where a crop has been planted and the player is not able to walk on.

## Medium Difficulty

Crops need to be surrounded by a scare crow or they can be potentially eaten by crows. The likelyhood of being eaten by crows changes with the number of crops present and is hard to deduce by the behaviour described [here](https://www.reddit.com/r/StardewValley/comments/7bndnd/does_anyone_know_exactly_how_crows_work/). Scarecrows can protect crops with an 8 cell radius.

## Hard Difficulty

There exists 3 different sprinkler systems in the game, each being better than the last. The cells they cover is dictated by `.` and a cell uncovered is dictated as `o`. If the sprinkler covers a given area then the player does not have to water it.

```
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
```

# Strategies

| symbol | description |
| -------| ----------- |
| x      | Clear area  |
| c      | Spot for a crop |

## Bruteforce

There's absilutely no way this is scalable. The only pruning going on is stopping the search if the board can not be traversed by the player. After being smart with a job system the time reduced from ~2m30s to ~30s on 8 core computer.

Running a 6x6 board causes me to run out of memory and have the program crash. Even if we had enough memory to run this, it still took a 5x5 30s to run. 6x6 - 5x5 = 36 - 25 = 11 new tiles to explore, each with the option to either be a crop or a square the player can walk in. That's 2^11 more solutions multiplied by the time it takes to explore the first 25. This ends up with 2^11 * 30 = 61440s = 1024m = 17 hours of running. 

![crash](https://i.imgur.com/Mea5qg6.png)

```
Search took 36.1950397s with 8 workers; chunking 10000
Explored 48408284 solutions
Total Crops: 21

c c x c c
c c x c c
c c x c c
c c x c c
c c c c c
```

## Pruning based on niave best guess layout

The next idea I came up with is using some strategy to come up with a pretty good layout, but not optimal. Then with this layout we know when exploring solutions, if it ever becomes impossible to become better than the guess, then we can stop exploring. The best guess algorithm for layout took the largest side, and would lay crops along it in pairs of 2, so if we had a 5x5 grid, then our best guess would return:

```
c c c c c
c c c c c 
x x x x x 
c c c c c 
c c c c c 
```

Which is 20 crops, close to the optimal 21 crops.

### Results

This actually ended up helping immensely, and allowed me to run a 6x6 for the first time.

5x5
```
Search took 148.9861ms with 8 workers; chunking 5000
Explored 31620 solutions
Total Crops: 21

c c x c c
c c x c c
c c x c c
c c x c c
c c c c c
```

6x6
```
Search took 7m37.0579751s with 8 workers; chunking 5000
Explored 384706620 solutions
Total Crops: 28

c x c c c c
c c c x c c
c c c x c c
x x x x x c
c c c c c c
c c c c c c
```

This current best guess does not scale however. If we could figure out the optimal number of crops before starting this run then it would help tremendousely, as all the program would have to do at that point is try placing them. After figuring out the optimal placement had 28 crops in a 6x6, I just made the best guess function return 28 to see the speed improvements. It went from taking 7m37 to JUST 2 SECONDS. Picking a number greater than 28 stops in less than 2 seconds with an incomplete board.

```
Search took 2.5670386s with 8 workers; chunking 5000
Explored 2102643 solutions
Total Crops: 28

c c c x c c
x c c x c c
c c c x c c
x x c x c c
c c c x c c
c c c c c c
```

I took this and hard coded the guess to 40 for a 7x7 board. The board never finishes so I'd knock down the hard coded guess by one until the board finished. As a result I finally got a valid board with 12 minutes of runtime.

7x7
```
Search took 12m17.9270004s with 8 workers; chunking 5000
Explored 547254570 solutions
Total Crops: 37

c c x c x c c
c c x c x x c
c c x c c c c
c c x c c c c
c x x x x x c
c c c c c c c
c c c c c c c
```

From this we see that even if we know the optimal number of crops for any board, coming up with a valid configuration will still take way too much time on larger boards. Which warrants the need for better strategies to be implemented along with this.

## Pruning bad sub-configurations

2019/01/03 19:02:04 Search took 4m43.4358307s with 8 workers; chunking 5000
2019/01/03 19:02:04 Explored 203492374 solutions2019/01/03 19:02:04 Total Crops: 37
2019/01/03 19:02:04
c c c x c c c
c x x x x x c
c c c c c c c
c c c c c c c
x x x x x x c
c c c c c c c
c c c c c c c

## Pruning based on uncovered regions

If a region is every left wholy empty then it can be shown that it can always be optimized. For example if a configuration contains a 3x3 region that is always empty, then we can always improve it by placing a crop in the center

```
x x x    should can be improved to    x x x
x x x  ============================>  x c x
x x x                                 x x x

```

I am making the assumption that this also applied to a 2x2 grid, but I can't prove it in my mind. However I have not seen a solution yet that ever contains a empty 2x2 grid in it. So I check to see if a grid every contains an empty 2x2 with no crops, and if it does then I prune it. Doing this reduced the 7x7 grid by 8 minutes (12 minutes to 4 minutes). But is far from what's needed if I want to start evaluating 8x8 grids.

```
Search took 4m43.4358307s with 8 workers; chunking 5000
Explored 203492374 solutions
Total Crops: 37

c c c x c c c
c x x x x x c
c c c c c c c
c c c c c c c
x x x x x x c
c c c c c c c
c c c c c c c
```