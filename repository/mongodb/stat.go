package mongodb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music-core/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStatRepository struct {
	db *mongo.Collection
}

func NewMongoStatRepository(db *mongo.Database) *MongoStatRepository {
	return &MongoStatRepository{db: db.Collection(StatCollection)}
}

var listenedTracks []uuid.UUID

func (sr *MongoStatRepository) Add(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error {
	if listenedTracks == nil {
		listenedTracks = make([]uuid.UUID, 0)
	}
	
	listenedTracks = append(listenedTracks, trackID)
	return nil
}

func (sr *MongoStatRepository) GetMostListenedMusicians(ctx context.Context, userID uuid.UUID, maxCnt int) ([]domain.UserMusiciansStat, error) {
	fmt.Printf("listened list: %v", listenedTracks)
	r := NewMongoMusicianRepository(sr.db.Database())
	var musiciansListenCount map[uuid.UUID]int64

	for _, trackID := range listenedTracks {
		musician, _ := r.GetByTrackID(ctx, trackID)
		musiciansListenCount[musician.ID] += 1
	}

	var musiciansStat []domain.UserMusiciansStat

	for musicianID, listenCount := range musiciansListenCount {
		musiciansStat = append(musiciansStat, domain.UserMusiciansStat{
			MusicianID:  musicianID,
			UserID:      userID,
			ListenCount: listenCount,
		})
	}

	return musiciansStat, nil
}

func (sr *MongoStatRepository) GetListenedGenres(ctx context.Context, userID uuid.UUID) ([]domain.UserGenresStat, error) {
	r := NewMongoGenreRepository(sr.db.Database())
	var genresListenCount map[uuid.UUID]int64

	for _, trackID := range listenedTracks {
		genres, _ := r.GetByTrackID(ctx, trackID)
		for _, genre := range genres {
			genresListenCount[genre.ID] += 1
		}
	}

	var genresStat []domain.UserGenresStat

	for genreID, listenCount := range genresListenCount {
		genresStat = append(genresStat, domain.UserGenresStat{
			GenreID:     genreID,
			UserID:      userID,
			ListenCount: listenCount,
		})
	}

	return genresStat, nil
}
