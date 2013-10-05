package main

import (
	"path/filepath"
	"sync"
	"strings"
	"github.com/ascherkus/go-id3/src/id3"
	"os"
	"log"
	"time"
)

type SongRecord struct {
	fullpath string
	filetype string
	title string
	artist string
}

type Playlist struct {
	path string
	modtime time.Time
	subdirs bool
}

var Playlists []*Playlist

func getSongs() []SongRecord {
	addfilechannel := make(chan SongRecord, 100)


	var w sync.WaitGroup

	for _, playlistRec :=  range(IcerayCfg.Playlist) {

		fext := strings.ToLower(filepath.Ext(playlistRec.Path))
		if fext == ".xspf" {
			// process XML playlist file
		} else if fext == ".m3u" {
			// process m3u playlist file
		} else {
			w.Add(1)
			go folderSearch(playlistRec.Path, playlistRec.Subdirs, addfilechannel, &w)
		}
	}

	// wait for song search(es) to finish up
	w.Wait()
	close(addfilechannel)
	
	var songs []SongRecord
	
	for mfile := range addfilechannel {

		// Read in MP3 info and save it
		
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

		songs = append(songs,mfile)
	}

	if IcerayCfg.Music.Shuffle {
		// Now (linear) shuffle it
		sc := len(songs)
		for i := range(songs) {
			j := randGen.Intn(sc)
			songs[i], songs[j] = songs[j], songs[i]
		}
	}

	return songs
}










