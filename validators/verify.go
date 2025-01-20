package validators

import "mime/multipart"

type VerifyFile struct {
	File multipart.FileHeader `form:"file" binding:"required"`
}
