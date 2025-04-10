package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"

	"github.com/eiannone/keyboard"
)

type applicationMode int

const (
	LengthMode applicationMode = iota
	TimeMode
)

type TypeAndSeek struct {
	mode    applicationMode
	length  int
	time    int
	symbols []string
	mu      sync.Mutex
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
					fmt.Println("Invalid argument for --length. Please provide a number.")
				} else {
					ts.length = length
					fmt.Printf("Setting the number of words to type and seek to \033[32;1m%d\033[0m.\n", length)
				}
			} else {
				fmt.Println("No argument provided for --length. Defaulting to \033[32;1m15\033[0m rounds.")
			}
		case "--time":
			ts.mode = TimeMode
			if len(os.Args) > 2 {
				time, err := strconv.Atoi(os.Args[2])
				if err != nil {
					ts.time = 15
					fmt.Println("Invalid argument for --time. Please provide a number.")
				} else {
					ts.time = time
					fmt.Printf("Setting the time amount to type and seek to \033[32;1m%d\033[0m seconds.\n", time)
				}
			} else {
				fmt.Println("No argument provided for --time. Defaulting to \033[32;1m15\033[0m seconds.")
			}
		default:
			fmt.Printf("Unknown command: \033[31;1m%s\033[0m. Type \033[32;1m--help\033[0m for assistance.\n", os.Args[1])
		}
	} else {
		ts.mode = LengthMode
		ts.length = 15
		fmt.Println("No command provided. Running the type and seek for \033[32;1m15\033[0m rounds.")
	}

	fmt.Println("")

	return ts
}

func (state *TypeAndSeek) startGame() {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("\033[2;0H\rMatch this symbol: %s\n", state.symbols[0])
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if string(char) == state.symbols[0] {
			state.symbols = state.symbols[1:]
			if len(state.symbols) > 0 {
				fmt.Print("\033[H\033[2J")
				fmt.Printf("\033[2;0H\rMatch this symbol: %s\n", state.symbols[0])
			} else {
				fmt.Print("\033[H\033[2J")
				fmt.Println("\033[2;0H\rAll symbols matched! Game completed.")
				break
			}
		}
		if key == keyboard.KeyEsc {
			continue
		}
		if key == keyboard.KeyCtrlC {
			break
		}
	}
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

	fmt.Println("")

	if state.mode == TimeMode {
		state.startGame()
	}
}

func main() {
	state := createApplicationState()

	input, err := openFile("input.txt")
	if err != nil {
		fmt.Println("No \033[31;1minput.txt\033[0m file found. Please create an \033[32;1minput.txt\033[0m file and add your own keys to 'Type And Seek'.")
		fmt.Println("Trying \033[32;1mexample-input.txt\033[0m instead...")
		state.runExampleInput()
		return
	}
	defer input.Close()

	state.createKeyOrder(input)
}
