package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/go-rat/chix"
	"github.com/go-rat/sessions"
	"github.com/spf13/cast"
)

// MustLogin 确保已登录
func MustLogin(session *sessions.Manager) func(next http.Handler) http.Handler {
	// 白名单
	whiteList := []string{
		"/api/user/key",
		"/api/user/login",
		"/api/user/logout",
		"/api/user/isLogin",
		"/api/dashboard/panel",
	}
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

			// 对白名单和非 API 请求放行
			if slices.Contains(whiteList, r.URL.Path) || !strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
				return
			}

			if sess.Missing("user_id") {
				render := chix.NewRender(w)
				defer render.Release()
				render.Status(http.StatusUnauthorized)
				render.JSON(chix.M{
					"message": "会话已过期，请重新登录",
				})
				return
			}

			userID := cast.ToUint(sess.Get("user_id"))
			if userID == 0 {
				render := chix.NewRender(w)
				defer render.Release()
				render.Status(http.StatusUnauthorized)
				render.JSON(chix.M{
					"message": "会话无效，请重新登录",
				})
				return
			}

			safeLogin := cast.ToBool(sess.Get("safe_login"))
			if safeLogin {
				safeClientHash := cast.ToString(sess.Get("safe_client"))
				ip, _, _ := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
				clientHash := fmt.Sprintf("%x", sha256.Sum256([]byte(ip)))
				if safeClientHash != clientHash || safeClientHash == "" {
					render := chix.NewRender(w)
					defer render.Release()
					render.Status(http.StatusUnauthorized)
					render.JSON(chix.M{
						"message": "客户端IP/UA变化，请重新登录",
					})
					return
				}
			}

			r = r.WithContext(context.WithValue(r.Context(), "user_id", userID)) // nolint:staticcheck
			next.ServeHTTP(w, r)
		})
	}
}
