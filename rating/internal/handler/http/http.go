package http

import (
	"encoding/json"
	"errors"
	"log"
	"movie-rating-app/rating/internal/controller/rating"
	model "movie-rating-app/rating/pkg"
	"net/http"
	"strconv"
)

type Handler struct {
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	recordId := model.RecordID(req.FormValue("id"))
	recordType := model.RecordType(req.FormValue("type"))
	if recordId == "" || recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	switch req.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(ctx, recordId, recordType)
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Error while getting aggregated rating %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("Error encoding response %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPut:
		userId := req.FormValue("userId")
		v, err := strconv.ParseFloat(req.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rating := &model.Rating{
			UserID:     model.UserID(userId),
			Value:      model.RatingValue(v),
			RecordID:   string(recordId),
			RecordType: string(recordType),
		}
		if err := h.ctrl.PutRating(ctx, recordId, recordType, rating); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
