package valid

import "regexp"

func IsEmail(e string) bool {
	r := regexp.MustCompile(`(\w+)@(\w)+\.(\w+)`)
	return r.MatchString(e)
}

func IsPassword(p string) bool {
	r := regexp.MustCompile(`^(.{8,})$`)
	return r.MatchString(p)
}

func IsFirstName(fn string) bool {
	r := regexp.MustCompile(`^\b[a-zA-Z]{3,}$`)
	return r.MatchString(fn)
}

func IsLastName(ln string) bool {
	r := regexp.MustCompile(`^\b[a-zA-Z]{3,}$`)
	return r.MatchString(ln)
}

func IsUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
