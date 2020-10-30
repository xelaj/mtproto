package tlparser

type CommentType string

const (
	CommentTypeConstructor CommentType = "constructor"
	CommentTypeEnum        CommentType = "enum"
	CommentTypeMethod      CommentType = "method"
	CommentTypeParam       CommentType = "param"
	CommentTypeType        CommentType = "type"
)

func (c CommentType) String() string {
	return "@" + string(c)
}
