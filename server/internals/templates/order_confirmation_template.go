package templates

import (
	"fmt"

	"github.com/AboloreDev/geritcht-restaurant/internals/events"
)

func OrderConfirmationTemplate(
	firstName string,
	orderID uint,
	amount int64,
	reference string,
	items []events.OrderItemPayload,
) string {

	total := float64(amount) / 100

	// build items HTML
	var itemsHTML string
	for _, item := range items {
		itemTotal := float64(item.Price) * float64(item.Quantity)

		itemsHTML += fmt.Sprintf(`
			<li>
				<strong>%s</strong> x%d — %.2f
			</li>
		`, item.Name, item.Quantity, itemTotal)
	}

	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Order Confirmation ✅</h2>

			<p>Hello %s,</p>

			<p>Your payment was successful and your order has been confirmed.</p>

			<h3>Order Details</h3>

			<ul>
				<li><strong>Order ID:</strong> %d</li>
				<li><strong>Reference:</strong> %s</li>
			</ul>

			<h3>Items</h3>
			<ul>
				%s
			</ul>

			<h3>Payment Summary</h3>
			<ul>
				<li><strong>Total Paid:</strong> %.2f</li>
			</ul>

			<p>Your order is now being processed.</p>

			<br/>

			<p>Best regards,</p>
			<p><strong>Your Team</strong></p>
		</body>
		</html>
	`,
		firstName,
		orderID,
		reference,
		itemsHTML,
		total,
	)
}
