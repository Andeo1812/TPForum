package models

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"db-performance-project/internal/models"
)

//go:generate easyjson -all -disallow_unknown_fields updatedetailsthread.go

type ThreadUpdateDetailsRequest struct {
	SlugOrID string
	Title    string `json:"title"`
	Message  string `json:"message"`
}

func NewThreadUpdateDetailsRequest() *ThreadUpdateDetailsRequest {
	return &ThreadUpdateDetailsRequest{}
}

func (req *ThreadUpdateDetailsRequest) Bind(r *http.Request) error {
	// if r.Header.Get("Content-Type") == "" {
	//	return pkg.ErrContentTypeUndefined
	// }
	//
	// if r.Header.Get("Content-Type") != pkg.ContentTypeJSON {
	//	return pkg.ErrUnsupportedMediaType
	// }

	vars := mux.Vars(r)

	req.SlugOrID = vars["slug_or_id"]

	body, _ := io.ReadAll(r.Body)
	// if err != nil {
	//	return pkg.ErrBadBodyRequest
	// }
	// defer func() {
	//	err = r.Body.Close()
	//	if err != nil {
	//		logrus.Error(err)
	//	}
	// }()

	// if len(body) == 0 {
	//	return pkg.ErrEmptyBody
	// }

	easyjson.Unmarshal(body, req)
	// err = easyjson.Unmarshal(body, req)
	// if err != nil {
	//	return pkg.ErrJSONUnexpectedEnd
	// }

	return nil
}

func (req *ThreadUpdateDetailsRequest) GetThread() *models.Thread {
	id, err := strconv.Atoi(req.SlugOrID)
	if err != nil {
		return &models.Thread{
			ID: uint32(id),
		}
	}

	return &models.Thread{
		Slug: req.SlugOrID,
	}
}

//easyjson:json
type ThreadUpdateDetailsResponse struct {
	ID      uint32 `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Slug    string `json:"slug"`
	Message string `json:"message"`
	Created string `json:"created"`
	Votes   int32  `json:"votes"`
}

func NewThreadUpdateDetailsResponse(thread *models.Thread) *ThreadUpdateDetailsResponse {
	return &ThreadUpdateDetailsResponse{
		ID:      thread.ID,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Slug:    thread.Slug,
		Message: thread.Message,
		Created: thread.Created,
		Votes:   thread.Votes,
	}
}