package user

import (
	"errors"
	"net/http"

	"github.com/Transform-Chiropratic/api/data"
	"github.com/Transform-Chiropratic/api/data/presenter"
	"github.com/Transform-Chiropratic/api/server/api"
	"github.com/go-chi/render"
)

var (
	maxBioLen      = 250
	maxHeadlineLen = 100

	ErrBioLen      = errors.New("bio length too long")
	ErrHeadlineLen = errors.New("headline length too long")
	ErrMissingName = errors.New("name cannot be empty")
)

func GetSessionUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value(api.SessionUserCtxKey).(*data.User)
	api.IgnoreError(render.Render(w, r, presenter.NewUser(ctx, user)))
}

type updateSessionUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
	Location  string `json:"location"`
}

func (p *updateSessionUserRequest) Bind(r *http.Request) error {
	return nil
}

func UpdateSessionUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value(api.SessionUserCtxKey).(*data.User)

	var payload updateSessionUserRequest
	if err := render.Bind(r, &payload); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInvalidRequest(err)))
		return
	}
	if payload.Location != "" {
		user.Location = payload.Location
	}
	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}

	if err := data.DB.Save(user); err != nil {
		api.IgnoreError(render.Render(w, r, api.ErrInternalServerError(err)))
		return
	}

	presented := presenter.NewUser(ctx, user)
	api.Render(w, r, presented)
}
