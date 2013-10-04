// iceray project main.go
package main

import (
	"time"
	"math/rand"
	"sync"
	"fmt"
	"os"
	"os/user"
	"flag"
	"io"
	"path"
	"path/filepath"
	"strings"
	"log"
	"github.com/ascherkus/go-id3/src/id3"
	"github.com/systemfreund/go-libshout"
	"code.google.com/p/gcfg"
)

type SongRecord struct {
	fullpath string
	filetype string
	title string
	artist string
}

// Setup some command line flags
type Config struct  {
	Server struct {
		Hostname string
		Port uint
		User string
		Password string
		Mount string
	}

	Music map[string] *struct {
		Playlist string
		Shuffle bool
		Subdirs bool
		Rootfolder string
	}
}

func sdir(folder string, subdirs bool, addfilechannel chan SongRecord, w *sync.WaitGroup) {
	defer w.Done()

	searchdir, eopen := os.Open(folder)
	if eopen != nil {
		log.Println("Error opening " + folder + " : " + eopen.Error())
		return
	}
	
	homefiles, eread := searchdir.Readdir(-1)
	if eread != nil {
		log.Println("Error reading " + folder + " : " + eopen.Error())
		return
	}

	for i := range homefiles {
		fname := homefiles[i].Name()
		if fname[0] == '.' {
			continue
		}

		if homefiles[i].IsDir() && subdirs {
			ndir := folder+"/"+fname
			w.Add(1)
			go sdir(ndir,subdirs,addfilechannel,w)
			continue
		}

		if !strings.Contains(fname,".mp3") {
			continue
		}
			
		if homefiles[i].Size() < 100 {
			continue
		}

		var fext = path.Ext(fname);

		var sr SongRecord
		sr.fullpath = folder+"/"+fname
		sr.filetype = strings.TrimPrefix(fext,".")

		addfilechannel <- sr
	}
}

func main() {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	configPath := usr.HomeDir + "/.iceray.gcfg"

	var cfg Config
	err = gcfg.ReadFileInto(&cfg,configPath)
	if err != nil {
		log.Fatal("Error opening config file: "+err.Error())
	}

	addfilechannel := make(chan SongRecord, 100)

	var w sync.WaitGroup

	for _, musicRec :=  range(cfg.Music) {
		fext := strings.ToLower(filepath.Ext(musicRec.Playlist))
		if fext == ".xspf" {
			// process XML playlist file
		} else if fext == ".m3u" {
			// process m3u playlist file
		} else {
			w.Add(1)
			go sdir(musicRec.Playlist,musicRec.Subdirs, addfilechannel, &w)
		}
	}

	// wait for song search to finish up
	w.Wait()
	close(addfilechannel)
	
	var songs []SongRecord
	
	for mfile := range addfilechannel {
		songs = append(songs,mfile)
	}
	
	songCount := len(songs)
	log.Printf("Found %d songs", songCount)
	
	log.Println(songs)
	
	// Now shuffle it
	for i := range(songs) {
		j := i + randGen.Intn(songCount-i)
		tmp := songs[i]
		songs[i] = songs[j]
		songs[j] = tmp
	}

	log.Println(songs)
	
	mountpoint := cfg.Server.Mount
	if mountpoint[0] != '/' {
		mountpoint = "/" + mountpoint
	}

	log.Printf("Connecting to %s:%d",cfg.Server.Hostname, cfg.Server.Port)
	
	hostname := flag.String("host", cfg.Server.Hostname, "shoutcast server name")
	port := flag.Uint("port", cfg.Server.Port, "shoutcast server source port")
	user := flag.String("user", cfg.Server.User, "source user name")
	password := flag.String("password", cfg.Server.Password, "source password")
	mount := flag.String("mountpoint", mountpoint, "mountpoint")

	flag.Parse()

	// Setup libshout parameters
	s := shout.Shout{
		Host:     *hostname,
		Port:     *port,
		User:     *user,
		Password: *password,
		Mount:    *mount,
		Format:   shout.FORMAT_MP3,
		Protocol: shout.PROTOCOL_HTTP,
	}

	defer s.Close()

	// Create a channel where we can send the data
	stream, err := s.Open()
	if err != nil {
		log.Fatal("Error opening server " + cfg.Server.Hostname + " : " + err.Error())
	}
	
	buffer := make([]byte, shout.BUFFER_SIZE)
	
	
	for {
		if len(songs) == 0 {
			break
		}

		songIdx := randGen.Intn(len(songs))
		mfile := songs[songIdx]
		
		fd,err := os.Open(mfile.fullpath)
		
		if err != nil {
			log.Println("Problem opening: " + mfile.fullpath)
			continue
		}
		
		// Read in MP3 tags
		mp3tags := id3.Read(fd)

		if ( mp3tags == nil ) {
			log.Println("Problems getting MP3 tags for " + mfile.fullpath)
			continue
		}

		if mp3tags.Artist == "" {
			log.Println("Artist tag missing for " + mfile.fullpath)
			continue
		}

		if  mp3tags.Name == "" {
			log.Println("Song tag missing for " + mfile.fullpath)
			continue
		}
		
		mfile.artist = mp3tags.Artist
		mfile.title = mp3tags.Name

		track := mfile.title + " by " + mfile.artist
		fmt.Println("Playing " + track)

		fd.Seek(0,0)
		
		// add track to the stream
		s.UpdateMetadata( "song", track )
		
		for {
			// Read from file
			n, err := fd.Read(buffer)
			if err != nil && err != io.EOF { panic(err) }
			if n == 0 { break }

			// Send to shoutcast server
			stream <- buffer
		}

		fd.Close()
	}
}
