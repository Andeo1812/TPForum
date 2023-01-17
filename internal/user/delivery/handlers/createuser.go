package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"db-performance-project/internal/pkg"
	"db-performance-project/internal/user/delivery/models"
	"db-performance-project/internal/user/service"
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

	user, err := h.userService.CreateUser(r.Context(), request.GetUser())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewUserCreateResponse(user)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}
