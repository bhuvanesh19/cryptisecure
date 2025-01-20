package validators

import "mime/multipart"

type SignFileRequest struct {
	ID   string               `form:"id" binding:"required"`
	File multipart.FileHeader `form:"file" binding:"required"`
}
