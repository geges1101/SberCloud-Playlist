package playlist

import "time"

type Playlist interface {
	Play() error
	Pause() error
	AddSong(song *Song) error
	Next() error
	Prev() error
}

type Song struct {
	Name     string
	Duration time.Duration
}

type playlist struct {
	head   *node
	tail   *node
	cur    *node
	length int
}

type node struct {
	song *Song
	prev *node
	next *node
}

func CreatePlaylist() Playlist {
	return &playlist{}
}

func (p *playlist) Play() error {
	if p.cur == nil {
		p.cur = p.head
	}
	go func() {
		time.Sleep(p.cur.song.Duration)
		if p.cur.next != nil {
			p.cur = p.cur.next
			p.Play()
		} else {
			p.cur = nil
		}
	}()
	return nil
}

func (p *playlist) Pause() error {
	// ничего не делать если песня не играет
	if p.cur == nil {
		return nil
	}
	p.cur = nil
	return nil
}

func (p *playlist) AddSong(song *Song) error {
	n := &node{
		song: song,
		prev: p.tail,
	}
	if p.head == nil {
		p.head = n
	} else {
		p.tail.next = n
	}
	p.tail = n
	p.length++
	return nil
}

func (p *playlist) Next() error {
	if p.cur != nil {
		p.cur = p.cur.next
	} else {
		p.cur = p.head
	}
	p.Pause()
	p.Play()
	return nil
}

func (p *playlist) Prev() error {
	if p.cur != nil {
		p.cur = p.cur.prev
	} else {
		p.cur = p.tail
	}
	p.Pause()
	p.Play()
	return nil
}
