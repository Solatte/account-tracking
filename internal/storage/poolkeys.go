package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/iqbalbaharum/sol-stalker/internal/types"
	"github.com/redis/go-redis/v9"
)

type PoolKeysStorage struct {
	client *redis.Client
}

func SetPoolKeys(client *redis.Client, pKey *types.RaydiumPoolKeys) error {
	ctx := context.Background()

	data, err := json.Marshal(pKey)

	if err != nil {
		return err
	}

	if err := client.HSet(ctx, pKey.ID.String(), KEY_POOLKEYS, data).Err(); err != nil {
		return err
	}

	return nil
}

func GetPoolKeys(client *redis.Client, ammId *solana.PublicKey) (*types.RaydiumPoolKeys, error) {
	ctx := context.Background()
	data, err := client.HGet(ctx, ammId.String(), KEY_POOLKEYS).Result()
	if err != nil {
		if err == redis.Nil {
			return &types.RaydiumPoolKeys{}, errors.New("key not found")
		}
		return &types.RaydiumPoolKeys{}, err
	}

	var pKey types.RaydiumPoolKeys
	if err := json.Unmarshal([]byte(data), &pKey); err != nil {
		return &types.RaydiumPoolKeys{}, err
	}

	return &pKey, nil
}
