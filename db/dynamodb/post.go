package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gumelarme/ctrlv/config"
	"github.com/gumelarme/ctrlv/db"
	"github.com/pkg/errors"
)

func (d *DynamoDB) GetPostById(ctx context.Context, id string) (*db.Post, error) {
	output, err := d.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: &config.DynamoDB.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: &id,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var post db.Post
	if err := dynamodbattribute.UnmarshalMap(output.Item, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (d *DynamoDB) GetPosts(ctx context.Context) ([]*db.Post, error) {
	result, err := d.db.ScanWithContext(ctx, &dynamodb.ScanInput{
		TableName: &config.DynamoDB.TableName,
		Limit:     aws.Int64(int64(config.DB.ItemsPerPage)),
	})

	if err != nil {
		panic(err)
	}

	var posts []*db.Post
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &posts); err != nil {
		panic(errors.Wrap(err, "error while processing post"))
	}
	return posts, nil
}

func (d *DynamoDB) CreatePost(ctx context.Context, post *db.Post) error {
	// TODO: save alias to index
	item, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return err
	}

	_, err = d.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: &config.DynamoDB.TableName,
		Item:      item,
	})

	return errors.Wrap(err, "error while creating new post")
}

func (d *DynamoDB) UpdatePost(ctx context.Context, post *db.Post) error {
	data := map[string]string{
		"Title":      post.Title,
		"Category":   string(post.Category),
		"Content":    post.Content,
		"Visibility": post.Visibility,
	}

	_, err := d.UpdatePostByMap(ctx, post.Id, data)
	return err
}

func (d *DynamoDB) UpdatePostByMap(ctx context.Context, id string, data map[string]string) (*db.Post, error) {
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
	_, err := d.db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: &config.DynamoDB.TableName,
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

	return d.GetPostById(ctx, id)
}

// Delete a post by id
func (d *DynamoDB) DeletePost(ctx context.Context, id string) error {
	_, err := d.db.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: &config.DynamoDB.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: &id,
			},
		},
	})
	return err
}
