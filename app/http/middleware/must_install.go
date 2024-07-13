package middleware

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/translation"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/internal/services"
)

// MustInstall 确保已安装插件
func MustInstall() http.Middleware {
	return func(ctx http.Context) {
		path := ctx.Request().Path()
		translate := facades.Lang(ctx)
		var slug string
		if strings.HasPrefix(path, "/api/panel/website") {
			slug = "openresty"
		} else if strings.HasPrefix(path, "/api/panel/container") {
			slug = "podman"
		} else {
			pathArr := strings.Split(path, "/")
			if len(pathArr) < 4 {
				ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
					"message": translate.Get("errors.plugin.notExist"),
				})
				return
			}
			slug = pathArr[3]
		}

		plugin := services.NewPluginImpl().GetBySlug(slug)
		installedPlugin := services.NewPluginImpl().GetInstalledBySlug(slug)
		installedPlugins, err := services.NewPluginImpl().AllInstalled()
		if err != nil {
			ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, http.Json{
				"message": translate.Get("errors.internal"),
			})
			return
		}

		if installedPlugin.Slug != plugin.Slug {
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
				"message": translate.Get("errors.plugin.notInstalled", translation.Option{
					Replace: map[string]string{
						"slug": slug,
					},
				}),
			})
			return
		}

		pluginsMap := make(map[string]bool)

		for _, p := range installedPlugins {
			pluginsMap[p.Slug] = true
		}

		for _, require := range plugin.Requires {
			_, requireFound := pluginsMap[require]
			if !requireFound {
				ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
					"message": translate.Get("errors.plugin.dependent", translation.Option{
						Replace: map[string]string{
							"slug":       slug,
							"dependency": require,
						},
					}),
				})
				return
			}
		}

		for _, exclude := range plugin.Excludes {
			_, excludeFound := pluginsMap[exclude]
			if excludeFound {
				ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
					"message": translate.Get("errors.plugin.incompatible", translation.Option{
						Replace: map[string]string{
							"slug":    slug,
							"exclude": exclude,
						},
					}),
				})
				return
			}
		}

		ctx.Request().Next()
	}
}
