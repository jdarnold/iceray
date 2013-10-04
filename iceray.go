// iceray project main.go
package iceray

import (
	"time"
	"math/rand"
	"fmt"
	"log"
	"github.com/systemfreund/go-libshout"
)

var randGen *rand.Rand

func main() {
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	
	readConfig()
	
	songs := getSongs()
	
	songCount := len(songs)

	if songCount == 0 {
		log.Fatal("No songs found!")
	}
	
	fmt.Printf("Found %d songs:\n", songCount)

	setupShout()
	
	defer closeShout()

	buffer := make([]byte, shout.BUFFER_SIZE)
		
	// maybe loop
	for {
		for songIdx := range(songs) {

			mfile := songs[songIdx]
			shoutSong(mfile, buffer)
		}

		if !IcerayCfg.Music.Loop {
			break
		}
	}
}

