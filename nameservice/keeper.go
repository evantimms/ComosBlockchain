package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter
// methods for various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper 

	storeKey sdk.Keeper // Unexposed key to access store from sdk

	cdc *codec.Codec // The wire codec for binary encoding/decoding

}

// Constructor
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey: storeKey,
		cdc: cdc,
	}
}

// Function to set the Whois a given name resolves to
// We first get the store object for the map[name]Whois using the
// storeKey from the Keeper
// sdk.Context holds functions to access a number of important pieces 
// of the state like blockHeight and chainID
func (k Keeper ) SetWhois(ctx sdk.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)

	// inset the <name, whois> pair into the store
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// Get the entire Whois metadata struct for a name
// Like the SetName method, first access the store using the StoreKey
// Next, instead of using the Set method on the store key, use the 
// Get method using the name casted to a []byte, return
// the result in the form of []byte
func (k Keeper) GetWhois(ctx, sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewWhois()
	}

	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
}

// The following functions get specific parameters from the store based
// on the name

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - gets the current price of a name.  If price doesn't exist yet, set to 1nametoken.
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// Get an iterator over all names in which the keys are the names and the 
// values are the whois
func (k Keeper) GetNamesIterator(ctx, sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}