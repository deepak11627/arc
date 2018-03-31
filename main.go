package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deepak11627/arc/arc"
	"github.com/deepak11627/arc/utils"
)

//CacheSize the number of maximum values to cache.
var CacheSize int

func main() {

	// Let's take cache size from user
	utils.Message("Please enter maximum number of keys which caching system should store. ")
	for CacheSize == 0 {
		SetCacheSize()
	}

	a := arc.NewARC(CacheSize)
	for { // Keep the program executing until user chooses to exit
		//prompt user to select an option
		option := showOptions()
		switch option {
		case 1:
			key := ReadCache()
			v, ok := a.Get(key)
			if ok {
				utils.Message(fmt.Sprintf("Value at %s is %s ", key, v))
			} else {
				utils.Message("No such key.")
			}
		case 2:
			// Send to LRU handler for option 1
			k, v := GetKeyValuePair()
			a.Put(k, v)
		case 3:
			a.Traverse()
		case 4:
			utils.Message("Thank you. Exiting...")
			os.Exit(0)
		default:
			utils.Message("Program error.")
			os.Exit(1)
		}
	}

}

func ReadCache() interface{} {
	utils.Message("Please enter key to read value from ")
	reader := bufio.NewReader(os.Stdin)
	k, _ := reader.ReadString('\n')
	k = strings.Replace(k, "\n", "", -1)
	return k
}

func GetKeyValuePair() (interface{}, interface{}) {
	utils.Message("Please enter key ")
	reader := bufio.NewReader(os.Stdin)
	k, _ := reader.ReadString('\n')
	k = strings.Replace(k, "\n", "", -1)
	utils.Message("Please enter value ")
	v, _ := reader.ReadString('\n')
	v = strings.Replace(v, "\n", "", -1)
	return k, v
}

//SetCacheSize takes input from user and sets value for CacheSize
func SetCacheSize() {
	reader := bufio.NewReader(os.Stdin)
	val, err := reader.ReadString('\n')
	if err != nil {
		utils.Message("Problem reading the entered value.")
		os.Exit(1)
	}
	val = strings.Replace(val, "\n", "", -1)
	ok, err := strconv.Atoi(val)
	if err != nil || ok < 0 {
		utils.Message("A positive number is the only accepted value.")
	} else {
		CacheSize = ok
	}
}

// showOptions prompts the user to select an option to proceed with the program
func showOptions() int {
	utils.RenderMessageHeading("Please select an operation to perform on cache.")
	utils.Message("Press 1 for getting a value from cache.")
	utils.Message("Press 2 for adding a value into cache.")
	utils.Message("Press 3 to view the cache items")
	utils.Message("Press 4 to Exit the program.")
	utils.RenderMessageEnd()
	notAnOption := true
	var selection int
	for notAnOption {
		reader := bufio.NewReader(os.Stdin)
		val, err := reader.ReadString('\n')
		if err != nil {
			utils.Message("Problem reading the entered value.")
			os.Exit(1)
		}
		val = strings.Replace(val, "\n", "", -1)
		selection, err = strconv.Atoi(val)
		if err != nil {
			utils.Message("1,2,3 or 4 are the only accepted values.")
		} else {
			if selection == 1 || selection == 2 || selection == 3 || selection == 4 {
				notAnOption = false
			} else {
				utils.Message("1,2,3 or 4 are the only accepted values.")
			}
		}

	}
	return selection
}
