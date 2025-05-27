package cache

import (
	"context"
	"errors"
	"time"

	redisClient "github.com/khuongdo95/go-pkg/caching/redis"
	"github.com/khuongdo95/go-service/internal/generated/ent"
	entuser "github.com/khuongdo95/go-service/internal/generated/ent/user"
	"github.com/khuongdo95/go-service/internal/infrastructure/global"
	redis "github.com/redis/go-redis/v9"
)

var _defaultTTL = 24 * time.Hour

type UserCache interface {
	Set(ctx context.Context, u *User) error
	Get(ctx context.Context, userID string) (*User, error)
	Delete(ctx context.Context, userID string) error
}

type userCache struct {
	ent   *ent.Client
	redis redisClient.EngineCaching
}

type User struct {
	Id string
}

func (u *User) GetId() string {
	return u.Id
}

func New(ent *ent.Client, redis redisClient.EngineCaching) UserCache {
	return &userCache{
		ent:   ent,
		redis: redis,
	}
}

func (n userCache) Get(ctx context.Context, userID string) (*User, error) {
	bytes, _, err := n.redis.Get(key(userID))
	if err == nil {
		return unmarshal(bytes)
	}
	if !errors.Is(redis.Nil, err) {
		global.Log.Error("could not get user on redis", err)
	}

	u, err := n.ent.User.Query().
		Where(entuser.ID(userID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	u0 := &User{
		Id: u.ID,
	}
	if err = n.Set(ctx, u0); err != nil {
		global.Log.Error("could not set cache", err)
	}

	return u0, nil
}

func (n userCache) Set(ctx context.Context, u *User) error {
	bytes, err := marshal(u)
	if err != nil {
		return err
	}

	_, err = n.redis.Set(key(u.GetId()), bytes, _defaultTTL)
	return err
}

func (n userCache) Delete(ctx context.Context, userID string) error {
	return n.redis.Del(key(userID))
}

func key(userID string) string {
	return "user:" + userID
}

func marshal(u *User) ([]byte, error) {
	return []byte(u.Id), nil
}

func unmarshal(bytes []byte) (*User, error) {
	return &User{
		Id: string(bytes),
	}, nil
}
