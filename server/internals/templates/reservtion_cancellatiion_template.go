package templates

import "fmt"

func ReservationCancellationTemplate(
	firstName string,
	date string,
	timeSlot string,
	tableName string,
	partySize int,
) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Reservation Cancelled ❌</h2>

			<p>Hello %s,</p>

			<p>Your reservation at <strong>Geritcht</strong> has been successfully cancelled.</p>

			<h3>Cancelled Reservation Details</h3>

			<ul>
				<li><strong>Date:</strong> %s</li>
				<li><strong>Time:</strong> %s</li>
				<li><strong>Table:</strong> %s</li>
				<li><strong>Party Size:</strong> %v</li>
			</ul>

			<p>We hope to serve you another time.</p>

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
