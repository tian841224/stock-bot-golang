package entity

type SubscriptionSymbol struct {
	ID             uint
	SubscriptionID uint
	SymbolID       uint
	Subscription   *Subscription
	StockSymbol    *StockSymbol
}
