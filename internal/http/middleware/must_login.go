package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
)

// MustLogin 确保已登录
func MustLogin(next http.Handler) http.Handler {
	// 白名单
	whiteList := []string{
		"/api/user/login",
		"/api/user/logout",
		"/api/user/isLogin",
		"/api/dashboard/panel",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := app.Session.GetSession(r)
		if err != nil {
			render := chix.NewRender(w)
			render.Status(http.StatusInternalServerError)
			render.JSON(chix.M{
				"message": err.Error(),
			})
		}

		// 对白名单和非 API 请求放行
		if slices.Contains(whiteList, r.URL.Path) || !strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
			return
		}

		if sess.Missing("user_id") {
			render := chix.NewRender(w)
			render.Status(http.StatusUnauthorized)
			render.JSON(chix.M{
				"message": "会话已过期，请重新登录",
			})
			return
		}

		userID := cast.ToUint(sess.Get("user_id"))
		if userID == 0 {
			render := chix.NewRender(w)
			render.Status(http.StatusUnauthorized)
			render.JSON(chix.M{
				"message": "会话无效，请重新登录",
			})
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user_id", userID)) // nolint:staticcheck
		next.ServeHTTP(w, r)
	})
}
