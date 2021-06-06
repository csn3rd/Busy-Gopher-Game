package main

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
	"strconv"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)


var RESET = "\033[0m"
var RED = "\033[31m"
var GREEN = "\033[32m"
var YELLOW = "\033[33m"
var BLUE = "\033[34m"
var PURPLE = "\033[35m"
var CYAN = "\033[36m"
var GRAY = "\033[37m"
var WHITE = "\033[97m"

func typeWrite(w io.Writer, sleep int, content string) {
	chars := strings.Split(content, "")

	for _, c := range chars {
		fmt.Fprint(w, c)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}

func colorType(w io.Writer, color string, sleep int, content string) {
	typeWrite(w, sleep, color+content+RESET)
}

func printEnd(w io.Writer, space int) {
	for i := 0; i < space; i++ {
		fmt.Fprint(w, "\n")
	}
}

func header(w io.Writer) {
	printEnd(w, 2)
	colorType(w, BLUE, 0, "                             #####    #####    #####   ###        #     #####    #####  \n                            #     #  #     #  #     #   #        ##    #     #  #     # \n                            #        #        #         #       # #    #        #     # \n                            #         #####   #         #         #    ######    ###### \n                            #              #  #         #         #    #     #        # \n                            #     #  #     #  #     #   #         #    #     #  #     # \n                             #####    #####    #####   ###      #####   #####    #####  ")
	printEnd(w, 3)
	colorType(w, BLUE, 0, "     #######  ###  #     #     #     #            ######   ######   #######        #  #######   #####   ####### \n     #         #   ##    #    # #    #            #     #  #     #  #     #        #  #        #     #     #    \n     #         #   # #   #   #   #   #            #     #  #     #  #     #        #  #        #           #    \n     #####     #   #  #  #  #     #  #            ######   ######   #     #        #  #####    #           #    \n     #         #   #   # #  #######  #            #        #   #    #     #  #     #  #        #           #    \n     #         #   #    ##  #     #  #            #        #    #   #     #  #     #  #        #     #     #    \n     #        ###  #     #  #     #  #######      #        #     #  #######   #####   #######   #####      #    ")
	printEnd(w, 3)
	colorType(w, PURPLE, 0, "               ____  _  _       __   ___   __  __  ___  ______ __  __  ___  __  __    __  __   ___  \n               || )) \\\\//       ||  // \\\\  ||\\ || // \\\\ | || | ||  || // \\\\ ||\\ ||    ||  ||  // \\\\ \n               ||=)   )/        || ((   )) ||\\\\|| ||=||   ||   ||==|| ||=|| ||\\\\||    ||==|| ((   ))\n               ||_)) //      |__||  \\\\_//  || \\|| || ||   ||   ||  || || || || \\||    ||  ||  \\\\_//")
	printEnd(w, 5)
	colorType(w, CYAN, 1, "               ____     __  __   _____ __  __          ______   ____     ____     __  __    ______    ____ \n              / __ )   / / / /  / ___/ \\ \\/ /         / ____/  / __ \\   / __ \\   / / / /   / ____/   / __ \\\n             / __  |  / / / /   \\__ \\   \\  /         / / __   / / / /  / /_/ /  / /_/ /   / __/     / /_/ /\n            / /_/ /  / /_/ /   ___/ /   / /         / /_/ /  / /_/ /  / ____/  / __  /   / /___    / _, _/ \n           /_____/   \\____/   /____/   /_/          \\____/   \\____/  /_/      /_/ /_/   /_____/   /_/ |_| ")
	printEnd(w, 1)
}

var tape = make([]string, 0)
var cards = 0
var symbols = 0

func initializeTape(size int) {
	tape = make([]string, 0)
	for i := 0; i < size; i++ {
		tape = append(tape, "0")
	}
}

func printTape(w io.Writer, pos int, crd string) {
	colorType(w, "\u001B[31m", 0, "  "+crd)
	for i := 0; i < 55-len(tape); i++ {
		colorType(w, "\u001b[31m", 0, " ")
	}
	for i := 0; i < len(tape); i++ {
		var col = "\u001B[31m"
		if (i == pos) {
			col += "\u001b[43m"
		} else if tape[i] == "1" {
			col += "\u001B[42m"
		}
		var delay = 0
		if i%4 == 0 {
			delay = 1
		}
		colorType(w, col, delay, tape[i])
		colorType(w, "", 0, " ")
	}
	printEnd(w, 1)
}

func scoreTape(w io.Writer) {
	var score = 0
	for i := 0; i < len(tape); i++ {
		if tape[i] == "1" {
			score++
		}
	}
	colorType(w, "", 0, "  Final Score: "+strconv.Itoa(score))
	printEnd(w, 1)
}

type card struct {
	id int
	overprints []int
	shifts []int
	transitions []int
}

func printCard(w io.Writer, c card) {
	colorType(w, "\033[4;36m", 0, fmt.Sprintf("%d:\n", c.id))
	for i := 0; i < len(c.overprints); i++ {
		colorType(w, "", 0, fmt.Sprintf("%d: %d%d%d\n", i, c.overprints[i], c.shifts[i], c.transitions[i]))
	}
	printEnd(w, 1)
}

func simulateTape(w io.Writer, pos int, crd card) (int, int) {
	printTape(w, pos, strconv.Itoa(crd.id))
	curr, _ := strconv.Atoi(tape[pos])
	tape[pos] = strconv.Itoa(crd.overprints[curr])
	nextpos := pos+1
	if crd.shifts[curr] == 0 {
		nextpos = pos-1
	}
	nextcard := crd.transitions[curr]
	return nextpos, nextcard
}

func game(w io.Writer, terminal *term.Terminal) int {
	printEnd(w, 2)
	log.Println("New Game")

	colorType(w, "", 0, "  How many cards would you like? ")
	t, _ := terminal.ReadLine()
	if t == ""{
		return -1
	}
	cards, _ = strconv.Atoi(t)
	if cards < 1 {
		return -1
	}

	colorType(w, "", 0, "  How many symbols would you like? ")
	t, _ = terminal.ReadLine()
	if t == ""{
		return -1
	}
	symbols, _ = strconv.Atoi(t)
	if symbols < 2 {
		return -1
	}

	printEnd(w, 1)

	// colorType(w, "", 0, "  How much tape would you like? ")
	// t, _ = terminal.ReadLine()
	// size, _ := strconv.Atoi(t)
	// printEnd(w, 1)

	size := 49

	shift := 200


	log.Println(fmt.Sprintf("Cards: %d     Symbols: %d     Tape Length: %d", cards, symbols, size))

	colorType(w, "", 0, "  Please enter transitions in the form \"overwite\" \"shift\" \"card transition\" with no spaces between them.\n")
	colorType(w, "", 0, "  Overwrite values must range from 0 to the number of symbols - 1.\n  Shift must be \"L\", \"l\", or \"0\" for left or \"R\", \"r\", or \"1\" for right.\n  Card transitions must be in the range from 0 to the number of cards - 1. -1 is the designated halting value.\n\n")

	var cardsList = make([]card, cards)
	for i := 0; i < cards; i++ {
		var printsList = make([]int, symbols)
		var shiftsList = make([]int, symbols)
		var transitionsList = make([]int, symbols)
		for j := 0; j < symbols; j++ {
			colorType(w, "", 0, fmt.Sprintf("  Transition for card %d on input %d: ", i, j))
			t, _ = terminal.ReadLine()
			if len(t) < 3 {
				return -1
			}
			transitionsList[j], _ = strconv.Atoi(t[2:])
			if transitionsList[j] >= cards || transitionsList[j] < -1 {
				return -1
			}
			if t[1:2] == "R" || t[1:2] == "r" || t[1:2] == "1" {
				shiftsList[j] = 1
			} else if t[1:2] == "L" || t[1:2] == "l" || t[1:2] == "0" {
				shiftsList[j] = 0
			} else {
				return -1
			}
			printsList[j], _ = strconv.Atoi(t[0:1])
			if printsList[j] >= symbols || printsList[j] < 0 {
				return -1
			}
		}
		cardsList[i] = card{i, printsList, shiftsList, transitionsList}
	}

	log.Println("\nCards:\n")
	for i := 0; i < cards; i++ {
		printCard(log.Writer(), cardsList[i])
	}

	printEnd(w, 3)

	initializeTape(size)
	pos := 24
	c := 0
	shiftScore := 0
	colorType(w, "", 0, "  ")
	colorType(w, "\033[4;31m", 0, "Card")
	colorType(w, "", 0, "                                                  ")
	colorType(w, "\033[4;31m", 0,"Tape")
	colorType(w, "", 0, "                                               ")
	printEnd(w, 2)
	for pos >= 0 && pos < size && c != -1 && shiftScore < shift {
		crd := cardsList[c]
		pos, c = simulateTape(w, pos, crd)
		shiftScore++
	}
	printTape(w, pos, "H")

	printEnd(w, 3)

	colorType(w, RED, 0, "  Final Tape:\n\n")
	printTape(w, -1, " ")
	printEnd(w, 1)

	if c == -1 {
		scoreTape(w)
		colorType(w, "", 0, fmt.Sprintf("  Number of shifts: %d\n", shiftScore))
		printEnd(w, 1)
	} else if shiftScore > shift {
		colorType(w, "", 0, fmt.Sprintf("  Final Score: 0 (Didn't halt within %d shifts)", shift))
		printEnd(w, 1)
		colorType(w, "", 0, fmt.Sprintf("  Number of shifts: >%d\n", shift))
		printEnd(w, 1)
	} else {
		colorType(w, "", 0, "  Final Score: 0 (Went beyond the tape)")
		printEnd(w, 1)
		colorType(w, "", 0, fmt.Sprintf("  Number of shifts: >%d\n", shiftScore))
		printEnd(w, 1)		
	}
	return 0
}

func footer(w io.Writer) {
	printEnd(w, 3)
	colorType(w, WHITE, 0, "              _______ _                 _           __                   _             _             _ \n             |__   __| |               | |         / _|                 | |           (_)           | |\n                | |  | |__   __ _ _ __ | | _____  | |_ ___  _ __   _ __ | | __ _ _   _ _ _ __   __ _| |\n                | |  | '_ \\ / _` | '_ \\| |/ / __| |  _/ _ \\| '__| | '_ \\| |/ _` | | | | | '_ \\ / _` | |\n                | |  | | | | (_| | | | |   <\\__ \\ | || (_) | |    | |_) | | (_| | |_| | | | | | (_| |_|\n                |_|  |_| |_|\\__,_|_| |_|_|\\_\\___/ |_| \\___/|_|    | .__/|_|\\__,_|\\__, |_|_| |_|\\__, (_)\n                                                                  | |             __/ |         __/ |  \n                                                                  |_|            |___/         |___/   ")
	printEnd(w, 2)
}

func main() {
	if runtime.GOOS == "windows" {
		RESET = ""
		RED = ""
		GREEN = ""
		YELLOW = ""
		BLUE = ""
		PURPLE = ""
		CYAN = ""
		GRAY = ""
		WHITE = ""
	}

	ssh.Handle(func(s ssh.Session) {
		log.Println("Started SSH session with " + s.User() + "\n")
		terminal := term.NewTerminal(s, "")
		header(s)
		var playing bool = true
		for playing {
			g := game(s, terminal)
			if g == -1 {
				footer(s)
				log.Println("Ended SSH session with " + s.User() + "\n")
				return
			}
			colorType(s, CYAN, 0, "  Would you like to play again? (yes/no) ");
			t, _ := terminal.ReadLine()
			if t == "" {
				footer(s)
				log.Println("Ended SSH session with " + s.User() + "\n")
				return
			}
			if t[0:1] != "y" && t[0:1] != "Y" {
				playing = false
			}
		}
		footer(s)
		log.Println("Ended SSH session with " + s.User() + "\n")
	})

	log.Println("Started SSH server on port 2222")
	ssh.ListenAndServe(":2222", nil)
	// log.Println("System: " + runtime.GOOS)
	// log.Fatal(ssh.ListenAndServe(":2222", nil))
}