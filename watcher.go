package main

import (
	"log"

	"github.com/hpcloud/tail"
)

const (
	tailOffset = 10
	tailWhence = 2
)

// Watcher watches file updates
type Watcher struct {
	hub   *Hub
	tail  *tail.Tail
	cache []string
}

// NewWatcher is Watcher constructor
func NewWatcher(hub *Hub, file string) (*Watcher, error) {
	config := tail.Config{
		Location: &tail.SeekInfo{Offset: tailOffset, Whence: tailWhence},
		ReOpen:   true,
		Poll:     true,
		Follow:   true,
	}

	tail, err := tail.TailFile(file, config)
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		hub:   hub,
		tail:  tail,
		cache: []string{},
	}
	return watcher, nil
}

func (w *Watcher) watch() {
	for {
		line := <-w.tail.Lines
		if line == nil {
			continue
		}
		if line.Err != nil {
			log.Println(line.Err)
		}
		u := &Update{
			Name:  w.tail.Filename,
			Lines: []string{line.Text},
		}
		w.hub.broadcast <- u
		w.cache = append(w.cache, line.Text)
	}
}

func (w *Watcher) truncateCache() {
	start := len(w.cache) - tailOffset
	if start < 0 {
		return
	}
	w.cache = w.cache[start:]
}

func (w *Watcher) lastUpdate() *Update {
	w.truncateCache()
	return &Update{
		Name:  w.tail.Filename,
		Lines: w.cache,
	}
}
