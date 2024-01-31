package models

import "mime/multipart"

type File struct {
	Path string
	File *multipart.FileHeader
}

var (
	FilePathNews       = "news"
	FilePathBrand      = "brands"
	FilePathVacancy    = "vacancy"
	FilePathAbout      = "about"
	FilePathFeedbacks  = "feedbacks"
	FilePathApplicants = "applicants"
	FilePathCategory   = "categories"
	FilePathApplicant  = "applicants"
	FilePathProducts   = "products"
	FilePathBanner     = "banner"
	FilePathService    = "service"
)
