package templates

import "fmt"

func VerificationEmailTemplate(token string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Email Verification</title>
	</head>
	<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif;">
		<table width="100%%" cellspacing="0" cellpadding="0">
			<tr>
				<td align="center" style="padding:40px 0;">
					<table width="600" cellspacing="0" cellpadding="0" style="background:#ffffff;border-radius:10px;padding:40px;">
						
						<tr>
							<td align="center">
								<h1 style="margin:0;color:#111827;">
									Verify Your Email
								</h1>
							</td>
						</tr>

						<tr>
							<td style="padding-top:20px;color:#4b5563;font-size:16px;line-height:24px;">
								Hello,
								<br><br>
								Thank you for signing up. Please use the verification code below to verify your email address.
							</td>
						</tr>

						<tr>
							<td align="center" style="padding:30px 0;">
								<div style="
									display:inline-block;
									padding:16px 32px;
									font-size:32px;
									font-weight:bold;
									letter-spacing:6px;
									background:#f3f4f6;
									border-radius:8px;
									color:#111827;
								">
									%s
								</div>
							</td>
						</tr>

						<tr>
							<td style="color:#6b7280;font-size:14px;line-height:22px;">
								This code will expire shortly. If you did not request this email, you can safely ignore it.
							</td>
						</tr>

						<tr>
							<td style="padding-top:30px;color:#9ca3af;font-size:12px;text-align:center;">
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
