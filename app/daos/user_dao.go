package daos

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"city-pulse/app/collections"
	"city-pulse/app/models"
)

type UserDAO struct {
	col *mongo.Collection
}

func NewUserDAO(db *mongo.Database) *UserDAO {
	return &UserDAO{col: db.Collection(collections.Users)}
}

func (d *UserDAO) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := d.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) Create(ctx context.Context, user *models.User) error {
	user.ID = bson.NewObjectID()
	user.CreatedAt = time.Now()
	_, err := d.col.InsertOne(ctx, user)
	return err
}
