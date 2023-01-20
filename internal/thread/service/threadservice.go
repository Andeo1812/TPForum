package service

import (
	"context"

	"github.com/pkg/errors"

	repoForum "db-performance-project/internal/forum/repository"
	"db-performance-project/internal/models"
	"db-performance-project/internal/pkg"
	repoPost "db-performance-project/internal/post/repository"
	repoThread "db-performance-project/internal/thread/repository"
	repoUser "db-performance-project/internal/user/repository"
)

type ThreadService interface {
	CreateThread(ctx context.Context, thread *models.Thread) (models.Thread, error)
	CreatePosts(ctx context.Context, thread *models.Thread, posts []*models.Post) ([]models.Post, error)
	GetDetailsThread(ctx context.Context, thread *models.Thread) (models.Thread, error)
	GetPosts(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error)
	UpdateThread(ctx context.Context, thread *models.Thread) (models.Thread, error)
}

type threadService struct {
	threadRepo repoThread.ThreadRepository
	forumRepo  repoForum.ForumRepository
	userRepo   repoUser.UserRepository
	postRepo   repoPost.PostRepository
}

func NewThreadService(rt repoThread.ThreadRepository, rf repoForum.ForumRepository, ru repoUser.UserRepository, rp repoPost.PostRepository) ThreadService {
	return &threadService{
		threadRepo: rt,
		forumRepo:  rf,
		userRepo:   ru,
		postRepo:   rp,
	}
}

func (t threadService) CreateThread(ctx context.Context, thread *models.Thread) (models.Thread, error) {
	// CheckAuthor
	resUser, err := t.userRepo.GetUserByNickname(ctx, &models.User{Nickname: thread.Author})
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "CreateForum")
	}
	thread.Author = resUser.Nickname

	// CheckForum
	resForum, err := t.forumRepo.GetDetailsForumBySlug(ctx, &models.Forum{Slug: thread.Forum})
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "CreateForum")
	}
	thread.Forum = resForum.Slug

	// CheckThread
	if thread.Slug != "" {
		var existThread models.Thread

		existThread, err = t.threadRepo.GetDetailsThreadBySlug(ctx, thread)
		if err == nil {
			return existThread, errors.Wrap(pkg.ErrSuchThreadExist, "CreateForum")
		}
	}

	// All valid - action
	res, err := t.threadRepo.CreateThread(ctx, thread)
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "CreateThread")
	}

	return res, err
}

func (t threadService) CreatePosts(ctx context.Context, thread *models.Thread, posts []*models.Post) ([]models.Post, error) {
	var err error

	threadSum := models.Thread{ID: thread.ID, Slug: thread.Slug}

	// Получаем ID по слагу
	if thread.Slug != "" {
		var threadID models.Thread

		threadID, err = t.threadRepo.GetThreadIDBySlug(ctx, thread)
		if err != nil {
			return []models.Post{}, errors.Wrap(err, "CreatePosts")
		}

		threadSum.ID = threadID.ID
	}

	var threadForum models.Thread

	threadForum, err = t.threadRepo.GetThreadForumByID(ctx, &threadSum)
	if err != nil {
		return []models.Post{}, errors.Wrap(err, "CreatePosts")
	}

	threadSum.Forum = threadForum.Forum

	// if posts[0].Parent != 0 {
	//	var postWithParent *models.Post
	//
	//	postWithParent, err = t.postRepo.GetParentPost(ctx, posts[0])
	//	if err != nil {
	//		return []models.Post{}, errors.Wrap(err, "CreatePosts")
	//	}
	//
	//	if postWithParent.Parent != thread.ID {
	//		return nil, errors.Wrap(pkg.ErrInvalidParent, "CreatePosts")
	//	}
	// }

	res, err := t.threadRepo.CreatePostsByID(ctx, &threadSum, posts)
	if err != nil {
		return []models.Post{}, errors.Wrap(err, "CreatePosts")
	}

	return res, nil
}

func (t threadService) GetDetailsThread(ctx context.Context, thread *models.Thread) (models.Thread, error) {
	var err error

	var resThread models.Thread

	// CheckAndGetThread
	if thread.Slug != "" {
		resThread, err = t.threadRepo.GetDetailsThreadBySlug(ctx, thread)
	} else {
		resThread, err = t.threadRepo.GetDetailsThreadByID(ctx, thread)
	}
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "GetDetailsThread")
	}

	return resThread, nil
}

func (t threadService) UpdateThread(ctx context.Context, thread *models.Thread) (models.Thread, error) {
	var err error

	resThread := models.Thread{}

	// CheckAndGetThread
	if thread.Slug != "" {
		resThread, err = t.threadRepo.GetDetailsThreadBySlug(ctx, thread)
	} else {
		resThread, err = t.threadRepo.GetDetailsThreadByID(ctx, thread)
	}
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "UpdateThread")
	}

	resThread.Title = thread.Title
	resThread.Message = thread.Message

	res, err := t.threadRepo.UpdateThreadByID(ctx, &resThread)
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "UpdateThread")
	}

	return res, nil
}

func (t threadService) GetPosts(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error) {
	var res []models.Post
	var err error

	resThread := models.Thread{}

	// CheckAndGetThread
	if thread.Slug != "" {
		resThread, err = t.threadRepo.GetDetailsThreadBySlug(ctx, thread)
	} else {
		resThread, err = t.threadRepo.GetDetailsThreadByID(ctx, thread)
	}
	if err != nil {
		return []models.Post{}, errors.Wrap(err, "GetPosts")
	}

	switch params.Sort {
	case pkg.TypeSortFlat:
		res, err = t.threadRepo.GetPostsByIDFlat(ctx, &resThread, params)
	case pkg.TypeSortTree:
		res, err = t.threadRepo.GetPostsByIDTree(ctx, &resThread, params)
	case pkg.TypeSortParentTree:
		res, err = t.threadRepo.GetPostsByIDParentTree(ctx, &resThread, params)
	default:
		return nil, errors.Wrap(pkg.ErrNoSuchRuleSortPosts, "GetPosts")
	}
	if err != nil {
		return nil, errors.Wrap(err, "GetPosts")
	}

	return res, nil
}
