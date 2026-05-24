package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"test_em/internal/domain"
	"test_em/internal/usecase"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	serv *usecase.Service
}

func NewHandler(serv *usecase.Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func ParseMonthYear(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse("01-2006", dateStr)
}

// @Summary      Create a subscription
// @Tags         subscriptions
// @Accept       json
// @Param        input body CreateSubscriptionInput true "Subscription Data"
// @Success      201 {string} string "Subscription successfully created!"
// @Router       /subscription/create [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input CreateSubscriptionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON body: " + err.Error()))
		return
	}

	userUUID, err := uuid.Parse(input.UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user id: " + err.Error()))
		return
	}

	StartDateMMYY, err := ParseMonthYear(input.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid start date: " + err.Error()))
		return
	}

	EndDateMMYY, err := ParseMonthYear(input.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid end date: " + err.Error()))
		return
	}

	sub := domain.Subscription{
		ID:          uuid.New(),
		UserID:      userUUID,
		ServiceName: input.ServiceName,
		Price:       input.Price,
		StartDate:   StartDateMMYY,
		EndDate:     EndDateMMYY,
	}

	err = h.serv.Create(r.Context(), sub)
	if err != nil {
		log.Printf("[ERROR] Failed to create subscription for user %s: %v", userUUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create subscription: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Subscription successfully created!"))
}

// @Summary      Update a subscription
// @Tags         subscriptions
// @Accept       json
// @Param        input body UpdateSubscriptionInput true "Updated Subscription Data"
// @Success      200 {string} string "Subscription successfully updated!"
// @Router       /subscription/update [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var input UpdateSubscriptionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON body: " + err.Error()))
		return
	}

	userUUID, err := uuid.Parse(input.UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user id: " + err.Error()))
		return
	}

	UUID, err := uuid.Parse(input.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid subscription id: " + err.Error()))
		return
	}

	StartDateMMYY, err := ParseMonthYear(input.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid start date: " + err.Error()))
		return
	}

	EndDateMMYY, err := ParseMonthYear(input.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid end date: " + err.Error()))
		return
	}

	sub := domain.Subscription{
		ID:          UUID,
		UserID:      userUUID,
		ServiceName: input.ServiceName,
		Price:       input.Price,
		StartDate:   StartDateMMYY,
		EndDate:     EndDateMMYY,
	}

	err = h.serv.Update(r.Context(), sub)
	if err != nil {
		log.Printf("[ERROR] Failed to update subscription %s: %v", UUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update subscription: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Subscription successfully updated!"))
}

// @Summary      Delete a subscription
// @Tags         subscriptions
// @Param        id query string true "Subscription UUID"
// @Success      200 {string} string "Subscription deleted"
// @Router       /subscription/delete [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing id parameter"))
		return
	}

	idUUID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id: " + err.Error()))
		return
	}

	err = h.serv.Delete(r.Context(), idUUID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete subscription %s: %v", idUUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete subscription: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Subscription deleted"))
}

// @Summary      Get a subscription
// @Tags         subscriptions
// @Param        id query string true "Subscription UUID"
// @Success      200 {object} domain.Subscription
// @Router       /subscription [get]
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing id parameter"))
		return
	}

	UUID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id: " + err.Error()))
		return
	}

	subs, err := h.serv.Get(r.Context(), UUID)
	if err != nil {
		log.Printf("[ERROR] Failed to get subscription %s: %v", UUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get subscription: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subs)
}

// @Summary      List user subscriptions
// @Tags         subscriptions
// @Param        user_id query string true "User UUID"
// @Success      200 {array} domain.Subscription
// @Router       /subscription/list [get]
func (h *Handler) ListSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("user_id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing user_id parameter"))
		return
	}

	UUID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user_id: " + err.Error()))
		return
	}

	subs, err := h.serv.List(r.Context(), UUID)
	if err != nil {
		log.Printf("[ERROR] Failed to list subscriptions for user %s: %v", UUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get subscription: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subs)
}

// @Summary      Get total cost
// @Tags         subscriptions
// @Param        user_id      query string true  "User UUID"
// @Param        start_date   query string true  "Start date (MM-YYYY)"
// @Param        end_date     query string false "End date (MM-YYYY)"
// @Param        service_name query string false "Filter by service name"
// @Success      200 {object} TotalCostResponse
// @Router      /subscription/total [get]
func (h *Handler) GetTotalCostSub(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("user_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	serviceName := r.URL.Query().Get("service_name")

	userId, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: invalid user_id: " + err.Error()))
		return
	}

	if startDateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: start_date parameter is required (format: MM-YYYY)"))
		return
	}

	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error parsing start_date: " + err.Error()))
		return
	}

	var endDate time.Time

	if endDateStr == "" {
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	} else {
		endDate, err = time.Parse("01-2006", endDateStr)
		if err == nil {
			endDate = endDate.AddDate(0, 1, 0).Add(-time.Second)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error parsing end_date: " + err.Error()))
			return
		}
	}

	var res int

	if serviceName == "" {
		res, err = h.serv.GetTotalCost(r.Context(), userId, startDate, endDate, "")
	} else {
		res, err = h.serv.GetTotalCost(r.Context(), userId, startDate, endDate, serviceName)
	}

	if err != nil {
		log.Printf("[ERROR] Failed to calculate total cost for user %s: %v", userId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to calculate total cost: " + err.Error()))
		return
	}

	response := TotalCostResponse{
		Total: res,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
