# Simple Classical XY Model

A simple simulation of the Classical XY model, written in go, visualised with ebiten.

Use up/down arrow keys to increase or decrease temperature respectively. Also use right/left to increase or decrease interaction strength respectively.

You can see a more interactive version of this on my [website](https://e74000.net/posts/xymodel/).

### Command line flags:

* `-t`: The initial temperature of the simulation.
* `-e`: The value for the external field
* `-i`: The initial interaction strength
* `-x`: X resolution
* `-y`: Y resolution
* `-s`: Scale (only works if not fullscreen)
* `-f`: Enable fullscreen
