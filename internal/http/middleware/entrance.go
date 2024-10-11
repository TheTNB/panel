package middleware

import (
	"net/http"
	"strings"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
)

// Entrance 确保通过正确的入口访问
func Entrance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.Session.GetSession(r)
		if err != nil {
			render := chix.NewRender(w)
			render.Status(http.StatusInternalServerError)
			render.JSON(chix.M{
				"message": err.Error(),
			})
		}

		entrance := app.Conf.String("http.entrance")
		if strings.TrimSuffix(r.URL.Path, "/") == strings.TrimSuffix(entrance, "/") {
			sess.Put("verify_entrance", true)
			render := chix.NewRender(w, r)
			render.Redirect("/login")
			return
		}

		if !app.Conf.Bool("app.debug") &&
			!cast.ToBool(sess.Get("verify_entrance", false)) &&
			r.URL.Path != "/robots.txt" {
			render := chix.NewRender(w)
			render.Status(http.StatusTeapot)
			render.JSON(chix.M{
				"message": "请通过正确的入口访问",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
