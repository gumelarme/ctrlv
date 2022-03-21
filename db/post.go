package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

var (
	postTableName = aws.String("paste")
	itemsPerPage  = int64(100)
)

type PostType string

const (
	PostNote PostType = "note"
	PostCode PostType = "code"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

type Post struct {
	Id         string   `json:",omitempty" form:"Id"`
	Type       PostType `form:"Type"`
	Title      string   `form:"Title"`
	Content    string   `form:"Content"`
	Visibility string   `form:"Visibility"`
	Alias      string   `json:",omitempty" form:"Alias"`
}

func NewPostNote(title, content, alias string) *Post {
	return &Post{
		Type:    PostNote,
		Title:   title,
		Content: content,
		Alias:   alias,
	}
}

func GetPosts(last map[string]*dynamodb.AttributeValue) []Post {
	result, err := db.Scan(&dynamodb.ScanInput{
		TableName:         postTableName,
		Limit:             &itemsPerPage,
		ExclusiveStartKey: last,
	})

	if err != nil {
		panic(err)
	}

	var posts []Post
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &posts); err != nil {
		panic(errors.Wrap(err, "error while processing post"))
	}
	return posts
}

func (p *Post) Save() (string, error) {
	contentLength := len(p.Content)
	if len(p.Title) == 0 && contentLength > 0 {
		if contentLength > 10 {
			p.Title = string([]rune(p.Content)[:10]) + "..."
		} else {
			p.Title = p.Content
		}
	}

	saveExcution := p.updatePost
	if len(p.Id) == 0 {
		p.Id = NewULID()
		saveExcution = p.saveNewPost
	}

	if err := saveExcution(); err != nil {
		return "", errors.Wrap(err, "error while saving post")
	}
	return p.Id, nil
}

func (p *Post) saveNewPost() error {
	// TODO: save alias to index
	item, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return err
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: postTableName,
		Item:      item,
	})

	return errors.Wrap(err, "error while creating new post")
}

func (p *Post) updatePost() error {
	// TODO: update alias
	exp := "set #type = :type, Title = :title, Content = :content, Visibility = :visibility"
	_, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: postTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(p.Id),
			},
		},
		UpdateExpression: &exp,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":type": {
				S: aws.String(string(p.Type)),
			},
			":title": {
				S: aws.String(p.Title),
			},
			":content": {
				S: aws.String(p.Content),
			},
			":visibility": {
				S: aws.String(p.Visibility),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#type": aws.String("Type"),
		},
	})
	return errors.Wrapf(err, "error while updating post %s", p.Id)
}
