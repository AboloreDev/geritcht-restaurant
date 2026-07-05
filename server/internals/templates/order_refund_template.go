package templates

import "fmt"

func OrderRefundTemplate(
	firstName string,
	orderID uint,
	reference string,
	amount int64,
	reason string,
) string {

	refund := float64(amount) / 100

	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Refund Processed 💸</h2>

			<p>Hello %s,</p>

			<p>Your refund has been successfully processed.</p>

			<h3>Refund Details</h3>

			<ul>
				<li><strong>Order ID:</strong> %d</li>
				<li><strong>Reference:</strong> %s</li>
				<li><strong>Refund Amount:</strong> %.2f</li>
				<li><strong>Reason:</strong> %s</li>
			</ul>

			<p>The amount will reflect in your account within a few business days depending on your bank.</p>

			<br/>

			<p>Best regards,</p>
			<p><strong>Your Team</strong></p>
		</body>
		</html>
	`,
		firstName,
		orderID,
		reference,
		refund,
		reason,
	)
}
