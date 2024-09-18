package middleware

import (
	"context"
	"net/http"

	"github.com/go-rat/chix"
	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/panel"
)

// MustLogin 确保已登录
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := panel.Session.GetSession(r)
		if err != nil {
			render := chix.NewRender(w)
			render.Status(http.StatusInternalServerError)
			render.JSON(chix.M{
				"message": err.Error(),
			})
		}

		if session.Missing("user_id") {
			render := chix.NewRender(w)
			render.Status(http.StatusUnauthorized)
			render.JSON(chix.M{
				"message": "会话已过期，请重新登录",
			})
			return
		}

		userID := cast.ToUint(session.Get("user_id"))
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
