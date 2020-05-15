package zk

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubSignals(t *testing.T) {
	p0 := PubSignals([]*big.Int{new(big.Int).SetUint64(0xf1f2f3f4f5), new(big.Int).SetUint64(0x01020304)})
	p0JSON, err := json.Marshal(p0)
	require.Nil(t, err)
	var p1 PubSignals
	err = json.Unmarshal(p0JSON, &p1)
	require.Nil(t, err)
	require.Equal(t, p0, p1)
}
