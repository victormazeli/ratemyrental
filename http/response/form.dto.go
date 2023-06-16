package response

import "mime/multipart"

type Form struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}
