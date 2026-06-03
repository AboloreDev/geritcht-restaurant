package payment

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/rpip/paystack-go"
)

type Paystack struct {
	client *paystack.Client
}

func NewClient(config *config.PaystackConfig) *Paystack {
	client := paystack.NewClient(config.PaystackSecretKey, nil)

	return &Paystack{
		client: client,
	}
}
