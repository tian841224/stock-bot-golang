package entity

type SubscriptionSymbol struct {
	ID           uint
	UserID       uint
	SymbolID     uint
	Subscription *Subscription
	StockSymbol  *StockSymbol
	User         *User
}
