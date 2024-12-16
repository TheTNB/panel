package middleware

import (
	"net/http"
	"strings"

	"github.com/go-rat/chix"
	"github.com/go-rat/sessions"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cast"
)

// Entrance 确保通过正确的入口访问
func Entrance(conf *koanf.Koanf, session *sessions.Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := session.GetSession(r)
			if err != nil {
				render := chix.NewRender(w)
				defer render.Release()
				render.Status(http.StatusInternalServerError)
				render.JSON(chix.M{
					"message": err.Error(),
				})
				return
			}

			entrance := conf.String("http.entrance")
			if strings.TrimSuffix(r.URL.Path, "/") == strings.TrimSuffix(entrance, "/") {
				sess.Put("verify_entrance", true)
				render := chix.NewRender(w, r)
				defer render.Release()
				render.Redirect("/login")
				return
			}

			if !conf.Bool("app.debug") &&
				!cast.ToBool(sess.Get("verify_entrance", false)) &&
				r.URL.Path != "/robots.txt" {
				render := chix.NewRender(w)
				defer render.Release()
				render.Status(http.StatusTeapot)
				render.JSON(chix.M{
					"message": "请通过正确的入口访问",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
