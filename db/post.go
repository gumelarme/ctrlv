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
	Id         string   `form:"Id"`
	Category   Category `json:",omitempty" form:"Category"`
	Title      string   `json:",omitempty" form:"Title"`
	Content    string   `json:",omitempty" form:"Content"`
	Visibility string   `json:",omitempty" form:"Visibility"`
	Alias      string   `json:",omitempty" form:"Alias"`
}

func NewPostNote(title, content, alias string) *Post {
	return &Post{
		Category: PostNote,
		Title:    title,
		Content:  content,
		Alias:    alias,
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

	saveFunc := func() error {
		return p.updatePost(p.Id, *p)
	}

	if len(p.Id) == 0 {
		p.Id = NewULID()
		saveFunc = p.saveNewPost
	}

	if err := saveFunc(); err != nil {
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

func (p *Post) updatePost(id string, post Post) error {
	// TODO: update alias
	exp := "set Category = :category, Title = :title, Content = :content, Visibility = :visibility"
	_, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: postTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(p.Id),
			},
		},
		UpdateExpression: &exp,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(p.Title),
			},
			":content": {
				S: aws.String(p.Content),
			},
			":visibility": {
				S: aws.String(p.Visibility),
			},
			":category": {
				S: aws.String(string(p.Category)),
			},
		},
	})

	if err != nil {
		return errors.Wrapf(err, "error while updating post %s", p.Id)
	}

	return nil
}

func (p *Post) Timestamp() string {
	return GetTimeFromId(p.Id).Format("Mon, 02 Jan 2006 15:04:05")
}

func GetPost(id string) (*Post, error) {
	output, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: postTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: &id,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var post Post
	if err := dynamodbattribute.UnmarshalMap(output.Item, &post); err != nil {
		return nil, err
	}

	return &post, nil
}
