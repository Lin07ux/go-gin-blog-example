package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/middleware/jwt"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	v1 "github.com/lin07ux/go-gin-example/routers/v1"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.RunMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")
	apiV1.Use(jwt.JWT())
	{
		// 文章标签接口
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.EditTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		// 文章接口
		apiV1.GET("/articles", v1.GetArticles)
		apiV1.POST("/articles", v1.AddArticle)
		apiV1.GET("/articles/:id", v1.GetArticle)
		apiV1.PUT("/articles/:id", v1.EditArticle)
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
