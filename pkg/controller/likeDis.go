package controller

import (
	"encoding/json"
	"forum/models"
	"net/http"
)

func (h *Handler) ratePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorMsg(w, http.StatusMethodNotAllowed, "")
		return
	} else if r.URL.Path != "/post/rate" {
		h.errorMsg(w, http.StatusNotFound, "")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var rate models.RatePost

	err := decoder.Decode(&rate)
	if err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusBadRequest, "")
		return
	}

	user_id := r.Context().Value(keyUser).(int)
	rate.User_ID = user_id

	if err := h.srv.Post.RatePost(rate); err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusBadRequest, "")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) rateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorMsg(w, http.StatusMethodNotAllowed, "")
		return
	} else if r.URL.Path != "/comment/rate" {
		h.errorMsg(w, http.StatusNotFound, "")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var rate models.RateComment

	err := decoder.Decode(&rate)
	if err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusBadRequest, "")
		return
	}

	user_id := r.Context().Value(keyUser).(int)
	rate.User_ID = user_id

	if err := h.srv.Comment.RateComment(rate); err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusBadRequest, "")
		return
	}
	w.WriteHeader(http.StatusOK)
}
