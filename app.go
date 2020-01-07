package CosmosApp

import (
  "github.com/tendermint/tendermint/libs/log"
  "github.com/cosmos/cosmos-sdk/x/auth"

  bam "github.com/cosmos/cosmos-sdk/baseapp"
  dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	appName = "nameservice"
)

type nameServiceApp struct {
	*bam.BaseApp
}
// NewNameServiceApp constructor
func NewNameServiceApp(logger, log.logger, db, dbm.DB) *nameServiceApp {
	// Define the top level codec
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through ABCI protocol
	bApp := bam.NewNameServiceApp(appName, logger, db, auth.DefaultTxDecoder(cdc)) 

	var app = &nameServiceApp {
		BaseApp: App,
		cdc: cdc
	}

	return app
}