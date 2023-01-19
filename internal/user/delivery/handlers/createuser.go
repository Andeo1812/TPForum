package handlers

import (
	"github.com/pkg/errors"
	"net/http"

	"db-performance-project/internal/pkg"
	"db-performance-project/internal/user/delivery/models"
	"db-performance-project/internal/user/service"
	"github.com/gorilla/mux"
)

type userCreateHandler struct {
	userService service.UserService
}

func NewUserCreateHandler(s service.UserService) pkg.Handler {
	return &userCreateHandler{
		s,
	}
}

func (h *userCreateHandler) Configure(r *mux.Router, mw *pkg.HTTPMiddleware) {
	r.HandleFunc("/user/{nickname}/create", h.Action).Methods(http.MethodPost)
}

func (h *userCreateHandler) Action(w http.ResponseWriter, r *http.Request) {
	request := models.NewUserCreateRequest()

	request.Bind(r)
	// err := request.Bind(r)
	// if err != nil {
	//	pkg.DefaultHandlerHTTPError(r.Context(), w, err)
	//	return
	// }

	users, err := h.userService.CreateUser(r.Context(), request.GetUser())
	if err != nil {
		if errors.Is(errors.Cause(err), pkg.ErrSuchUserExist) {
			response := models.NewUsersCreateResponse(users)

			pkg.Response(r.Context(), w, http.StatusConflict, response)

			return
		}

		pkg.DefaultHandlerHTTPError(r.Context(), w, err)

		return
	}

	response := models.NewUserCreateResponse(users[0])

	pkg.Response(r.Context(), w, http.StatusCreated, response)
}
