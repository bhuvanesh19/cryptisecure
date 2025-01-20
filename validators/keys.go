package validators

import "mime/multipart"

type AddKey struct {
	Id  string                `form:"id" binding:"required"`
	Key *multipart.FileHeader `form:"key" binding:"required"`
}
