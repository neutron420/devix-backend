package media

import "mime/multipart"

type UploadResponse struct {
	URL      string `json:"url"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}

type FileUpload struct {
	File     multipart.File
	Header   *multipart.FileHeader
	FileType string
	MimeType string
}
