// iceray project main.go
package main

import (
	"math/rand"
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/systemfreund/go-libshout"
)

var randGen *rand.Rand

func isFileModified(path string, mtime *time.Time, subdirs bool) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Println("Problems getting modtime for " + path)
		return true
	}
	
	if *mtime != fi.ModTime() {
		*mtime = fi.ModTime()
		return true
	}

	return false
}

func isPlaylistModified( ) bool {
	rv := false
	for _, playlistRec :=  range(Playlists) {
		if isFileModified(playlistRec.path, &playlistRec.modtime, playlistRec.subdirs) {
			// need to go through them all to update modtimes
			rv = true
		}
	}

	return rv
}

func main() {
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	
	readConfig()
	
	setupShout()
	
	defer closeShout()

	buffer := make([]byte, shout.BUFFER_SIZE)

	songs := []SongRecord{}
	
	// maybe loop
	for {
		if len(songs) == 0 {
			// populate songs
			songs = getSongs()
	
			songCount := len(songs)
			
			if songCount == 0 {
				log.Fatal("No songs found!")
			}
	
			fmt.Printf("Found %d songs:\n", songCount)
		}

		for songIdx := range(songs) {
			if isPlaylistModified() {
				// empty songs to go get new list
				fmt.Println("Getting new songs")
				songs = []SongRecord{}
				break
			}
			
			mfile := songs[songIdx]
			shoutSong(mfile, buffer)

			time.Sleep(5 * time.Second)
		}

		if !IcerayCfg.Music.Loop {
			break
		}
	}
}

