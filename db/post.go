package db

import (
	"fmt"
	"strings"

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
		return p.updatePost(*p)
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

func (p *Post) updatePost(post Post) error {
	data := map[string]string{
		"Title":      post.Title,
		"Category":   string(post.Category),
		"Content":    post.Content,
		"Visibility": post.Visibility,
	}
	_, err := UpdatePostByMap(p.Id, data)
	return err
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

func UpdatePostByMap(id string, data map[string]string) (*Post, error) {
	keys := []string{"Title", "Category", "Content", "Visibility"}

	var columns []string
	attrValue := make(map[string]*dynamodb.AttributeValue)
	for _, k := range keys {
		if val, ok := data[k]; ok {
			attrValue[":"+k] = &dynamodb.AttributeValue{S: &val}
			columns = append(columns, fmt.Sprintf("%s = :%s", k, k))
		}
	}

	// TODO: update alias
	_, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: postTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		UpdateExpression:          aws.String("set " + strings.Join(columns, ", ")),
		ExpressionAttributeValues: attrValue,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "error while updating post %s", id)
	}

	return GetPost(id)
}
