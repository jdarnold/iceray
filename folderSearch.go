package main

import (
	"sync"
	"log"
	"path"
	"strings"
	"os"
)

func folderSearch(folder string, subdirs bool, addfilechannel chan SongRecord, w *sync.WaitGroup) {
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

	// get mod time for later check
	fi, err := os.Stat(folder)
	if err != nil {
		log.Println("Problems getting modtime for " + folder)
		return
	}

	var playlist Playlist
	playlist.path = folder
	playlist.modtime = fi.ModTime()
	playlist.subdirs = subdirs
	Playlists = append(Playlists, &playlist)
	for i := range homefiles {
		fname := homefiles[i].Name()
		if fname[0] == '.' {
			continue
		}

		if homefiles[i].IsDir() && subdirs {
			ndir := folder+"/"+fname
			w.Add(1)
			go folderSearch(ndir,subdirs,addfilechannel,w)
			continue
		}

		if !strings.HasSuffix(fname,".mp3") {
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



