package middleware

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/TheTNB/panel/pkg/types"
)

// Status 检查程序状态
func Status(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch types.Status {
		case types.StatusUpgrade:
			render := chix.NewRender(w)
			render.Status(http.StatusServiceUnavailable)
			render.JSON(chix.M{
				"message": "面板升级中，请稍后刷新",
			})
			return
		case types.StatusMaintain:
			render := chix.NewRender(w)
			render.Status(http.StatusServiceUnavailable)
			render.JSON(chix.M{
				"message": "面板正在运行维护任务，请稍后刷新",
			})
			return
		case types.StatusClosed:
			render := chix.NewRender(w)
			render.Status(http.StatusForbidden)
			render.JSON(chix.M{
				"message": "面板已关闭",
			})
			return
		case types.StatusFailed:
			render := chix.NewRender(w)
			render.Status(http.StatusInternalServerError)
			render.JSON(chix.M{
				"message": "面板运行出错，请检查排除或联系支持",
			})
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}
