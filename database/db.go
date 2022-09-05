package database

import (
	"audit/util"
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	uri, err := util.GetEnvRequired("MONGODB_URI")
	if err != nil {
		return nil, err
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
	Database = client.Database(util.GetEnvDefault("MONGODB_DATABASE", "auditBot"))
	Collections = &AuditCollections{
		Guilds: &GuildsCollection{Database.Collection(util.GetEnvDefault("MONGODB_GUILDS_COLLECTION", "guilds"))},
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

func Disconnect(ctx context.Context) error {
	log.Debug().Msg("MongoDB disconnect..")
	return Client.Disconnect(ctx)
}
