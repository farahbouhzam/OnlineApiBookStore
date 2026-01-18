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

type CustomerHandler struct {
	CustomerStore interfaces.CustomerStore
}

func NewCustomerHandler(CustomerStore interfaces.CustomerStore) *CustomerHandler {
	return &CustomerHandler{
		CustomerStore: CustomerStore,
	}
}

// /customers
func (h *CustomerHandler) CustomersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getCustomers(w, r)
	case http.MethodPost:
		h.createCustomer(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CustomerHandler) getCustomers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	customers, err := h.CustomerStore.GetAllCustomers(ctx)
	if err != nil {
		log.Printf("ERROR fetching customers: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to fetch customers")
		return
	}

	resp, err := json.Marshal(customers)
	if err != nil {
		log.Printf("ERROR serializing customers: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize customers")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// /customers/{id}
func (h *CustomerHandler) CustomersByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getCustomerByID(w, r)
	case http.MethodPut:
		h.updateCustomer(w, r)
	case http.MethodDelete:
		h.deleteCustomer(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CustomerHandler) getCustomerByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/customers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	customer, err := h.CustomerStore.GetCustomer(ctx, id)
	if err != nil {
		log.Printf("ERROR fetching customer %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "customer not found")
		return
	}

	resp, err := json.Marshal(customer)
	if err != nil {
		log.Printf("ERROR serializing customer %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize customer")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *CustomerHandler) updateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/customers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading update customer body: %v", err)
		WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var customer models.Customer
	err = json.Unmarshal(body, &customer)
	if err != nil {
		log.Printf("ERROR unmarshalling customer %d: %v", id, err)
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	updatedCustomer, err := h.CustomerStore.UpdateCustomer(ctx, id, customer)
	if err != nil {
		log.Printf("ERROR updating customer %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "customer not found")
		return
	}

	// significant business event
	log.Printf("CUSTOMER UPDATED id=%d", id)

	resp, err := json.Marshal(updatedCustomer)
	if err != nil {
		log.Printf("ERROR serializing updated customer %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize customer")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading create customer body: %v", err)
		WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var customer models.Customer
	err = json.Unmarshal(body, &customer)
	if err != nil {
		log.Printf("ERROR unmarshalling customer: %v", err)
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	createdCustomer, err := h.CustomerStore.CreateCustomer(ctx, customer)
	if err != nil {
		log.Printf("ERROR creating customer: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	// significant business event
	log.Printf("CUSTOMER CREATED id=%d email=%s", createdCustomer.ID, createdCustomer.Email)

	resp, err := json.Marshal(createdCustomer)
	if err != nil {
		log.Printf("ERROR serializing created customer: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize customer")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *CustomerHandler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/customers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	err = h.CustomerStore.DeleteCustomer(ctx, id)
	if err != nil {
		log.Printf("ERROR deleting customer %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "customer not found")
		return
	}

	// significant business event
	log.Printf("CUSTOMER DELETED id=%d", id)

	w.WriteHeader(http.StatusNoContent)
}
