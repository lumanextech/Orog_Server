package x_err

const (
	NotFoundMarket      = 200001
	MarketTradeError    = 200002
	GetAccountError     = 200003
	InsufficientBalance = 200004
	ErrPageInvalid      = 200005
	ErrSizeInvalid      = 200006
	ErrMaxSizeInvalid   = 200007
	ErrDBQueryError     = 200008
	ErrRedisGetError    = 200009

	ErrInvalidParam = 200010

	ErrCodeInvalidMarketAddress    = 200011
	ErrCodeRedisGet                = 200012
	ErrCodeJsonMarshal             = 200013
	ErrCodeInvalidQuoteMintAddress = 200014
)

func init() {
	initBase()

	message[NotFoundMarket] = "market not found"
	message[MarketTradeError] = "market trade error"
	message[GetAccountError] = "get account error"
	message[InsufficientBalance] = "insufficient balance"
	message[ErrPageInvalid] = "page invalid"
	message[ErrSizeInvalid] = "size invalid"
	message[ErrMaxSizeInvalid] = "max size invalid"

	message[ErrDBQueryError] = "db query error"
	message[ErrCodeInvalidMarketAddress] = "invalid market address"
}
