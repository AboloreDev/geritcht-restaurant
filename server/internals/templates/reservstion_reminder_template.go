package templates

import "fmt"

func ReservationReminderTemplate(
	firstName string,
	date string,
	timeSlot string,
	tableName string,
	partySize int,
) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Reservation Reminder 🍽️</h2>

			<p>Hello %s,</p>

			<p>This is a friendly reminder about your upcoming reservation at <strong>Geritcht</strong>.</p>

			<h3>Reservation Details</h3>

			<ul>
				<li><strong>Date:</strong> %s</li>
				<li><strong>Time:</strong> %s</li>
				<li><strong>Table:</strong> %s</li>
				<li><strong>Party Size:</strong> %v</li>
			</ul>

			<p>We look forward to welcoming you soon.</p>

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
