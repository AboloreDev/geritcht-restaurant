package templates

import "fmt"

func ReservationNoShowTemplate(
	firstName string,
	date string,
	timeSlot string,
	tableName string,
	partySize int,
) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Reservation Marked as No-Show ⚠️</h2>

			<p>Hello %s,</p>

			<p>Our records show that your reservation at <strong>Geritcht</strong> was marked as a no-show.</p>

			<h3>Reservation Details</h3>

			<ul>
				<li><strong>Date:</strong> %s</li>
				<li><strong>Time:</strong> %s</li>
				<li><strong>Table:</strong> %s</li>
				<li><strong>Party Size:</strong> %v</li>
			</ul>

			<p>If you believe this was a mistake, please contact our support team.</p>

			<br/>

			<p>Best regards,</p>
			<p><strong>Geritcht Team</strong></p>
		</body>
		</html>
	`,
		firstName,
		date,
		timeSlot,
		tableName,
		partySize,
	)
}
