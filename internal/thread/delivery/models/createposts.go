package models

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"db-performance-project/internal/models"
)

//go:generate easyjson -disallow_unknown_fields createposts.go

//easyjson:json
type PostRequest struct {
	Parent  uint32 `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

//easyjson:json
type PostsRequestList []PostRequest

type ThreadCreatePostsRequest struct {
	SlugOrID string
	Posts    PostsRequestList
}

func NewThreadCreatePostsRequest() *ThreadCreatePostsRequest {
	return &ThreadCreatePostsRequest{}
}

func (req *ThreadCreatePostsRequest) Bind(r *http.Request) error {
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

	easyjson.Unmarshal(body, &req.Posts)
	// err = easyjson.Unmarshal(body, req)
	// if err != nil {
	//	return pkg.ErrJSONUnexpectedEnd
	// }

	return nil
}

func (req *ThreadCreatePostsRequest) GetPosts() []*models.Post {
	res := make([]*models.Post, len(req.Posts))

	for idx, value := range req.Posts {
		res[idx] = &models.Post{
			Parent:  value.Parent,
			Message: value.Message,
			Author: models.User{
				Nickname: value.Author,
			},
		}
	}

	return res
}

func (req *ThreadCreatePostsRequest) GetThread() *models.Thread {
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
type PostResponse struct {
	ID       uint32 `json:"id"`
	Parent   uint32 `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   uint32 `json:"thread"`
	Created  string `json:"created"`
}

//easyjson:json
type PostsResponseList []PostResponse

func NewThreadCreatePostsResponse(posts []*models.Post) PostsResponseList {
	res := make([]PostResponse, len(posts))

	for idx, value := range posts {
		res[idx] = PostResponse{
			ID:       value.ID,
			Parent:   value.Parent,
			Author:   value.Author.Nickname,
			Forum:    value.Forum,
			IsEdited: value.IsEdited,
			Message:  value.Message,
			Created:  value.Created,
			Thread:   value.Thread,
		}
	}

	return res
}