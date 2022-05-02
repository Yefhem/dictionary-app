package repository

import (
	"context"

	"github.com/Yefhem/mongo/dictionary/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WordRepository interface {
	Insert(w models.Word) error
	Get(id primitive.ObjectID) (models.Word, error)
	GetAll() ([]models.Word, error)
	Update(id primitive.ObjectID, word models.Word) error
	Delete(id primitive.ObjectID) error
	Search(searchParam string, opt *options.FindOptions) ([]models.Word, int64, error)
}

type wordConnect struct {
	ctx context.Context
	col *mongo.Collection
}

func NewWordRepository(ctx context.Context, col *mongo.Collection) WordRepository {
	return &wordConnect{
		ctx: ctx,
		col: col,
	}
}

// --------> Methods
// --------> Insert to New Word
func (c *wordConnect) Insert(w models.Word) error {
	_, err := c.col.InsertOne(c.ctx, w)
	if err != nil {
		return err
	}
	return nil
}

// --------> Get to Single Word
func (c *wordConnect) Get(id primitive.ObjectID) (models.Word, error) {
	word := models.Word{}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	if err := c.col.FindOne(c.ctx, filter).Decode(&word); err != nil {
		return word, err
	}
	return word, nil
}

// --------> Get to All Words
func (c *wordConnect) GetAll() ([]models.Word, error) {
	var word models.Word
	var words []models.Word

	filter := bson.D{{}}

	cursor, err := c.col.Find(c.ctx, filter)
	if err != nil {
		defer cursor.Close(c.ctx)
		return words, err
	}

	for cursor.Next(c.ctx) {
		if err := cursor.Decode(&word); err != nil {
			return words, err
		}
		words = append(words, word)
	}

	cursor.Close(c.ctx)

	if len(words) == 0 {
		return words, mongo.ErrNoDocuments
	}

	return words, nil
}

// --------> Update to Word
func (c *wordConnect) Update(id primitive.ObjectID, word models.Word) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	updates := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "english", Value: word.English}, {Key: "turkish", Value: word.Turkish}, {Key: "abbreviation", Value: word.Abbreviation}, {Key: "description", Value: word.Description}}}}

	_, err := c.col.UpdateMany(
		c.ctx,
		filter,
		updates,
	)
	if err != nil {
		return err
	}

	return nil

}

// --------> Delete to Word
func (c *wordConnect) Delete(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := c.col.DeleteOne(c.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// --------> Search to Word
func (c *wordConnect) Search(searchParam string, opt *options.FindOptions) ([]models.Word, int64, error) {

	var word models.Word
	var words []models.Word

	filter := bson.M{
		"$or": []bson.M{
			{
				"english": bson.M{
					"$regex": primitive.Regex{
						Pattern: searchParam,
						Options: "i",
					},
				},
			},
			{
				"turkish": bson.M{
					"$regex": primitive.Regex{
						Pattern: searchParam,
						Options: "i",
					},
				},
			},
		},
	}

	total, _ := c.col.CountDocuments(c.ctx, filter)
	if total == 0 {
		return words, total, mongo.ErrNoDocuments
	}

	cursor, err := c.col.Find(c.ctx, filter, opt)
	if err != nil {
		defer cursor.Close(c.ctx)
		return words, total, err
	}
	for cursor.Next(c.ctx) {
		if err := cursor.Decode(&word); err != nil {
			return words, total, err
		}
		words = append(words, word)
	}

	cursor.Close(c.ctx)

	return words, total, nil
}
