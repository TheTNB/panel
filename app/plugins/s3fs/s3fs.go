package s3fs

var (
	Name        = "S3fs"
	Description = "S3fs 通过 FUSE 挂载兼容 S3 标准的存储桶，例如Amazon S3、阿里云OSS、腾讯云COS、七牛云Kodo等。"
	Slug        = "s3fs"
	Version     = "1.9"
	Requires    = []string{}
	Excludes    = []string{}
	Install     = `bash /www/panel/scripts/s3fs/install.sh`
	Uninstall   = `bash /www/panel/scripts/s3fs/uninstall.sh`
	Update      = `bash /www/panel/scripts/s3fs/install.sh`
)
