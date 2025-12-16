package repository

import (
	"context"
	"errors"
	"time"
	"uas/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAchievementRepository interface {
	FindByHexIDs(ctx context.Context, ids []string) ([]models.AchievementDetail, error)
	Insert(ctx context.Context, a models.AchievementDetail) (*mongo.InsertOneResult, error)
	UpdateByHexID(ctx context.Context, hexID string, update bson.M) error
	DeleteByHexID(ctx context.Context, hexID string) error
	PushHistoryByHexID(ctx context.Context, hexID string, status string) error
	FindByHexID(ctx context.Context, hexID string) (models.AchievementDetail, error)
	GetHistory(ctx context.Context, hexID string) ([]bson.M, error)
}

type mongoAchievementRepo struct {
	col *mongo.Collection
}

func (r *mongoAchievementRepo) FindByHexIDs(
	ctx context.Context,
	ids []string,
) ([]models.AchievementDetail, error) {

	var objIDs []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	filter := bson.M{
		"_id": bson.M{"$in": objIDs},
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.AchievementDetail
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func NewMongoAchievementRepository(db *mongo.Database) MongoAchievementRepository {
	return &mongoAchievementRepo{
		col: db.Collection("achievements"),
	}
}

func (r *mongoAchievementRepo) Insert(ctx context.Context, a models.AchievementDetail) (*mongo.InsertOneResult, error) {
	return r.col.InsertOne(ctx, a)
}

func (r *mongoAchievementRepo) UpdateByHexID(ctx context.Context, hexID string, update bson.M) error {
	objectId, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}

	_, err = r.col.UpdateByID(ctx, objectId, update)
	return err
}

func (r *mongoAchievementRepo) DeleteByHexID(ctx context.Context, hexID string) error {
	objectId, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": objectId})
	return err
}

func (r *mongoAchievementRepo) PushHistoryByHexID(ctx context.Context, hexID string, status string) error {
	objectId, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"history": bson.M{
				"status":    status,
				"timestamp": primitive.NewDateTimeFromTime(time.Now()),
				"changedBy": "system",
			},
		},
	}

	_, err = r.col.UpdateByID(ctx, objectId, update)
	return err
}

func (r *mongoAchievementRepo) FindByHexID(ctx context.Context, hexID string) (models.AchievementDetail, error) {
	objectId, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return models.AchievementDetail{}, err
	}

	var detail models.AchievementDetail
	err = r.col.FindOne(ctx, bson.M{"_id": objectId}).Decode(&detail)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return detail, nil
	}

	return detail, err
}

func (r *mongoAchievementRepo) GetHistory(ctx context.Context, hexID string) ([]bson.M, error) {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}
	projection := bson.M{"history": 1}

	var result struct {
		History []bson.M `bson:"history"`
	}

	err = r.col.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.History, nil
}
