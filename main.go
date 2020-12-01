package main

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/nickalie/go-webpbin"
)

var ext = []string{".jpg", ".png"}

func isValidFile(name string) bool {
	for _, e := range ext {
		if strings.Contains(name, e) {
			return true
		}
	}
	return false
}
func getNewName(name string) string {
	if strings.Contains(name, ".jpg") {
		return strings.Replace(name, "jpg", "webp", 1)
	}
	if strings.Contains(name, ".png") {
		return strings.Replace(name, "png", "webp", 1)
	}
	return name
}
func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Create && isValidFile(event.Name) {
					log.Println("processing:", event.Name)
					err := webpbin.NewCWebP().
						Quality(80).
						InputFile(event.Name).
						OutputFile(getNewName(event.Name)).
						Run()
					if err != nil {
						log.Fatal(err.Error())
					}
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("new file created:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("watch")
	if err != nil {
		log.Fatal(err)
	}
	<-done

}
