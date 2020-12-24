package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pi-sin/go-repo-structure/domain"
	"net/http"
	"strconv"

	validator "gopkg.in/go-playground/validator.v9"
)

// ArticleHandler  represent the http handler for article
type ArticleHandler struct {
	AManager domain.ArticleManager
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewHttpArticleHandler(e *gin.Engine, am domain.ArticleManager) {
	handler := &ArticleHandler{
		AManager: am,
	}
	e.GET("/articles", handler.FetchArticle)
	e.POST("/articles", handler.Store)
	e.GET("/articles/:id", handler.GetByID)
	e.DELETE("/articles/:id", handler.Delete)
}

// FetchArticle will fetch the article based on given params
func (a *ArticleHandler) FetchArticle(c *gin.Context) {
	numS := c.Param("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.Param("cursor")
	ctx := c.Request.Context()

	listAr, nextCursor, err := a.AManager.Fetch(ctx, cursor, int64(num))
	if err != nil {
		//log error if required
		c.JSON(http.StatusBadRequest, domain.Response{
			Status:  domain.Error,
			Message: domain.ErrorMessageBadRequest,
		})
	}

	c.Writer.Header().Set("X-Cursor", nextCursor)
	c.JSON(http.StatusOK, domain.Response{
		Status:  domain.Success,
		Message: listAr,
	})
}

// GetByID will get article by given id
func (a *ArticleHandler) GetByID(c *gin.Context) {
	idP, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		//log error if required
		c.JSON(http.StatusBadRequest, domain.Response{
			Status:  domain.Error,
			Message: domain.ErrorMessageBadRequest,
		})
	}

	id := int64(idP)
	ctx := c.Request.Context()

	art, err := a.AManager.GetByID(ctx, id)
	if err != nil {
		//log error if required
		c.JSON(http.StatusNotFound, domain.Response{
			Status:  domain.Error,
			Message: domain.ErrorMessageArticleNotFound,
		})
	}

	c.JSON(http.StatusOK, domain.Response{
		Status:  domain.Success,
		Message: art,
	})
}

func isRequestValid(m *domain.Article) bool {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		//log error if required
		return false
	}
	return true
}

// Store will store the article by given request body
func (a *ArticleHandler) Store(c *gin.Context) {
	var article domain.Article
	if err := c.BindJSON(&article); err != nil {
		//log error if required
		c.JSON(http.StatusUnprocessableEntity, domain.Response{
			Status:  domain.Error,
			Message: domain.ErrorMessageBadRequest,
		})
	}

	if ok := isRequestValid(&article); !ok {
		c.JSON(http.StatusBadRequest, domain.Response{
			Status:  domain.Error,
			Message: domain.ErrorMessageBadRequest,
		})
	}

	ctx := c.Request.Context()
	if err := a.AManager.Store(ctx, &article); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status:  domain.Success,
			Message: domain.ErrorMessageArticleNotSaved,
		})
	}

	c.JSON(http.StatusCreated, domain.Response{
		Status:  domain.Success,
		Message: article,
	})
}

// Delete will delete article by given param
func (a *ArticleHandler) Delete(c *gin.Context) {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Status:  domain.Error,
			Message: domain.ErrorMessageBadRequest,
		})
	}

	id := int64(idP)
	ctx := c.Request.Context()

	if err = a.AManager.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status:  domain.Success,
			Message: domain.ErrorMessageArticleNotDeleted,
		})
	}

	c.Status(http.StatusNoContent)
}
