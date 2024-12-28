package request

type ContainerImageID struct {
	ID string `json:"id" form:"id"`
}

type ContainerImagePull struct {
	Name     string `form:"name" json:"name" validate:"required"`
	Auth     bool   `form:"auth" json:"auth"`
	Username string `form:"username" json:"username" validate:"requiredIf:Auth,true"`
	Password string `form:"password" json:"password" validate:"requiredIf:Auth,true"`
}
