# Type And Seek

A simple terminal-based typing game where you match displayed symbols with keyboard input.

## Description

Type And Seek is a minimalist typing game that challenges you to match a series of symbols as quickly as possible. Each time you correctly match a symbol, a new one appears until you complete all rounds. The game supports two modes of play:

- **Length Mode**: Play through a specific number of symbols
- **Time Mode**: Play for a set amount of time

## Installation

### Prerequisites

- Go (1.16 or later recommended)
- Terminal with ANSI color support
- The [github.com/eiannone/keyboard](https://github.com/eiannone/keyboard) package

### Installing

1. Clone this repository
```
git clone https://github.com/pedersandvoll/Type-And-Seek.git
cd type-and-seek
```

2. Install dependencies
```
go mod tidy
```

3. Build the program
```
go build
```

## Usage

Run the program with one of the following commands:

```
# Default: Run in Length Mode with 15 symbols
./type-and-seek

# Set a specific length (number of symbols)
./type-and-seek --length 20

# Run in Time Mode for a specific duration (in seconds)
./type-and-seek --time 30

# Display help
./type-and-seek --help
```

## Input Files

The program reads symbols from input files:

- **input.txt**: Create this file with one symbol per line to use your custom symbols
- **example-input.txt**: A fallback file used if input.txt is not found

Each symbol should be a single character on its own line.

Example input.txt:
```
a
b
c
d
e
f
g
h
i
j
```

## Controls

- Type the displayed symbol to progress
- Press `Esc` or `Ctrl+C` to exit the game

## Features

- Two game modes: Length Mode and Time Mode
- Performance tracking (time taken in Length Mode, symbols typed in Time Mode)
- Colorful terminal output for better visibility
- Randomized symbol order for varied gameplay
- Support for custom symbol sets

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
