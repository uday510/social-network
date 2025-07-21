package main

import (
	"github.com/uday510/go-crud-app/internal/store"
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	// Override defaults with query values if present
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate after parsing
	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Get user feed
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(905), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Respond with feed
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
