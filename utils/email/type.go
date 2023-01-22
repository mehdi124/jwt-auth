package email


type EmailData struct {
	Email string
	Subject   string
	Data map[string]string
}

//TODO handle email types struct such as register and other stuff just for jwt auth

func checkEmailType(emailTemplate string) string {

	switch emailTemplate {

	case "register":
		return "verificationCode.html"
	case "withdraw":
		return "withdrawCode.html"
	default:
		return ""
	}

}