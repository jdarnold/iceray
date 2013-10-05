iceray
======

Icecast/Shoutcast client. Sends MP3 files to the specified server, via the Shoutcast protocol

Copy icecast.gcfg.sample to $HOME/.icecast.gcfg and edit parameters to reflect your server and music sources. 

Then simply run it:

```
$ ./icecast
```

Requirements
--

* [go-libshout](https://github.com/systemfreund/go-libshout)
* [libshout](http://www.icecast.org/download.php)
* [go-id3](https://github.com/ascherkus/go-id3)
* [gcfg](http://code.google.com/p/gcfg/)
