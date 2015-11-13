package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Store map[string]string
type Count map[string]int
type Block struct {
	store Store
	count Count
}

// Store the current values
var store = Store{}

// Store counts for fast NUMEQUALTO
var count = Count{}

// Store previous values/counts for each block, so they can be rolled back
var blocks = []*Block{}

/*
	Utility
*/
func cache(key string, oldVal string, newVal string) {
	if len(blocks) > 0 {
		currBlock := *blocks[len(blocks)-1]
		// Make sure not to overwrite the stored value
		if _, found := currBlock.store[key]; !found {
			currBlock.store[key] = oldVal
		}
		// Empty string values are always found, so will never be cached
		if _, found := currBlock.count[oldVal]; !found {
			currBlock.count[oldVal] = count[oldVal]
		}
		if _, found := currBlock.count[newVal]; !found {
			currBlock.count[newVal] = count[newVal]
		}
	}
}

/*
	Operations
*/
func set(key string, val string) {
	if key != "" && val != "" {
		oldVal := store[key]
		cache(key, oldVal, val)
		store[key] = val
		if count[oldVal] > 0 {
			count[oldVal] = count[oldVal] - 1
		}
		count[val] = count[val] + 1
	}

}

func unset(key string) {
	if key != "" {
		oldVal := store[key]
		cache(key, oldVal, "")
		delete(store, key)
		if count[oldVal] > 0 {
			count[oldVal] = count[oldVal] - 1
		}
	}
}

func get(key string) {
	if key != "" {
		if val, ok := store[key]; ok && val != "" {
			fmt.Println(val)
		} else {
			fmt.Println("NULL")
		}
	}
}

func begin() {
	blocks = append(blocks, &Block{Store{}, Count{}})
}

func rollback() {
	if len(blocks) > 0 {
		currBlock := *blocks[len(blocks)-1]
		for sk, sv := range currBlock.store {
			store[sk] = sv
		}
		for ck, cv := range currBlock.count {
			count[ck] = cv
		}
		blocks = blocks[:len(blocks)-1]
	} else {
		fmt.Println("NO TRANSACTION")
	}
}

func commit() {
	if len(blocks) > 0 {
		blocks = []*Block{}
	} else {
		fmt.Println("NO TRANSACTION")
	}
}

/*
	Command line
*/
func parseArgs(args []string) (cmd string, key string, val string) {
	pos := [2]string{}
	copy(pos[:], args[1:])
	cmd = strings.ToUpper(args[0])
	key, val = strings.TrimSpace(pos[0]), strings.TrimSpace(pos[1])
	return cmd, key, val
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
Loop:
	for scanner.Scan() {
		args := strings.SplitN(scanner.Text(), " ", 3)
		cmd, key, val := parseArgs(args)
		switch cmd {
		case "SET":
			set(key, val)
		case "UNSET":
			unset(key)
		case "GET":
			get(key)
		case "NUMEQUALTO", "NEQ":
			fmt.Println(count[key])
		case "BEGIN":
			begin()
		case "ROLLBACK":
			rollback()
		case "COMMIT":
			commit()
		case "END":
			break Loop
		default:
			fmt.Println("INVALID COMMAND")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "I/0 ERROR:", err)
	}
	fmt.Println("HAVE A NICE DAY")
}
