package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//CacheSize the number of maximum values to cache.
var CacheSize int

// LRUCache is Least recently Used cache
type LRUCache struct {
	*list.List
	Count int
	Heap  map[string]string
}

// Option are func to dynamically set LRUCache values
type Option func(*LRUCache)

// NewLRUCache create a new LRUCache cache
func NewLRUCache(opts ...Option) *LRUCache {

	lru := &LRUCache{
		Count: 0,
		List:  list.New(),
		Heap:  make(map[string]string, 0),
	}
	for _, opt := range opts {
		opt(lru)
	}

	return lru
}

func (c *LRUCache) Traverse() {
	var items []interface{}
	// Iterate through list and print its contents.
	for e := c.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}

	RenderMessageHeading("Current cached items are")
	// Iterate through list and print its contents.
	for _, item := range items {
		switch item.(type) {
		case int:
			fmt.Printf("  | %d |  ", item.(int))
		case float32:
			fmt.Printf("  | %f |  ", item.(float32))
		default:
			fmt.Print(fmt.Sprintf("  | %s  |  ", item))
		}
	}

	RenderMessageEnd()
}

func main() {

	// Let's take cache size from user
	Message("Please enter maximum number of keys which caching system should store. ")
	for CacheSize == 0 {
		SetCacheSize()
	}

	// we have cache size, now prompt user to select an option
	option := showOptions()
	switch option {
	case 1:
		// Send to LRU handler for option 1
		Message("you selected 1")
	case 2:
		Message("you selected 2")
	case 3:
		Message("Thank you. Exiting...")
		os.Exit(0)
	default:
		Message("Program error.")
		os.Exit(1)
	}

}

//SetCacheSize takes input from user and sets value for CacheSize
func SetCacheSize() {
	reader := bufio.NewReader(os.Stdin)
	val, err := reader.ReadString('\n')
	if err != nil {
		Message("Problem reading the entered value.")
		os.Exit(1)
	}
	val = strings.Replace(val, "\n", "", -1)
	ok, err := strconv.Atoi(val)
	if err != nil || ok < 0 {
		Message("A positive number is the only accepted value.")
	} else {
		CacheSize = ok
	}
}

// showOptions prompts the user to select an option to proceed with the program
func showOptions() int {
	RenderMessageHeading("Please select an operation to perform on cache.")
	Message("Press 1 for getting a value from cache.")
	Message("Press 2 for adding a value into cache.")
	Message("Press 3 to Exit the program.")
	RenderMessageEnd()
	notAnOption := true
	var selection int
	for notAnOption {
		reader := bufio.NewReader(os.Stdin)
		val, err := reader.ReadString('\n')
		if err != nil {
			Message("Problem reading the entered value.")
			os.Exit(1)
		}
		val = strings.Replace(val, "\n", "", -1)
		selection, err = strconv.Atoi(val)
		if err != nil {
			Message("1,2 or 3 are the only accepted values.")
		} else {
			if selection == 1 || selection == 2 || selection == 3 {
				notAnOption = false
			} else {
				Message("1,2 or 3 are the only accepted values.")
			}
		}

	}
	return selection
}
