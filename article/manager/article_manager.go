package manager

import (
	"context"
	"errors"
	"github.com/pi-sin/go-repo-structure/domain"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type articleManager struct {
	articleRepo    domain.ArticleRepository
	authorRepo     domain.AuthorRepository
	contextTimeout time.Duration
}

// NewArticleManager will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewArticleManager(a domain.ArticleRepository, ar domain.AuthorRepository, timeout time.Duration) domain.ArticleManager {
	return &articleManager{
		articleRepo:    a,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *articleManager) fillAuthorDetails(c context.Context, data []domain.Article) ([]domain.Article, error) {
	g, ctx := errgroup.WithContext(c)

	// Get the author's id
	mapAuthors := map[int64]domain.Author{}

	for _, article := range data {
		mapAuthors[article.Author.ID] = domain.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan domain.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != (domain.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *articleManager) Fetch(c context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *articleManager) GetByID(c context.Context, id int64) (res domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}
	res.Author = resAuthor
	return
}

func (a *articleManager) Update(c context.Context, ar *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleManager) GetByTitle(c context.Context, title string) (res domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *articleManager) Store(c context.Context, m *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, _ := a.GetByTitle(ctx, m.Title)
	if existedArticle != (domain.Article{}) {
		return errors.New("Your Item already exist")
	}

	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *articleManager) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (domain.Article{}) {
		return errors.New("Your requested Item is not found")
	}
	return a.articleRepo.Delete(ctx, id)
}
