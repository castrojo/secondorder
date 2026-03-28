package handlers

import (
	"fmt"
	"net/http"
	"sync"
)

type SSEHub struct {
	mu       sync.RWMutex
	clients  map[chan string]struct{}
	shutdown chan struct{}
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients:  make(map[chan string]struct{}),
		shutdown: make(chan struct{}),
	}
}

func (h *SSEHub) Close() {
	close(h.shutdown)
	h.mu.Lock()
	for ch := range h.clients {
		close(ch)
	}
	h.clients = make(map[chan string]struct{})
	h.mu.Unlock()
}

func (h *SSEHub) Broadcast(event, data string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	msg := fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
	for ch := range h.clients {
		select {
		case ch <- msg:
		default:
			// drop if client is slow
		}
	}
}

func (h *SSEHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := make(chan string, 16)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, ch)
		h.mu.Unlock()
	}()

	// Send initial keepalive
	fmt.Fprintf(w, ": keepalive\n\n")
	flusher.Flush()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprint(w, msg)
			flusher.Flush()
		case <-h.shutdown:
			return
		case <-r.Context().Done():
			return
		}
	}
}
