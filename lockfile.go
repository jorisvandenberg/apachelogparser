package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var lockfile *os.File
var lockfileMutex sync.Mutex

func lockUnlock(lock bool, args Args) {
	// Acquire the lockfile mutex to ensure that only one instance of the program
	// can create or release the lockfile at a time.
	lockfileMutex.Lock()
	defer lockfileMutex.Unlock()

	if lock {
		// Get the temporary directory for the current operating system.
		tmpDir := os.TempDir()

		// Create the lock file path.
		lockfilePath := filepath.Join(tmpDir, "apachelogparser_"+strconv.FormatUint(HashStruct(args), 10)+".lock")

		// Open the lock file for writing, creating it if it doesn't exist.
		f, err := os.OpenFile(lockfilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			fmt.Println("Failed to create lock file:", err)
			os.Exit(1)
		}

		// Set the global lockfile variable.
		lockfile = f
	} else {
		// Close the lockfile.
		err := lockfile.Close()
		if err != nil {
			fmt.Println("Failed to close lock file:", err)
			os.Exit(1)
		}

		// Remove the lockfile.
		err = os.Remove(lockfile.Name())
		if err != nil {
			fmt.Println("Failed to remove lock file:", err)
			os.Exit(1)
		}
	}
}

func HashStruct(s Args) uint64 {
	// Convert struct to JSON byte slice
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	// Create FNV-1a hash object
	h := fnv.New64a()

	// Write JSON byte slice to hash object
	_, err = h.Write(b)
	if err != nil {
		panic(err)
	}

	// Return hash as uint64
	return h.Sum64()
}
