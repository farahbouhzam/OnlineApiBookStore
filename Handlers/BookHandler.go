package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"online_bookStore/Interfaces"
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

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// read query parameters
	q := r.URL.Query()
	title := q.Get("title")
	genre := q.Get("genre")

	minPrice, _ := strconv.ParseFloat(q.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(q.Get("max_price"), 64)
	authorId, _ := strconv.Atoi(q.Get("author_id"))

	criteria := models.SearchCriteria{
		Title:    title,
		Genre:    genre,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		AuthorId: authorId,
	}

	books, err := h.bookStore.SearchBooks(ctx, criteria)
	if err != nil {
		log.Printf("ERROR fetching books: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to fetch books")
		return
	}

	resp, err := json.Marshal(books)
	if err != nil {
		log.Printf("ERROR serializing books: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize books")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading create book body: %v", err)
		WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var book models.Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Printf("ERROR unmarshalling book: %v", err)
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	createdBook, err := h.bookStore.CreateBook(ctx, book)
	if err != nil {
		log.Printf("ERROR creating book: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to create book")
		return
	}

	// significant business event
	log.Printf("BOOK CREATED id=%d title=%s", createdBook.ID, createdBook.Title)

	resp, err := json.Marshal(createdBook)
	if err != nil {
		log.Printf("ERROR serializing created book: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize book")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

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
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	book, err := h.bookStore.GetBook(ctx, id)
	if err != nil {
		log.Printf("ERROR fetching book %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "book not found")
		return
	}

	resp, err := json.Marshal(book)
	if err != nil {
		log.Printf("ERROR serializing book %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize book")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading update book body: %v", err)
		WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var book models.Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Printf("ERROR unmarshalling book %d: %v", id, err)
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	updatedBook, err := h.bookStore.UpdateBook(ctx, id, book)
	if err != nil {
		log.Printf("ERROR updating book %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "book not found")
		return
	}

	//  significant business event
	log.Printf("BOOK UPDATED id=%d", id)

	resp, err := json.Marshal(updatedBook)
	if err != nil {
		log.Printf("ERROR serializing updated book %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize book")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	err = h.bookStore.DeleteBook(ctx, id)
	if err != nil {
		log.Printf("ERROR deleting book %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "book not found")
		return
	}

	// significant business event
	log.Printf("BOOK DELETED id=%d", id)

	w.WriteHeader(http.StatusNoContent)
}
