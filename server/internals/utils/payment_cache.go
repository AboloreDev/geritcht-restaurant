package utils

import (
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
)

func BuildPaymentCacheKey(filter *dto.PaymentFilterRequest) string {
	return fmt.Sprintf(
		"payments:all:page:%d:size:%d:amount:%s:status:%s:reference:%s",
		filter.Page,
		filter.PageSize,
		normalize(filter.Amount),
		normalize(filter.Status),
		normalize(filter.Reference),
	)
}

func GetPaymentCacheTTL(filter *dto.PaymentFilterRequest) time.Duration {
	hasFilter := filter.Amount != 0 || filter.Status != "" || filter.Reference != ""

	if hasFilter {
		return 50 * time.Second
	}
	return 5 * time.Minute
}

func BuildUserPaymentCacheKey(userID uint, filter *dto.PaymentFilterRequest) string {
	return fmt.Sprintf(
		"reservations:user:%d:page:%d:size:%d:amount:%s:status:%s:reference:%s",
		userID,
		filter.Page,
		filter.PageSize,
		normalize(filter.Amount),
		normalize(filter.Status),
		normalize(filter.Reference),
	)
}
