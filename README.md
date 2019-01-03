# Stardew Valley Farming

Attempting to optimize crop placement in stardew valley

I'm pretty sure this is an np problem.

This is an imense amount of different configurations that need to be checked. Moving this to cuda would be advantageous, atleast for the validation of a farming configuration. I'm running into memory issues, so something like c/++ is probably required. Any tips for figuring out a configuration is going to be sub optimal before exploring it further would be extremely advantagous.

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

Which is 20 crops, close to the optimal 21 crops. This current best guess does not scale however. If we could figure out the optimal number of crops before starting this run then it would help tremendousely, as all the program would have to do at that point is try placing them.

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