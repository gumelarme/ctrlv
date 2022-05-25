package db

type Category string

const (
	PostNote Category = "note"
	PostCode Category = "code"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

type PostData struct {
	Category   Category `json:",omitempty" bson:",omitempty" form:"Category"`
	Title      string   `json:",omitempty" bson:",omitempty" form:"Title"`
	Content    string   `json:",omitempty" bson:",omitempty" form:"Content"`
	Visibility string   `json:",omitempty" bson:",omitempty" form:"Visibility"`
	Alias      string   `json:",omitempty" bson:",omitempty" form:"Alias"`
}

type Post struct {
	Id string `form:"Id"`
	PostData
}

// NewPostNote create and return a post object pointer
func NewPostNote(title, content, alias string) *Post {
	return &Post{
		PostData: PostData{
			Category: PostNote,
			Title:    title,
			Content:  content,
			Alias:    alias,
		},
	}
}
