package media

import "mime/multipart"

// UploadResponse is returned after a successful file upload.
type UploadResponse struct {
	URL      string `json:"url"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}

// FileUpload represents a parsed file upload.
type FileUpload struct {
	File     multipart.File
	Header   *multipart.FileHeader
	FileType string // "image" or "video"
	MimeType string
}
