package helpers

import "regexp"

func LoginValidate(email, password string) (bool, string) {
	if !isEmailValid(email) {
		return false, "email hatalı"
	}
	if !isPassValid(password) {
		return false, "password hatalı"
	}

	return true, "okey"
}

func NewWordValidate(en, tr, abb, desc string) (bool, string) {
	if len(en) == 0 {
		return false, "Terim Boş Olamaz!"
	}
	if len(tr) == 0 {
		return false, "Türkçe Karşılığı Boş Olamaz!"
	}
	if len(abb) == 0 {
		return false, "Kısaltma Boş Olamaz!"
	}
	if len(desc) == 0 {
		return false, "Açıklama Boş Olamaz!"
	}

	return true, "okey"
}

func isEmailValid(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return re.MatchString(email)
}

func isPassValid(pass string) bool {
	if len(pass) < 6 || len(pass) > 36 {
		return false
	}
	return true
}
