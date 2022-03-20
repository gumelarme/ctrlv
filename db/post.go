package db

import (
	"log"

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

type Post struct {
	Id      string `json:",omitempty"`
	Type    PostType
	Title   string
	Content string
	Alias   string `json:",omitempty"`
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
	log.Println(*result.Count)
	if dynamodbattribute.UnmarshalListOfMaps(result.Items, &posts) != nil {
		log.Println()
		panic(err)
	}
	return posts
}

func (p *Post) Save() (string, error) {
	if len(p.Title) == 0 {
		p.Title = string([]rune(p.Content)[:10]) + "..."
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
	exp := "set Type = :type, Title = :title, Content = :content"
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
		},
	})
	return errors.Wrapf(err, "error while updating post %s", p.Id)
}
