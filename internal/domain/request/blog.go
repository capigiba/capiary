package request

import "github.com/capigiba/capiary/internal/domain/constant"

type CreateBlogPostRequest struct {
	Title      string               `json:"title"`
	Blocks     []CreateBlockRequest `json:"blocks"`
	AuthorID   int                  `json:"author_id"`
	Categories []string             `json:"categories"`
}

// CreateBlockRequest describes a single block in the request.
type CreateBlockRequest struct {
	ID    int                `json:"id"`
	Type  constant.MediaType `json:"type"`
	Order int                `json:"order"`

	// for text blocks
	Paragraphs []CreateParagraphRequest `json:"paragraphs,omitempty"`

	// for heading blocks
	HeadingLevel *int    `json:"heading_level,omitempty"`
	Text         *string `json:"text,omitempty"`

	// for image/video
	Filename string `json:"filename,omitempty"`
	FileData string `json:"file_data,omitempty"` // If using base64, for instance
}

// CreateParagraphRequest mirrors the entity.Paragraph
type CreateParagraphRequest struct {
	ID      int                   `json:"id"`
	Text    string                `json:"text"`
	Formats []CreateFormatRequest `json:"formats"`
	Align   string                `json:"align"`
}

// CreateFormatRequest mirrors the entity.Format struct
type CreateFormatRequest struct {
	Type      string  `json:"type"`
	Start     int     `json:"start"`
	End       int     `json:"end"`
	Hyperlink *string `json:"hyperlink,omitempty"`
}
