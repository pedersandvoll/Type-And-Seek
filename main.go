package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

type applicationMode int

const (
	LengthMode applicationMode = iota
	TimeMode
)

type TypeAndSeek struct {
	mode         applicationMode
	length       int
	time         int
	symbols      []string
	symbolsTyped int
	timeTaken    int
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  --help         Show this help message.")
	fmt.Println("  --length <arg> Set the number of words to type and seek.")
	fmt.Println("  --time   <arg> Set the time amount for typing and seeking.")
	os.Exit(0)
}

func createApplicationState() *TypeAndSeek {
	ts := &TypeAndSeek{}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help":
			printHelp()
		case "--length":
			ts.mode = LengthMode
			if len(os.Args) > 2 {
				length, err := strconv.Atoi(os.Args[2])
				if err != nil {
					ts.length = 15
				} else {
					ts.length = length
				}
			}
		case "--time":
			ts.mode = TimeMode
			if len(os.Args) > 2 {
				time, err := strconv.Atoi(os.Args[2])
				if err != nil {
					ts.time = 15
				} else {
					ts.time = time
				}
			}
		default:
			fmt.Printf("Unknown command: \033[31;1m%s\033[0m. Type \033[32;1m--help\033[0m for assistance.\n", os.Args[1])
			os.Exit(1)
		}
	} else {
		ts.mode = LengthMode
		ts.length = 15
	}

	ts.symbolsTyped = 0

	return ts
}

func (state *TypeAndSeek) ticker(startTime time.Time, done chan bool, ticker *time.Ticker) {
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				elapsed := time.Since(startTime)
				elapsedSeconds := math.Round(elapsed.Seconds())
				state.timeTaken = int(elapsedSeconds)

				if state.mode == TimeMode {
					if int(elapsedSeconds) < state.time {
						fmt.Printf("\033[1;0H\rTime: %.0f / %d seconds\n", elapsedSeconds, state.time)
					} else {
						fmt.Print("\033[H\033[2J")
						fmt.Println("\033[1;0H\rTime's up!")
						fmt.Printf("\033[2;0H\rSymbols typed: %d\n", state.symbolsTyped)
						done <- true
						return
					}
				} else {
					fmt.Printf("\033[1;0H\rTime: %s\n", elapsed.Round(time.Second).String())
				}
			}
		}
	}()
}

func insert(slice []string, index int, value string) []string {
	slice = append(slice, value)
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}

func (state *TypeAndSeek) startGame() {
	fmt.Print("\033[H\033[2J")
	startTime := time.Now()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	defer close(done)

	state.ticker(startTime, done, ticker)

	fmt.Printf("\033[2;0H\rMatch this symbol: %s\n", state.symbols[0])
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	keyChan := make(chan keyboard.Key)
	charChan := make(chan rune)
	errChan := make(chan error)

	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				errChan <- err
				return
			}
			charChan <- char
			keyChan <- key
		}
	}()

gameLoop:
	for {
		select {
		case <-done:
			break gameLoop
		case err := <-errChan:
			panic(err)
		case char := <-charChan:
			if string(char) == state.symbols[0] {
				if state.mode == TimeMode {
					randomIndex := rand.Intn(len(state.symbols))
					state.symbols = insert(state.symbols, randomIndex, state.symbols[0])
				}
				state.symbols = state.symbols[1:]
				state.symbolsTyped = state.symbolsTyped + 1
				if len(state.symbols) > 0 {
					fmt.Printf("\033[2;0H\rMatch this symbol: %s\n", state.symbols[0])
				} else {
					fmt.Print("\033[H\033[2J")
					fmt.Println("\033[1;0H\rAll symbols typed!")
					fmt.Printf("\033[2;0H\rTime taken: %d\n", state.timeTaken)
					break gameLoop
				}
			}
		case key := <-keyChan:
			if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
				break gameLoop
			}
		}
	}
	fmt.Println("\033[3;0H\rExiting Game...")
	os.Exit(0)
}

func (state *TypeAndSeek) createKeyOrder(f *os.File) {
	arr := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 1 {
			arr = append(arr, line)
		}
	}

	if state.mode == LengthMode {
		for len(arr) < state.length {
			randomIndex := rand.Intn(len(arr))
			pick := arr[randomIndex]
			arr = append(arr, pick)
		}
	}

	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}

	if state.mode == LengthMode {
		if len(arr) > state.length {
			arr = arr[:state.length]
		}
	}

	state.symbols = arr
}

func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filename, err)
	}
	return file, nil
}

func (state *TypeAndSeek) runExampleInput() {
	exampleInput, err := openFile("example-input.txt")
	if err != nil {
		fmt.Println("\033[31;1mexample-input.txt\033[0m not found.")
		os.Exit(1)
	}
	defer exampleInput.Close()

	state.createKeyOrder(exampleInput)

	state.startGame()
}

func main() {
	state := createApplicationState()

	input, err := openFile("input.txt")
	if err != nil {
		state.runExampleInput()
		return
	}
	defer input.Close()

	state.createKeyOrder(input)

	fmt.Println("")

	state.startGame()
}
