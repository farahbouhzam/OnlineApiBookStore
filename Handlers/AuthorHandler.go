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

type AuthorHandler struct {
	AuthorStore interfaces.AuthorStore
}

func NewAuthorHandler(authorStore interfaces.AuthorStore) *AuthorHandler {
	return &AuthorHandler{
		AuthorStore: authorStore,
	}
}

/*
	ROUTE: /authors
*/
func (h *AuthorHandler) AuthorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAuthors(w, r)
	case http.MethodPost:
		h.createAuthor(w, r)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

/*
	GET /authors
*/
func (h *AuthorHandler) getAuthors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	authors, err := h.AuthorStore.GetAllAuthors(ctx)
	if err != nil {
		log.Printf("ERROR fetching authors: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to fetch authors")
		return
	}

	resp, err := json.Marshal(authors)
	if err != nil {
		log.Printf("ERROR serializing authors: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize authors")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

/*
	ROUTE: /authors/{id}
*/
func (h *AuthorHandler) AuthorsByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAuthorByID(w, r)
	case http.MethodPut:
		h.updateAuthor(w, r)
	case http.MethodDelete:
		h.deleteAuthor(w, r)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

/*
	GET /authors/{id}
*/
func (h *AuthorHandler) getAuthorByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	id, err := parseID(r.URL.Path, "/authors/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid author id")
		return
	}

	author, err := h.AuthorStore.GetAuthor(ctx, id)
	if err != nil {
		log.Printf("ERROR fetching author %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "author not found")
		return
	}

	resp, err := json.Marshal(author)
	if err != nil {
		log.Printf("ERROR serializing author %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize author")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

/*
	PUT /authors/{id}
*/
func (h *AuthorHandler) updateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	id, err := parseID(r.URL.Path, "/authors/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid author id")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading update author body: %v", err)
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var author models.Author
	if err := json.Unmarshal(body, &author); err != nil {
		log.Printf("ERROR unmarshalling author %d: %v", id, err)
		WriteError(w, http.StatusBadRequest, "invalid author payload")
		return
	}

	updatedAuthor, err := h.AuthorStore.UpdateAuthor(ctx, id, author)
	if err != nil {
		log.Printf("ERROR updating author %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "author not found")
		return
	}

	resp, err := json.Marshal(updatedAuthor)
	if err != nil {
		log.Printf("ERROR serializing updated author %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize author")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

/*
	POST /authors
*/
func (h *AuthorHandler) createAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading create author body: %v", err)
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var author models.Author
	if err := json.Unmarshal(body, &author); err != nil {
		log.Printf("ERROR unmarshalling author: %v", err)
		WriteError(w, http.StatusBadRequest, "invalid author payload")
		return
	}

	createdAuthor, err := h.AuthorStore.CreateAuthor(ctx, author)
	if err != nil {
		log.Printf("ERROR creating author: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to create author")
		return
	}

	//  significant business log
	log.Printf(
		"AUTHOR CREATED id=%d name=%s %s",
		createdAuthor.ID,
		createdAuthor.FirstName,
		createdAuthor.LastName,
	)

	resp, err := json.Marshal(createdAuthor)
	if err != nil {
		log.Printf("ERROR serializing created author: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize author")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

/*
	DELETE /authors/{id}
*/
func (h *AuthorHandler) deleteAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	id, err := parseID(r.URL.Path, "/authors/")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid author id")
		return
	}

	if err := h.AuthorStore.DeleteAuthor(ctx, id); err != nil {
		log.Printf("ERROR deleting author %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "author not found")
		return
	}

	//  significant business log
	log.Printf("AUTHOR DELETED id=%d", id)

	w.WriteHeader(http.StatusNoContent)
}

/*
	HELPER: parse ID from URL
*/
func parseID(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	return strconv.Atoi(idStr)
}
