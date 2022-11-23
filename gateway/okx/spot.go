package okx

type Spot struct {
	Market
}

func NewSpot() Spot {
	return Spot{NewMarketWS()}
}
