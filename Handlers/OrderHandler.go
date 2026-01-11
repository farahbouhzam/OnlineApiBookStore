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

type OrderHandler struct {
	OrderStore interfaces.OrderStore
}

func NewOrderHandler(OrderStore interfaces.OrderStore) *OrderHandler {
	return &OrderHandler{
		OrderStore: OrderStore,
	}
}

// /books
func (h *OrderHandler) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getOrders(w, r)
	case http.MethodPost:
		h.createOrder(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) getOrders(w http.ResponseWriter, r *http.Request)

// /books/{id}
func (h *OrderHandler) OrdersByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getOrderByID(w, r)
	case http.MethodPut:
		h.updateOrder(w, r)
	case http.MethodDelete:
		h.deleteOrder(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) getOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/Orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Order, err := h.OrderStore.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *OrderHandler) updateOrder(w http.ResponseWriter, r *http.Request) {
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

	var Order models.Order
	err = json.Unmarshal(body, &Order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedBook, err := h.OrderStore.UpdateOrderStatus(id,Order.Status)
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

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var Order models.Order
	err = json.Unmarshal(body, &Order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdOrder, err := h.OrderStore.CreateOrder(Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(createdOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *OrderHandler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.OrderStore.DeleteOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
