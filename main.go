package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fsnotify/fsnotify"
	"github.com/nickalie/go-webpbin"
)

var (
	dir      = flag.String("dir", "watch", "folder to watch for incoming images")
	existing = flag.Bool("convex", false, "convert already existing image in dir")
	quality  = flag.Uint("quality", 80, "quality of webp conversion between 1 and 100")
)

var ext = []string{".jpg", ".png"}

func main() {
	flag.Parse()
	figure.NewColorFigure("wellp", "isometric1", "blue", true).Print()

	if *existing {
		convExisting(*dir, *quality)
	}
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
					webp(event.Name, *quality)
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
func convExisting(dir string, quality uint) {
	imgs, _ := getImages(dir)
	for _, img := range imgs {
		err := webp(img, quality)
		if err != nil {
			log.Printf("could not convert %s to webp", img)
		}
	}
}
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

func webp(path string, quality uint) error {
	log.Printf("converting %s with quality %d", path, quality)
	err := webpbin.NewCWebP().
		Quality(quality).
		InputFile(path).
		OutputFile(getNewName(path)).
		Run()
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
func getImages(dir string) ([]string, error) {
	f, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	images := make([]string, 0)
	for _, fi := range f {
		if !fi.IsDir() && isValidFile(fi.Name()) {
			images = append(images, fmt.Sprintf("%s/%s", dir, fi.Name()))
		}
	}
	return images, nil
}
