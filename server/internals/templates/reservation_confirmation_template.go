package templates

import "fmt"

func ReservationConfirmationTemplate(
	firstName string,
	date string,
	timeSlot string,
	tableName string,
) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Reservation Confirmed 🎉</h2>

			<p>Hello %s,</p>

			<p>Your reservation has been successfully confirmed.</p>

			<h3>Reservation Details</h3>

			<ul>
				<li><strong>Date:</strong> %s</li>
				<li><strong>Time:</strong> %s</li>
				<li><strong>Table:</strong> %s</li>
			</ul>

			<p>We look forward to serving you.</p>

			<p>Thank you for choosing Geritcht Restaurant.</p>

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
	)
}
