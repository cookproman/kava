package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/bep3/keeper"
	"github.com/kava-labs/kava/x/bep3/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type KeeperTestSuite struct {
	suite.Suite

	keeper keeper.Keeper
	app    app.TestApp
	ctx    sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, abci.Header{Height: 1, Time: tmtime.Now()})
	keeper := tApp.GetBep3Keeper()
	suite.app = tApp
	suite.ctx = ctx
	suite.keeper = keeper
	return
}

func (suite *KeeperTestSuite) TestGetSetHtlt() {
	swapID, err := types.CalculateSwapID(randomNumberHashes[0], binanceAddrs[0], "")
	suite.NoError(err)

	heightSpan := int64(1000)
	expirationBlock := uint64(suite.ctx.BlockHeight()) + uint64(heightSpan)
	htlt := types.NewHTLT(swapID, binanceAddrs[0], kavaAddrs[0], "", "", randomNumberHashes[0], timestamps[0], coinsSingle, "50000bnb", heightSpan, false, expirationBlock)
	suite.keeper.SetHTLT(suite.ctx, htlt)

	h, found := suite.keeper.GetHTLT(suite.ctx, swapID)
	suite.True(found)
	suite.Equal(htlt, h)

	fakeSwapID, err := types.CalculateSwapID(htlt.RandomNumberHash, kavaAddrs[1], "otheraddress")
	suite.NoError(err)
	_, found = suite.keeper.GetHTLT(suite.ctx, fakeSwapID)
	suite.False(found)

	suite.keeper.DeleteHTLT(suite.ctx, swapID)
	_, found = suite.keeper.GetHTLT(suite.ctx, swapID)
	suite.False(found)
}

func (suite *KeeperTestSuite) TestIterateHtlts() {
	htlts := htlts(4)
	for _, h := range htlts {
		suite.keeper.SetHTLT(suite.ctx, h)
	}
	res := suite.keeper.GetAllHtlts(suite.ctx)
	suite.Equal(4, len(res))
}

func TestHtltTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func htlts(count int) types.HTLTs {
	var htlts types.HTLTs

	var swapIDs [][]byte
	for i := 0; i < count; i++ {
		swapID, _ := types.CalculateSwapID(randomNumberHashes[i], binanceAddrs[i], "")
		swapIDs = append(swapIDs, swapID)
	}
	h1 := types.NewHTLT(swapIDs[0], binanceAddrs[0], kavaAddrs[0], "", "", randomNumberHashes[0], timestamps[0], coinsSingle, "50000bnb", 50500, false, uint64(50500+1000))
	h2 := types.NewHTLT(swapIDs[1], binanceAddrs[1], kavaAddrs[1], "", "", randomNumberHashes[1], timestamps[1], coinsSingle, "50000bnb", 61500, false, uint64(61500+1000))
	h3 := types.NewHTLT(swapIDs[2], binanceAddrs[2], kavaAddrs[2], "", "", randomNumberHashes[2], timestamps[2], coinsSingle, "50000bnb", 72500, false, uint64(72500+1000))
	h4 := types.NewHTLT(swapIDs[3], binanceAddrs[3], kavaAddrs[3], "", "", randomNumberHashes[3], timestamps[3], coinsSingle, "50000bnb", 83500, false, uint64(83500+1000))
	htlts = append(htlts, h1, h2, h3, h4)
	return htlts
}
