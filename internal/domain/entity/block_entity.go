package entity

// BlockType enumerates what kind of block this is (text, image, or video).
type BlockType string

const (
	BlockTypeText    BlockType = "text"
	BlockTypeImage   BlockType = "image"
	BlockTypeVideo   BlockType = "video"
	BlockTypeHeading BlockType = "heading"
)

// Block is a generic container that references the actual content
// (text, image, or video) in the fields below.
// Only one of Text, Image, or Video should be set per block.
type Block struct {
	ID    int       `json:"id"`
	Type  BlockType `json:"type"`
	Order int       `json:"order"`

	Text    *TextBlock    `json:"text,omitempty"`
	Image   *ImageBlock   `json:"image,omitempty"`
	Video   *VideoBlock   `json:"video,omitempty"`
	Heading *HeadingBlock `json:"heading,omitempty"`
}
