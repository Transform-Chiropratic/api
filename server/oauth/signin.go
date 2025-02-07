package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/Transform-Chiropratic/api/data"
	"github.com/Transform-Chiropratic/api/data/presenter"
	"github.com/Transform-Chiropratic/api/server/api"
	"github.com/Transform-Chiropratic/api/server/auth"
	"github.com/upper/db/v4"
)

type signinRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (o *OAuth) Signin(w http.ResponseWriter, r *http.Request) {
	var payload signinRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Add email to log file
	api.LogEntrySetField(r, "email", payload.Email)

	// fetch user via email
	user, err := data.DB.User.FindByEmail(payload.Email)
	if err != nil {
		if err == db.ErrNoMoreRows {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !auth.VerifyPassword(user.PasswordHash, payload.Password) {
		http.Error(w, api.ErrPermissionDenied.Error(), http.StatusUnauthorized)
		return
	}
	presented := presenter.NewAuthUser(r.Context(), user)
	api.Render(w, r, presented)
}
