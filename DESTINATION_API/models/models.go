package models

import "mime/multipart"

type FileMeta struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
