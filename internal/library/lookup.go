package bot

import (
	"log"

	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/iqbalbaharum/sol-stalker/internal/adapter"
	"github.com/iqbalbaharum/sol-stalker/internal/generators"
	"github.com/iqbalbaharum/sol-stalker/internal/rpc"
	"github.com/iqbalbaharum/sol-stalker/internal/storage"
)

type LookupIndex struct {
	LookupTableIndex uint8
	LookupTableKey   string
}

// TODO: Lookup table is not store in redis!!!
func GetLookupTable(addr solana.PublicKey) (*addresslookuptable.AddressLookupTableState, error) {
	redisClient, err := adapter.GetRedisClient(3)
	if err != nil {
		log.Fatalf("Failed to get LookupTableStorage instance: %v", err)
	}

	account, err := storage.GetLookup(redisClient, addr.String())

	if err != nil {
		if err.Error() != "key not found" {
			return &addresslookuptable.AddressLookupTableState{}, err
		}
	}

	if account.Authority != nil {
		return &account, nil
	}

	resp, err := rpc.GetLookupTable(addr)

	if err != nil {
		return &addresslookuptable.AddressLookupTableState{}, err
	}

	storage.SetLookup(redisClient, addr.String(), resp)

	return &resp, nil
}

func SetLookupTable(ammId solana.PublicKey, lookup *addresslookuptable.AddressLookupTableState) error {
	redisClient, err := adapter.GetRedisClient(3)
	if err != nil {
		log.Fatalf("Failed to get LookupTableStorage instance: %v", err)
		return err
	}

	err = storage.SetLookup(redisClient, ammId.String(), *lookup)

	return nil
}

func GenerateTableLookup(addressTableLookups []generators.TxAddressTableLookup) []LookupIndex {
	var lookupIndexes []LookupIndex

	for _, lookup := range addressTableLookups {
		for _, index := range lookup.WritableIndexes {
			lookupIndexes = append(lookupIndexes, LookupIndex{
				LookupTableIndex: index,
				LookupTableKey:   lookup.AccountKey,
			})
		}
		for _, index := range lookup.ReadonlyIndexes {
			lookupIndexes = append(lookupIndexes, LookupIndex{
				LookupTableIndex: index,
				LookupTableKey:   lookup.AccountKey,
			})
		}
	}

	return lookupIndexes
}
