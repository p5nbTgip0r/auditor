package database

import (
	"audit/database/schema"
	"context"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jellydator/ttlcache/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	guildCache = ttlcache.New[discord.GuildID, *schema.Guild](
		ttlcache.WithTTL[discord.GuildID, *schema.Guild](12 * time.Hour),
	)
)

func init() {
	go guildCache.Start()
}

type GuildsCollection struct {
	*mongo.Collection
}

func guildIdFilter(id discord.GuildID) interface{} {
	return bson.D{{"guildID", id}}
}

func (c *GuildsCollection) GetGuild(id discord.GuildID) (*schema.Guild, error) {
	if cache := guildCache.Get(id); cache != nil {
		return cache.Value(), nil
	}

	g := &schema.Guild{ID: id}
	err := c.FindOne(context.Background(), guildIdFilter(id)).Decode(g)
	if err == nil {
		guildCache.Set(id, g, ttlcache.DefaultTTL)
	}
	return g, err
}

func (c *GuildsCollection) SetGuild(id discord.GuildID, value schema.Guild) error {
	g := &value
	if !g.ID.IsValid() {
		g.ID = id
	}
	result := c.FindOneAndReplace(context.Background(), guildIdFilter(id), *g, options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After))
	err := result.Err()
	if err == nil {
		guildCache.Set(id, g, ttlcache.DefaultTTL)
	}
	return err
}
