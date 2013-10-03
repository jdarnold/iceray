// iceray project main.go
package main

import (
	"fmt"
	"os"
	"flag"
	"io"
	"path"
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

	Music struct {
		Folder []string
	}
}

func sdir(folder string, addfilechannel chan SongRecord) {
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

		if homefiles[i].IsDir() {
			ndir := folder+"/"+fname
			fmt.Println("Searching " + ndir)
			go sdir(ndir,addfilechannel)
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
		sr.title = "Unknown"
		sr.artist = "Unknown"
		
		addfilechannel <- sr
	}
}

func main() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg,"iceray.gcfg")
	if err != nil {
		panic("Config error: "+err.Error())
		
	}

	log.Println(cfg.Server)
	log.Println(cfg.Music)
	
	addfilechannel := make(chan SongRecord, 10)
	defer close(addfilechannel)

	for folderIdx :=  range(cfg.Music.Folder) {
		go sdir(cfg.Music.Folder[folderIdx],addfilechannel)
	}

	mountpoint := cfg.Server.Mount
	if mountpoint[0] != '/' {
		mountpoint = "/" + mountpoint
	}

	log.Println(cfg.Server)

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
		panic("Error opening server " + cfg.Server.Hostname + " : " + err.Error())
	}
	
	buffer := make([]byte, shout.BUFFER_SIZE)
	
	for {
		mfile := <- addfilechannel
		fd,err := os.Open(mfile.fullpath)
		defer fd.Close()
		
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

	}
}

















