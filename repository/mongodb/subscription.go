package mongodb

import (
	"context"
	"github.com/hanoys/sigma-music-core/domain"
	"github.com/hanoys/sigma-music-core/ports"
	"github.com/hanoys/sigma-music-core/util"
	entity2 "github.com/hanoys/sigma-music-repository/repository/mongodb/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSubscriptionRepository struct {
	db *mongo.Collection
}

func NewMongoSubscriptionRepository(db *mongo.Database) *MongoSubscriptionRepository {
	return &MongoSubscriptionRepository{db: db.Collection(SubscriptionCollection)}
}

func (sr *MongoSubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	mongoSub := entity2.NewMongoSubscription(sub)
	_, err := sr.db.InsertOne(ctx, mongoSub)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.Subscription{}, util.WrapError(ports.ErrSubDuplicate, err)
		}
		return domain.Subscription{}, util.WrapError(ports.ErrInternalSubRepo, err)
	}

	return sub, nil
}
