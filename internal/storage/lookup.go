package storage

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/gagliardetto/solana-go"
	lookup "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/redis/go-redis/v9"
)

var (
	lookupTableStorageInstance *LookupTableStorage
	once                       sync.Once
)

type LookupTableAccount struct {
	Address []solana.PublicKey
}

type LookupTableStorage struct {
	client  *redis.Client
	storage map[string]lookup.AddressLookupTableState
}

func NewLookupTableStorage(client *redis.Client) *LookupTableStorage {
	return &LookupTableStorage{}
}

func GetLutInstance(db int) (*LookupTableStorage, error) {
	var err error
	once.Do(func() {
		client, clientErr := adapter.GetRedisClient(db)
		if clientErr != nil {
			err = clientErr
			return
		}
		lookupTableStorageInstance = &LookupTableStorage{
			client:  client,
			storage: make(map[string]lookup.AddressLookupTableState),
		}
	})
	if err != nil {
		return nil, err
	}
	return lookupTableStorageInstance, nil
}

func SetLookup(client *redis.Client, ammId string, lut lookup.AddressLookupTableState) error {
	ctx := context.Background()

	// Serialize the AddressLookupTableState to JSON
	data, err := json.Marshal(lut)
	if err != nil {
		return err
	}

	// Store the serialized data in Redis
	if err := client.HSet(ctx, ammId, KEY_LOOKUP, data).Err(); err != nil {
		return err
	}

	return nil
}

func GetLookup(client *redis.Client, ammId string) (lookup.AddressLookupTableState, error) {
	ctx := context.Background()

	// Retrieve the data from Redis
	data, err := client.HGet(ctx, KEY_LOOKUP, ammId).Result()
	if err != nil {
		if err == redis.Nil {
			return lookup.AddressLookupTableState{}, errors.New("key not found")
		}
		return lookup.AddressLookupTableState{}, err
	}

	// Deserialize the JSON data back into AddressLookupTableState
	var account lookup.AddressLookupTableState
	if err := json.Unmarshal([]byte(data), &account); err != nil {
		return lookup.AddressLookupTableState{}, err
	}

	return account, nil
}
