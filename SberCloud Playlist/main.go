package main

import (
	"SberCloud_Playlist/model"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

func getPlaylist(w http.ResponseWriter, r *http.Request) {
	playlist.Mutex.Lock()
	defer playlist.Mutex.Unlock()

	jsonBytes, err := json.Marshal(playlist)
	if err != nil {
		http.Error(w, "Failed to encode model as JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func deleteSong(w http.ResponseWriter, r *http.Request) {
	playlist.Mutex.Lock()
	defer playlist.Mutex.Unlock()

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID parameter is missing", http.StatusBadRequest)
		return
	}

	// Находим песню по id
	var nodeToDelete *model.Node
	for node := playlist.Head; node != nil; node = node.Next {
		if node.Song.ID == id {
			nodeToDelete = node
			break
		}
	}

	if nodeToDelete == nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Удаляем элемент из плейлиста
	if nodeToDelete == playlist.Curr {
		playlist.Curr = nodeToDelete.Next
	}
	if nodeToDelete == playlist.Head {
		playlist.Head = nodeToDelete.Next
		if playlist.Head != nil {
			playlist.Head.Prev = nil
		}
	}
	if nodeToDelete == playlist.Tail {
		playlist.Tail = nodeToDelete.Prev
		if playlist.Tail != nil {
			playlist.Tail.Next = nil
		}
	}
	if nodeToDelete.Prev != nil {
		nodeToDelete.Prev.Next = nodeToDelete.Next
	}
	if nodeToDelete.Next != nil {
		nodeToDelete.Next.Prev = nodeToDelete.Prev
	}

	playlist.Length--
	w.WriteHeader(http.StatusOK)
}

func getSong(w http.ResponseWriter, r *http.Request) {
	playlist.Mutex.Lock()
	defer playlist.Mutex.Unlock()

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID parameter is missing", http.StatusBadRequest)
		return
	}

	// Находим песню по id и возвращаем
	var song *model.Song
	for node := playlist.Head; node != nil; node = node.Next {
		if node.Song.ID == id {
			song = node.Song
			break
		}
	}

	if song == nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Marshal the song into JSON bytes
	jsonBytes, err := json.Marshal(song)
	if err != nil {
		http.Error(w, "Failed to encode song as JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func addSong(w http.ResponseWriter, r *http.Request) {
	playlist.Mutex.Lock()
	defer playlist.Mutex.Unlock()

	// декодируем песню
	var newSong model.Song
	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// проверяем используется ли уже данный id песни
	for node := playlist.Head; node != nil; node = node.Next {
		if node.Song.ID == newSong.ID {
			http.Error(w, "Song ID already in use", http.StatusConflict)
			return
		}
	}

	// Создаем узел для новой песни
	newNode := &model.Node{
		Song: &newSong,
		Next: nil,
		Prev: playlist.Tail,
	}

	// Добавляем узел
	if playlist.Head == nil {
		playlist.Head = newNode
	}
	if playlist.Tail != nil {
		playlist.Tail.Next = newNode
	}
	playlist.Tail = newNode
	playlist.Length++

	// Возвращаем json добавленной песни
	jsonBytes, err := json.Marshal(newSong)
	if err != nil {
		http.Error(w, "Failed to encode song as JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func updateSong(w http.ResponseWriter, r *http.Request) {
	var updatedSong model.Song
	err := json.NewDecoder(r.Body).Decode(&updatedSong)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// достаем id песни
	id := chi.URLParam(r, "id")

	playlist.Mutex.Lock()
	defer playlist.Mutex.Unlock()

	node := playlist.Head
	for node != nil {
		if node.Song.ID == id {
			// Обновляем данные песни
			node.Song.Title = updatedSong.Title
			node.Song.Artist = updatedSong.Artist
			node.Song.Duration = updatedSong.Duration

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(node.Song)
			return
		}
		node = node.Next
	}

	// Если песня не найдена
	http.NotFound(w, r)
}

var playlist = model.Playlist{
	Head:    nil,
	Tail:    nil,
	Curr:    nil,
	Length:  0,
	Playing: false,
	Mutex:   sync.Mutex{},
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/model", getPlaylist).Methods("GET")
	r.HandleFunc("/model/{id}", getSong).Methods("GET")
	r.HandleFunc("/model", addSong).Methods("POST")
	r.HandleFunc("/model/{id}", updateSong).Methods("PUT")
	r.HandleFunc("/model/{id}", deleteSong).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
