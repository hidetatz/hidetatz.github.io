package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/atotto/clipboard"
	"github.com/fsnotify/fsnotify"
)

func runServer() {
	ctx := context.Background()

	port := "8080"
	directory := "./docs"
	ip := getLocalIP()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	server := &http.Server{Addr: ":8080", Handler: http.FileServer(http.Dir(directory))}

	go func() {
		server.ListenAndServe()
	}()

	fmt.Printf("Serving at %s:%s (the address is automatically copied into the clipboard)\n", ip, port)
	clipboard.WriteAll(fmt.Sprintf("%s:%s", ip, port))

	// live reload with fsnotify
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write != fsnotify.Write {
					continue
				}

				fmt.Println("files edited. hot reloading...")
				server.Shutdown(ctx)
				time.Sleep(time.Second) // wait for shutdown

				removeAllFiles("./docs/")
				gen()
				server = &http.Server{Addr: ":8080", Handler: http.FileServer(http.Dir(directory))}
				// don't take care of goroutine cancel since this is not a long-running script and
				// want to make code look easy
				go func() {
					server.ListenAndServe()
				}()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	watcher.Add("./data/articles/ja")
	watcher.Add("./data/articles/")
	watcher.Add("./")
	<-done
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
