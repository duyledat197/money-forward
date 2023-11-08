package deliveries

type AuthDelivery interface{}

type authDelivery struct {
}

func NewAuthDelivery() AuthDelivery {
	return &authDelivery{}
}
