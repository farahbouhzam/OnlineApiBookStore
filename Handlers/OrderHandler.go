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

type OrderHandler struct {
	OrderStore interfaces.OrderStore
}

func NewOrderHandler(OrderStore interfaces.OrderStore) *OrderHandler {
	return &OrderHandler{
		OrderStore: OrderStore,
	}
}

// /orders
func (h *OrderHandler) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllOrders(w, r)
	case http.MethodPost:
		h.createOrder(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// /orders/{id}
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
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.OrderStore.GetOrder(ctx, id)
	if err != nil {
		log.Printf("ERROR fetching order %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "order not found")
		return
	}

	resp, err := json.Marshal(order)
	if err != nil {
		log.Printf("ERROR serializing order %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize order")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *OrderHandler) updateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading update order body: %v", err)
		WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var order models.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		log.Printf("ERROR unmarshalling order %d: %v", id, err)
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	updatedOrder, err := h.OrderStore.UpdateOrderStatus(ctx, id, order.Status)
	if err != nil {
		log.Printf("ERROR updating order %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "order not found")
		return
	}

	// significant business event
	log.Printf("ORDER UPDATED id=%d status=%s", id, order.Status)

	resp, err := json.Marshal(updatedOrder)
	if err != nil {
		log.Printf("ERROR serializing updated order %d: %v", id, err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize order")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading create order body: %v", err)
		WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	var order models.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		log.Printf("ERROR unmarshalling order: %v", err)
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	createdOrder, err := h.OrderStore.CreateOrder(ctx, order)
	if err != nil {
		log.Printf("ERROR creating order: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	//  significant business event
	log.Printf(
		"ORDER CREATED id=%d customer=%d total=%.2f",
		createdOrder.ID,
		createdOrder.Customer.ID,
		createdOrder.TotalPrice,
	)

	resp, err := json.Marshal(createdOrder)
	if err != nil {
		log.Printf("ERROR serializing created order: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize order")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *OrderHandler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	err = h.OrderStore.DeleteOrder(ctx, id)
	if err != nil {
		log.Printf("ERROR deleting order %d: %v", id, err)
		WriteError(w, http.StatusNotFound, "order not found")
		return
	}

	// significant business event
	log.Printf("ORDER DELETED id=%d", id)

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	orders, err := h.OrderStore.GetAllOrders(ctx)
	if err != nil {
		log.Printf("ERROR fetching orders: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to fetch orders")
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {
		log.Printf("ERROR serializing orders: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to serialize orders")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
