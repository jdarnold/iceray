package main

import (
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

	err := gcfg.ReadFileInto(&IcerayCfg,*ConfigPath)
	if err != nil {
		log.Fatal("Error opening config file: "+err.Error())
	}

	if *Hostname == DEFAULT_HOSTNAME && IcerayCfg.Server.Hostname != "" {
		*Hostname = IcerayCfg.Server.Hostname
	}

	if *Port == DEFAULT_PORT && IcerayCfg.Server.Port != 0 {
		*Port = IcerayCfg.Server.Port
	}
	
	mountpoint := IcerayCfg.Server.Mount
	if mountpoint[0] != '/' {
		mountpoint = "/" + mountpoint
	}

	if *Mount == DEFAULT_MOUNTPOINT && IcerayCfg.Server.Mount != "" {
		*Mount = IcerayCfg.Server.Mount;
	}

	if *Username == DEFAULT_USERNAME && IcerayCfg.Server.User != "" {
		*Username = IcerayCfg.Server.User;
	}

	if *Password == DEFAULT_PASSWORD && IcerayCfg.Server.Password != "" {
		*Password = IcerayCfg.Server.Password;
	}
}



















