package templates

import (
	"fmt"
	"strings"

	"github.com/AboloreDev/geritcht-restaurant/internals/events"
)

func LowStockAlertTemplate(
	adminName string,
	items []events.LowStockPayload,
) string {
	var rows strings.Builder

	for _, item := range items {
		rows.WriteString(fmt.Sprintf(`
			<tr>
				<td style="padding:10px;border:1px solid #ddd;">%s</td>
				<td style="padding:10px;border:1px solid #ddd;text-align:center;color:#d32f2f;"><strong>%.2f</strong></td>
				<td style="padding:10px;border:1px solid #ddd;text-align:center;">%.2f</td>
			</tr>
		`,
			item.Name,
			item.CurrentStock,
			item.MinThreshold,
		))
	}

	return fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif; line-height:1.6; background:#f7f7f7; padding:30px;">

		<div style="max-width:650px; margin:auto; background:#fff; border-radius:8px; padding:30px;">

			<h2 style="color:#d32f2f;">⚠️ Low Stock Alert</h2>

			<p>Hello <strong>%s</strong>,</p>

			<p>
				One or more ingredients in your inventory have reached or fallen
				below their minimum stock threshold.
			</p>

			<p>
				Please restock the following item(s) as soon as possible to avoid
				disruptions to restaurant operations.
			</p>

			<table style="width:100%%; border-collapse:collapse; margin-top:20px;">
				<tr style="background:#f2f2f2;">
					<th style="padding:10px;border:1px solid #ddd;">Ingredient</th>
					<th style="padding:10px;border:1px solid #ddd;">Current Stock</th>
					<th style="padding:10px;border:1px solid #ddd;">Minimum Threshold</th>
				</tr>

				%s

			</table>

			<br>

			<p>
				Keeping your inventory stocked helps ensure uninterrupted service
				and a great customer experience.
			</p>

			<br>

			<p>Best regards,</p>
			<p><strong>Geritcht Restaurant Management System</strong></p>

		</div>

	</body>
	</html>
	`,
		adminName,
		rows.String(),
	)
}
