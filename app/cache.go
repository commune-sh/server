package app

import (
	config "commune/config"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/tidwall/buntdb"
)

type Cache struct {
	VerificationCodes *buntdb.DB
	GIFs              *buntdb.DB
	Events            *redis.Client
	System            *redis.Client
	Notifications     *redis.Client
}

func NewCache(conf *config.Config) (*Cache, error) {

	pdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.PostsDB,
	})

	sdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.SystemDB,
	})

	ndb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.NotificationsDB,
	})

	db, err := buntdb.Open(":memory:")
	if err != nil {
		panic(err)
	}

	gifdb, err := buntdb.Open(":memory:")
	if err != nil {
		panic(err)
	}

	c := &Cache{
		GIFs:              gifdb,
		VerificationCodes: db,
		Events:            pdb,
		System:            sdb,
		Notifications:     ndb,
	}

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("mykey", "myvalue", nil)
		return err
	})

	return c, nil
}

func (c *App) AddCodeToCache(key string, t any) error {

	serialized, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.Cache.VerificationCodes.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, string(serialized), &buntdb.SetOptions{Expires: true, TTL: time.Minute * 60})
		return err
	})
	log.Println("added to cache: ", key, t)
	return nil
}
