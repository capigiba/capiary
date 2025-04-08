package entity

// VideoBlock holds data for a video-based block.
type VideoBlock struct {
	ID       int     `json:"id"`
	Filename string  `json:"filename"`
	Link     *string `json:"link,omitempty"`
}
