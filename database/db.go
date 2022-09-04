package database

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type AuditCollections struct {
	Guilds *GuildsCollection
}

var (
	Client      *mongo.Client
	Database    *mongo.Database
	Collections *AuditCollections
)

func Initialize(ctx context.Context) (*mongo.Client, error) {
	uri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		return nil, fmt.Errorf("environment variable 'MONGODB_URI' must be set")
	}

	log.Debug().Msg("Creating MongoDB client..")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("MongoDB client creation successful")

	setupVariables(client)
	err = createIndexes(ctx)
	if err != nil {
		return nil, err
	}

	return Client, nil
}

func setupVariables(client *mongo.Client) {
	Client = client
	Database = client.Database(getEnvDefault("MONGODB_DATABASE", "auditBot"))
	Collections = &AuditCollections{
		Guilds: &GuildsCollection{Database.Collection(getEnvDefault("MONGODB_GUILDS_COLLECTION", "guilds"))},
	}
}

func createIndexes(ctx context.Context) error {
	log.Debug().Msg("Creating MongoDB indexes..")
	guildIndex := mongo.IndexModel{
		Keys:    bson.D{{"guildID", -1}},
		Options: options.Index(),
	}
	guildIndex.Options.SetUnique(true)
	i, err := Collections.Guilds.Indexes().CreateOne(ctx, guildIndex)
	if err != nil {
		return err
	}

	log.Debug().Interface("indexName", i).Msg("Created index for guilds")
	return nil
}

func getEnvDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func Disconnect(ctx context.Context) error {
	log.Debug().Msg("MongoDB disconnect..")
	return Client.Disconnect(ctx)
}
