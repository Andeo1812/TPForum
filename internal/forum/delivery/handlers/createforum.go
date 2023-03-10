package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"db-performance-project/internal/forum/delivery/models"
	"db-performance-project/internal/forum/service"
	"db-performance-project/internal/pkg"
)

type forumCreateHandler struct {
	forumService service.ForumService
}

func NewForumCreateHandler(s service.ForumService) pkg.Handler {
	return &forumCreateHandler{
		s,
	}
}

func (h *forumCreateHandler) Configure(r *mux.Router, mw *pkg.HTTPMiddleware) {
	r.HandleFunc("/api/forum/create", h.Action).Methods(http.MethodPost)
}

func (h *forumCreateHandler) Action(w http.ResponseWriter, r *http.Request) {
	request := models.NewForumCreateRequest()

	request.Bind(r)
	// err := request.Bind(r)
	// if err != nil {
	//	pkg.DefaultHandlerHTTPError(r.Context(), w, err)
	//	return
	// }

	forum, err := h.forumService.CreateForum(r.Context(), request.GetForum())
	if err != nil {
		if errors.Is(errors.Cause(err), pkg.ErrSuchForumExist) {
			response := models.NewForumCreateResponse(forum)

			pkg.Response(r.Context(), w, http.StatusConflict, response)

			return
		}

		pkg.DefaultHandlerHTTPError(r.Context(), w, err)

		return
	}

	response := models.NewForumCreateResponse(forum)

	pkg.Response(r.Context(), w, http.StatusCreated, response)
}
