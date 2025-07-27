package main

import (
	"net/http"

	"github.com/uday510/go-crud-app/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	// Override defaults with query values if present
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	// Validate after parsing
	if err := Validate.Struct(fq); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	// Get user feed
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(905), fq)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	// Respond with feed
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}
