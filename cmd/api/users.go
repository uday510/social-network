package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/uday510/go-crud-app/internal/store"
	"net/http"
	"strconv"
)

type userKy string

const userCtx postKey = "user"

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	ctx := r.Context()
	var payload FollowUser

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Revert back to auth userID from ctx
	app.store.Users.Follow(ctx, followerUser.ID)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		postId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Users.GetByID(ctx, postId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.badRequestError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
