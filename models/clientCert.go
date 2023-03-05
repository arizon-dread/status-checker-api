package models

import (
	"mime/multipart"
)

type ClientCert struct {
	ID       int    `json:"id" gorm:"type:int"`
	Name     string `json:"name"`
	P12      []byte `json:"-" gorm:"serializer:json"`
	Password string `json:"password"`
}

type CertUploadForm struct {
	Name     string                `form:"name" binding:"required"`
	P12      *multipart.FileHeader `form:"file" binding:"required"`
	Password string                `form:"password" binding:"required"`
}
