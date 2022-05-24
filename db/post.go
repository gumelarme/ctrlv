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

type Post struct {
	Id         string   `form:"Id", bson:"_id,omitempty"` // need to be a *primitive.ObjectID
	Category   Category `json:",omitempty" bson:",omitempty" form:"Category"`
	Title      string   `json:",omitempty" bson:",omitempty" form:"Title"`
	Content    string   `json:",omitempty" bson:",omitempty" form:"Content"`
	Visibility string   `json:",omitempty" bson:",omitempty" form:"Visibility"`
	Alias      string   `json:",omitempty" bson:",omitempty" form:"Alias"`
}

// NewPostNote create and return a post object pointer
func NewPostNote(title, content, alias string) *Post {
	return &Post{
		Category: PostNote,
		Title:    title,
		Content:  content,
		Alias:    alias,
	}
}
