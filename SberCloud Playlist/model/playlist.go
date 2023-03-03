package model

import (
	"sync"
	"time"
)

// Playlist структура плейлиста
type Playlist struct {
	Head    *Node
	Tail    *Node
	Curr    *Node
	Length  int
	Playing bool
	Mutex   sync.Mutex
}

// Node узел в двусвязном списке
type Node struct {
	Song *Song
	Next *Node
	Prev *Node
}

// Song песня в плейлисте
type Song struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Artist   string        `json:"artist"`
	Duration time.Duration `json:"duration"`
}

// CreatePlaylist создает новый плейлист
func CreatePlaylist() *Playlist {
	return &Playlist{}
}

// Play возобновляет проигрывание плейлиста
func (p *Playlist) Play() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if !p.Playing {
		p.Playing = true
		return
	}

	if p.Curr == nil {
		p.Curr = p.Head
	} else {
		p.Curr = p.Curr.Next
	}

	if p.Curr != nil {
		go func() {
			time.Sleep(p.Curr.Song.Duration)
			p.Play()
		}()
	}
}

// Pause останавливает проигрывание плейлиста
func (p *Playlist) Pause() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.Playing = false
}

// AddSong добавляет песню в конец плейлиста
func (p *Playlist) AddSong(song *Song) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	node := &Node{
		Song: song,
	}

	if p.Head == nil {
		p.Head = node
		p.Tail = node
	} else {
		p.Tail.Next = node
		node.Prev = p.Tail
		p.Tail = node
	}
	p.Length++
}

// Next переключается на следующую песню в плейлисте
func (p *Playlist) Next() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if p.Curr != nil {
		p.Curr = p.Curr.Next
		p.Playing = true
	}

	if p.Curr != nil {
		go func() {
			p.Play()
		}()
	}
}

// Prev переключается на предыдущую песню в плейлисте
func (p *Playlist) Prev() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if p.Curr.Prev != nil {
		p.Curr = p.Curr.Prev
		p.Playing = true
		go func() {
			p.Play()
		}()
	}
}
