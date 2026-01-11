package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"online_bookStore/interfaces"
	"online_bookStore/models"
)

type BookHandler struct {
	bookStore interfaces.BookStore
}

func NewBookHandler(bookStore interfaces.BookStore) *BookHandler {
	return &BookHandler{
		bookStore: bookStore,
	}
}

// /books
func (h *BookHandler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getBooks(w, r)
	case http.MethodPost:
		h.createBook(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request)

// /books/{id}
func (h *BookHandler) BookByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getBookByID(w, r)
	case http.MethodPut:
		h.updateBook(w, r)
	case http.MethodDelete:
		h.deleteBook(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) getBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.bookStore.GetBook(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedBook, err := h.bookStore.UpdateBook(id, book)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(updatedBook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdBook, err := h.bookStore.CreateBook(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(createdBook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.bookStore.DeleteBook(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
