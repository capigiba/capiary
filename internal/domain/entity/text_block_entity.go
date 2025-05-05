package entity

// TextBlock contains paragraphs and formatting.
type TextBlock struct {
	Paragraphs []Paragraph `json:"paragraphs"`
}

// Paragraph holds text and zero or more format ranges (bold/italic/hyperlink/etc.).
type Paragraph struct {
	ID       int       `json:"id"`
	Text     string    `json:"text"`
	Formats  []Format  `json:"formats"`
	Headings []Heading `json:"headings,omitempty"`
	Align    string    `json:"align"`
}

// FormatType enumerates the kinds of formatting that can be applied to text in a paragraph.
type FormatType string

const (
	FormatTypeBold      FormatType = "bold"
	FormatTypeItalic    FormatType = "italic"
	FormatTypeUnderline FormatType = "underline"
	FormatTypeHyperlink FormatType = "hyperlink"
)

// Format describes a particular formatting (e.g., bold, italic, hyperlink) for a substring
// from Start to End in the parent Paragraph.
type Format struct {
	Type      FormatType `json:"type"`
	Start     int        `json:"start"`
	End       int        `json:"end"`
	Hyperlink *string    `json:"hyperlink,omitempty"` // Used only if Type == "hyperlink"
}

type Heading struct {
	Type  string `json:"type"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}
