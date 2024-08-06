package types

import "math/big"

type TokenChunk struct {
	Total     *big.Int
	Remaining *big.Int
	Chunk     *big.Int
}
