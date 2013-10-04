package iceray

import (
	"os/user"
	"code.google.com/p/gcfg"
	"log"
)

type Config struct  {
	Server struct {
		Hostname string
		Port uint
		User string
		Password string
		Mount string
	}

	Music struct {
		Shuffle bool
		Loop bool
	}
	
	Playlist map[string] *struct {
		Path string
		Subdirs bool
		Rootfolder string
	}
}

var IcerayCfg Config

func readConfig() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	configPath := usr.HomeDir + "/.iceray.gcfg"

	err = gcfg.ReadFileInto(&IcerayCfg,configPath)
	if err != nil {
		log.Fatal("Error opening config file: "+err.Error())
	}

}
