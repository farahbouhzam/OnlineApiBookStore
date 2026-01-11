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

type AuthorHandler struct {
	AuthorStore interfaces.AuthorStore
}

func NewAuthorHandler(AuthorStore interfaces.AuthorStore) *AuthorHandler {
	return &AuthorHandler{
		AuthorStore: AuthorStore,
	}
}

// /books
func (h *AuthorHandler) AuthorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAuthors(w, r)
	case http.MethodPost:
		h.createAuthor(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *AuthorHandler) getAuthors(w http.ResponseWriter, r *http.Request)

// /books/{id}
func (h *AuthorHandler) AuthorsByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAuthorByID(w, r)
	case http.MethodPut:
		h.updateAuthor(w, r)
	case http.MethodDelete:
		h.deleteAuthor(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *AuthorHandler) getAuthorByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/authors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	author, err := h.AuthorStore.GetAuthor(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *AuthorHandler) updateAuthor(w http.ResponseWriter, r *http.Request) {
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

	var author models.Author
	err = json.Unmarshal(body, &author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedBook, err := h.AuthorStore.UpdateAuthor(id,author)
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

func (h *AuthorHandler) createAuthor(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var author models.Author
	err = json.Unmarshal(body, &author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdAuthor, err := h.AuthorStore.CreateAuthor(author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(createdAuthor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *AuthorHandler) deleteAuthor(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/authors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.AuthorStore.DeleteAuthor(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
