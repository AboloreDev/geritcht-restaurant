package templates

import "fmt"

func WaitlistNotificationTemplate(
	firstName string,
	date string,
	timeSlot string,
	partySize int,
	tableName string,
) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>🎉 A Table Is Available!</h2>

			<p>Hello %s,</p>

			<p>Great news! A table has become available for your waitlist request at <strong>Geritcht</strong>.</p>

			<p>Please confirm your reservation within <strong>10 minutes</strong>. If we don't hear from you before the confirmation window expires, the table will be offered to the next guest on the waitlist.</p>

			<h3>Reservation Details</h3>

			<ul>
				<li><strong>Date:</strong> %s</li>
				<li><strong>Time:</strong> %s</li>
				<li><strong>Table:</strong> %s</li>
				<li><strong>Party Size:</strong> %d</li>
			</ul>

			<p>We look forward to welcoming you to Geritcht!</p>

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
