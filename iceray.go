// iceray project main.go
package main

import (
	"math/rand"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"
	"flag"
	
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

var DEFAULT_HOSTNAME="localhost"
var DEFAULT_PORT uint =8000
var DEFAULT_USERNAME="source"
var DEFAULT_PASSWORD=""
var DEFAULT_MOUNTPOINT="/steam.mp3"

var Hostname, Username, Password, Mount, ConfigPath *string
var Port *uint

func parseCommandline() {
	Hostname = flag.String("host", DEFAULT_HOSTNAME, "shoutcast server name")
	Port = flag.Uint("port", DEFAULT_PORT, "shoutcast server source port")
	Username = flag.String("user", DEFAULT_USERNAME, "source user name")
	Password = flag.String("password", DEFAULT_PASSWORD, "source password")
	Mount = flag.String("mountpoint", DEFAULT_MOUNTPOINT, "server mountpoint")

	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	cpath := usr.HomeDir + "/.iceray.gcfg"
	ConfigPath = flag.String("config", cpath, "full path to config file")

	flag.Parse()
}

func main() {

	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	
	parseCommandline()
		
	readConfig()
	
	setupShout()
	
	defer closeShout()

	buffer := make([]byte, shout.BUFFER_SIZE)

	songs := []SongRecord{}
	
	// maybe loop
	for {
		if len(songs) == 0 {
			// (re)populate songs
			songs = getSongs()
	
			songCount := len(songs)
			if songCount == 0 {
				log.Fatal("No music found!")
			}
	
			fmt.Printf("Found %d songs:\n", songCount)
		}

		for songIdx := range(songs) {
			if isPlaylistModified() {
				// empty songs and go get new list
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

