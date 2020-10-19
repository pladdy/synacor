# synacor-challenge

[The challenge](https://challenge.synacor.com/)

Details that would be spoilers are intentionally left out.

## Getting Started

Start the Synacor challenge!

## The Journey

Spoilers below.  Codes have been omitted from the repo.

### Signing up

Code 1 is gained by simply signing up at the site.  Easy enough!

### Download the project and read the arch-spec

Code 2!  Read the arch-spec, you get a code, and a new task: implement a virtual
machine to this specification.

### VM Implemented

When the implemented VM runs the program downloaded with the project, a suite of
tests are run first.  If the VM functions correctly the tests will pass and
yield code 3.

### Playing the game

No not really, but the program, after the tests, turns out to be an RPG style
text game.  If you 'take tablet' when the program starts the game, you'll find
code 4 after you 'use tablet'.

### The grue cave

While exploring and playing the game, you'll find code 5 on the wall of a cave.

### The coin challenge

Solving this yields a teleporter.  When you take it and use it you get the 6th
code.

### Teleporter, part 1

Once you use the teleporter you arrive at a dead end.  You have to explore and
read a book to learn that the teleporter needs to be hacked...

### Teleporter, part 2

At this point in my endeavors, I got by with little debugging and introspection
into the VM.  Once I ran into the teleporter though, I needed to figure out how
to debug, and ultimately change the program so it would get me to code 7.

At this point...I became blocked.  I just hadn't this before.  So I did some
googling (the challenge is over/old now and the author says he doesn't give
hints any more on the synacor challenge site), and got some hints.

Ultimately the best hints and explanation came from [pankdm](http://pankdm.github.io/synacor-challenge.html).

They provide a walk through and good explanation.  Sadly, I couldn't make much
sense of their explanation on how they came up with their teleporter hack.  This
is my fault.  

I spent a long time on trying to get this code.  I added better debugging so
I could get output while the program ran and then I could see op codes, their
arguments, the registers, and the stack.

I had to kept going back for hints though; I just couldn't get my mind around
the program, how it was running in the VM and what optimizations I had to make.

After more hints...
-   [Redit thread](https://www.reddit.com/r/adventofcode/comments/3wyz4g/synacor_teleporter_challenge/)
-   [paiv](http://paiv.github.io/blog/2016/04/24/synacor-challenge.html)
-   [sparky8342](https://github.com/sparky8342/synacor_challenge)
-   [zach2good](https://github.com/zach2good/synacor-challenge/blob/master/vm.cpp)
-   [Ben Congdon](https://benjamincongdon.me/blog/2016/12/18/Taking-on-the-Synacor-Challenge/)
-   [Perl example](https://github.com/sparky8342/synacor_challenge/blob/master/teleport.pl)

I did get the program hacked to the point where I could disable the teleporter
check, set the 8th register, and get to the next portion of the game.  I even
got a code!  However that code was wrong.  Turns out you need the right value
in the 8th register to get the correct 7th code.

I ended up learning a lot but got blocked again.  I chose to use the [c++ implementation](https://github.com/pankdm/synacor-challenge/blob/master/teleport.cpp) pankdm provided (I wrote mine in Go) to make the calculations and
get the final result for the 8th register.

For more reading related to the function checking 8th register and optimizing
it:
-   [Ackermann Function](https://en.wikipedia.org/wiki/Ackermann_function)
-   [Ackermann function examples](https://rosettacode.org/wiki/Ackermann_function)

Tail call optimization might be needed for the teleporter function to be fast
enough.
-   [Functional Go](https://medium.com/@geisonfgfg/functional-go-bc116f4c96a4)
-   [Ackermann recursive & iterative](https://gist.github.com/Sebbyastian/9bf5551f915b2694c77e)

### Vault

Graphs make my brain hurt.  I brute forced this unsuccessfully by just trying
different combinations of movements to the vault.  I got close several times
but my solutions were too many steps.

I finally tried looking at [pankm's implementation](https://github.com/pankdm/synacor-challenge/blob/master/gdb.py)
to see how he did the breadth first search, and when I looked at the source code
I found the answer in an imported library...

Can't unsee that unfortunately.  In hindsight, I'm a little surprised I didn't
figure it out with my brute force tries AND I'm very sad I'm so bad at
algorithms I couldn't implement BFS and just find all the paths on my own.

TODO: Be better at graphs.

NOTE: the 'vault' solution is incomplete and bad.

## Developing the VM

-   Golang 1.14

### Installation / Setup

`make install`

## Usage

`make run`

### API

## Testing

`make test`

## Docs

`make docs`

## Versioning

Versioning is done using [Semver](https://semver.org/)

Tagging: TBD

## References
-   [Write a VM](https://justinmeiners.github.io/lc3-vm/)
-   [How to write go](https://golang.org/doc/code.html)
    -   it's been a while and I've forgotten...
-   [Reading a bit](https://stackoverflow.com/questions/29583024/reading-8-bits-from-a-reader-in-golang)
    -   [Example](https://play.golang.org/p/Wyr_K9YAro)
-   [encoding/binary](https://golang.org/pkg/encoding/binary/)
-   [Endianess](https://en.wikipedia.org/wiki/Endianness)
