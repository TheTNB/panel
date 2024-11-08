module github.com/TheTNB/panel

go 1.23

require (
	github.com/bddjr/hlfhr v1.1.3
	github.com/beevik/ntp v1.4.3
	github.com/creack/pty v1.1.24
	github.com/expr-lang/expr v1.16.9
	github.com/glebarez/sqlite v1.11.0
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-gormigrate/gormigrate/v2 v2.1.3
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.22.1
	github.com/go-rat/chix v1.1.4
	github.com/go-rat/gormstore v1.0.6
	github.com/go-rat/sessions v1.0.11
	github.com/go-rat/utils v1.0.7
	github.com/go-resty/resty/v2 v2.15.3
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-cz/httplog v0.0.0-20241002114323-98e09d6f537a
	github.com/gorilla/websocket v1.5.3
	github.com/hashicorp/go-version v1.7.0
	github.com/knadh/koanf/parsers/yaml v0.1.0
	github.com/knadh/koanf/providers/file v1.1.2
	github.com/knadh/koanf/v2 v2.1.2
	github.com/lib/pq v1.10.9
	github.com/libdns/alidns v1.0.3
	github.com/libdns/cloudflare v0.1.1
	github.com/libdns/huaweicloud v0.2.2
	github.com/libdns/libdns v0.2.2
	github.com/libdns/tencentcloud v1.1.0
	github.com/mholt/acmez/v2 v2.0.3
	github.com/orandin/slog-gorm v1.4.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/samber/lo v1.47.0
	github.com/sethvargo/go-limiter v1.0.0
	github.com/shirou/gopsutil v2.21.11+incompatible
	github.com/spf13/cast v1.7.0
	github.com/stretchr/testify v1.9.0
	github.com/tufanbarisyildirim/gonginx v0.0.0-20241013191809-e73b7dd454e8
	github.com/urfave/cli/v3 v3.0.0-alpha9.2
	go.uber.org/zap v1.27.0
	go.uber.org/zap/exp v0.3.0
	golang.org/x/crypto v0.29.0
	golang.org/x/net v0.31.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/gorm v1.25.12
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/devhaozi/huaweicloud-sdk-go-v3 v0.0.0-20241018211007-bbebb6de5db7 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-rat/securecookie v1.0.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gofiber/schema v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jaevor/go-nanoid v1.4.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.1033 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.1033 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	modernc.org/libc v1.60.1 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/sqlite v1.32.0 // indirect
)

replace github.com/mholt/acmez/v2 => github.com/TheTNB/acmez/v2 v2.0.0-20241025203320-cc718c4c870b
