package templates

import "fmt"

func PasswordChangedTemplate(firstName string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Password Changed</title>
	</head>

	<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif;">
		<table width="100%%" cellspacing="0" cellpadding="0">
			<tr>
				<td align="center" style="padding:40px 0;">

					<table width="600" cellspacing="0" cellpadding="0"
						style="
							background:#ffffff;
							border-radius:12px;
							padding:40px;
						">

						<tr>
							<td align="center">
								<h1 style="
									margin:0;
									color:#111827;
									font-size:28px;
								">
									Password Updated Successfully
								</h1>
							</td>
						</tr>

						<tr>
							<td style="
								padding-top:24px;
								color:#4b5563;
								font-size:16px;
								line-height:26px;
							">
								Hello %s,
								<br><br>

								This is a confirmation that your account password was successfully changed.
							</td>
						</tr>

						<tr>
							<td style="
								padding-top:20px;
								color:#6b7280;
								font-size:14px;
								line-height:24px;
							">
								If you made this change, no further action is required.
								<br><br>

								If you did not change your password, please reset your password immediately and contact support as soon as possible.
							</td>
						</tr>

						<tr>
							<td style="
								padding-top:40px;
								border-top:1px solid #e5e7eb;
								color:#9ca3af;
								font-size:12px;
								text-align:center;
							">
								© 2026 Geritcht Restaurant. All rights reserved.
							</td>
						</tr>

					</table>

				</td>
			</tr>
		</table>
	</body>
	</html>
	`, firstName)
}
