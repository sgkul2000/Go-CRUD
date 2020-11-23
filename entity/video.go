package entity

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VideoController contains methods on a video
type VideoController interface {
	Create(video *Video) (primitive.ObjectID, error)
}

// Video returns a video element
type Video struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	URL         string             `json:"url" bson:"url,omitempty"`
}

// Create creates a new video
func Create(video *Video) (primitive.ObjectID, error) {
	db, cancel := GetDB()
	defer cancel()
	result, err := db.Collection("videos").InsertOne(context.TODO(), video)
	if err != nil {
		log.Printf("Could not create Video: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

// Find finds all the videos in database
func Find(search string) ([]Video, error) {
	db, cancel := GetDB()
	defer cancel()

	var videos []Video
	var filter = make(map[string]map[string]string)
	if search != "" {
		filter["title"] =
			map[string]string{
				"$regex":   search,
				"$options": "i",
			}
	}
	cur, err := db.Collection("videos").Find(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var video Video
		err := cur.Decode(&video)
		if err != nil {
			log.Fatal(err)
			return videos, err
		}

		// add item our array
		videos = append(videos, video)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return videos, err
	}
	return videos, nil
}

// Update function updates an existing video
func Update(id *string, video *Video) error {
	val, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(val)
	db, cancel := GetDB()
	defer cancel()
	cur, err := db.Collection("videos").ReplaceOne(context.TODO(), bson.M{
		"_id": val,
	}, bson.M{
		"title":       video.Title,
		"description": video.Description,
		"url":         video.URL,
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println(cur.ModifiedCount)
	if cur.ModifiedCount < 1 {
		return errors.Errorf("Update failed")
	}

	return nil
}

// Delete function deletes a document by id
func Delete(id primitive.ObjectID) (bool, error) {
	db, cancel := GetDB()
	defer cancel()

	cur, err := db.Collection("videos").DeleteOne(context.TODO(), bson.M{
		"_id": id,
	})

	if err != nil {
		log.Fatalf("Delete failed: %v", err)
		return false, err
	}
	fmt.Println(cur.DeletedCount)
	if cur.DeletedCount < 1 {
		return false, nil
	}
	return true, nil
}
