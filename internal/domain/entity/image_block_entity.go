package entity

// ImageBlock holds data for an image-based block.
type ImageBlock struct {
	ID       int     `json:"id"`
	Filename string  `json:"filename"`
	Link     *string `json:"link,omitempty"`

	FileData string `json:"file_data,omitempty"`
}
