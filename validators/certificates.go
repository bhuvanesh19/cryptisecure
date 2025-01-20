package validators

import "mime/multipart"

type AddCertificateRequest struct {
	Name   string                `form:"name" binding:"required"`
	Group  string                `form:"group" binding:"required"`
	Domain string                `form:"domain" binding:"required"`
	File   *multipart.FileHeader `form:"file" binding:"required"`
}

type GetCertificateRequest struct {
	ID    string `form:"id"`
	Domin string `form:"domain"`
	Name  string `form:"name"`
	Group string `form:"group"`
}
