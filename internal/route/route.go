package route

import (
	"github.com/jianchengwang/coderunner/internal/auth"
	"github.com/jianchengwang/coderunner/internal/route/fetch"
	tmpl "html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	log "unknwon.dev/clog/v2"

	"github.com/jianchengwang/coderunner/frontend"
	task "github.com/jianchengwang/coderunner/internal/route/task"
	"github.com/jianchengwang/coderunner/public"
	"github.com/jianchengwang/coderunner/templates"
)

// New returns a new gin router.
func New() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Content-type", "User-Agent"},
		AllowCredentials: true,
		AllowOrigins:     []string{os.Getenv("APP_URL")},
	}))

	// Session
	store := cookie.NewStore([]byte(randstr.String(50)))
	r.Use(sessions.Sessions("coderunner", store))

	// Templates
	tpl := tmpl.Must(tmpl.New("").Delims("[[","]]").ParseFS(templates.FS, "*"))
	r.SetHTMLTemplate(tpl)

	r.GET("/", IndexHandler)

	run := r.Group("/r")
	{
		run.GET("/new", task.NewTaskHandler)
		run.GET("/:uid", task.EditorHandler)
		run.GET("/:uid/init", __(task.InitTaskHandler))
		run.POST("/:uid",  __(task.RunTaskHandler))
	}

	fet := r.Group("/f")
	{
		fet.GET("/gitRep", __(fetch.FetchGitRep))
	}

	api := r.Group("/api")
	managerApi := api.Group("/m")
	managerApi.Use(auth.LoginMiddleware)
	{
		managerApi.POST("/login", __(auth.LoginHandler))
		managerApi.POST("/logout", __(auth.LogoutHandler))
		managerApi.GET("/status", __(auth.CheckStatusHandlers))
	}

	// /m will be created by CI.
	fe, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Failed to sub path `dist`: %v", err)
	}
	r.StaticFS("/m", http.FS(fe))
	r.StaticFS("/static", http.FS(public.FS))
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})
	return r
}

func __(handler func(*gin.Context) (int, interface{})) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(handler(c))
	}
}
