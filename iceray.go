// iceray project main.go
package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"log"
	"github.com/jdarnold/go-id3"
)

type SongRecord struct {
	fullpath string
	filetype string
	title string
	artist string
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
	folder := "."

	if len(os.Args) > 1 {
		folder = os.Args[1]
	}
	
	addfilechannel := make(chan SongRecord, 10)
	defer close(addfilechannel)

	go sdir(folder,addfilechannel)

	for {
		
		mfile := <- addfilechannel
		fmt.Println(mfile)
	}
}










