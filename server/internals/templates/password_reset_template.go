package templates

import "fmt"

func PasswordResetTemplate(token string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Password Reset</title>
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
									Reset Your Password
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
								Hello,
								<br><br>

								We received a request to reset your password.
								Use the verification code below to continue.
							</td>
						</tr>

						<tr>
							<td align="center" style="padding:35px 0;">

								<div style="
									display:inline-block;
									background:#f3f4f6;
									padding:18px 36px;
									border-radius:10px;
									font-size:32px;
									font-weight:bold;
									letter-spacing:8px;
									color:#111827;
								">
									%s
								</div>

							</td>
						</tr>

						<tr>
							<td style="
								color:#6b7280;
								font-size:14px;
								line-height:22px;
							">
								This code will expire shortly.
								If you did not request a password reset,
								you can safely ignore this email.
							</td>
						</tr>

						<tr>
							<td style="
								padding-top:40px;
								text-align:center;
								font-size:12px;
								color:#9ca3af;
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
	`, token)
}
