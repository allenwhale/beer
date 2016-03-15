package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	last_mod := time.Now()
	has_event := true
	ext := ".swp"
	appName := "app"
	usrName := "tachien"
	rootDir := appName + "/src/" + usrName + "/" + appName

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if time.Now().Sub(last_mod) > time.Second {
					if !strings.Contains(event.Name, ext) {
						// log.Println(event)
						has_event = true
					}
				}
				last_mod = time.Now()
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for _ = range ticker.C {
			if has_event {
				log.Println("restart")
				has_event = false
			}
		}
	}()
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != rootDir {
			err = watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
