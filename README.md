# Stardew Valley Farming

Attempting to optimize crop placement in stardew valley

I'm pretty sure this is an np problem.

This is an imense amount of different configurations that need to be checked. Moving this to cuda would be advantageous, atleast for the validation of a farming configuration. I'm running into memory issues, so something like c/++ is probably required. Any tips for figuring out a configuration is going to be sub optimal before exploring it further would be extremely advantagous.

Running a 6x6 board causes me to run out of memory and have the program crash

![crash](https://i.imgur.com/Mea5qg6.png)

## Results

| symbol | description |
| -------| ----------- |
| x      | Clear area  |
| c      | Spot for a crop |

### Runs
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