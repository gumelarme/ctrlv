package dynamodb

import (
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/oklog/ulid"
)

type DynamoDB struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBAPI() *DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &DynamoDB{db: dynamodb.New(sess)}
}

func NewULID() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)
	return ulid.MustNew(ulid.Now(), entropy).String()
}

func GetTimeFromId(id string) time.Time {
	milli := ulid.MustParse(id).Time()
	return time.UnixMilli(int64(milli))
}
