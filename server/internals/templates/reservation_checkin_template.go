package templates

import "fmt"

func ReservationCheckInTemplate(
	firstName string,
	date string,
	timeSlot string,
	tableName string,
	partySize int,
) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Check-In Successful ✅</h2>

			<p>Hello %s,</p>

			<p>You have successfully checked in for your reservation at <strong>Geritcht</strong>.</p>

			<h3>Reservation Details</h3>

			<ul>
				<li><strong>Date:</strong> %s</li>
				<li><strong>Time:</strong> %s</li>
				<li><strong>Table:</strong> %s</li>
				<li><strong>Party Size:</strong> %v</li>
			</ul>

			<p>Our team is preparing your dining experience, and we look forward to serving you.</p>

			<p>Thank you for dining with Geritcht.</p>

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
