package controller

import (
	"forum/models"
	"net/http"
	"net/url"
)

func (h *Handler) filteredPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorMsg(w, http.StatusMethodNotAllowed, "")
		return
	} else if r.URL.Path != "/posts-by-user" {
		h.errorMsg(w, http.StatusNotFound, "")
		return
	}

	var data models.HomePage

	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusBadRequest, "")
		return
	}

	filter := queryParams.Get("filter")

	id := r.Context().Value(keyUser)
	if id == nil {
		h.errorMsg(w, http.StatusUnauthorized, "Need to log in")
		return
	}
	user, err := h.srv.GetUserByID(id.(int))
	if err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusInternalServerError, "")
		return
	}
	data.User = user

	posts, err := h.srv.GetFilteredByUserPosts(user.ID, filter)
	if err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusInternalServerError, "")
		return
	}

	data.Posts = posts

	categories, err := h.srv.GetCategories()
	if err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusInternalServerError, "")
		return
	}
	data.Categories = categories

	if err = templates["home"].Execute(w, data); err != nil {
		h.errLog.Println(err.Error())
		h.errorMsg(w, http.StatusInternalServerError, "")
		return
	}
}
