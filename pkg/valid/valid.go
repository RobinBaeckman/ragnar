package valid

import "regexp"

func Email(e string) bool {
	r := regexp.MustCompile(`(\w+)@(\w)+\.(\w+)`)
	return r.MatchString(e)
}

func Password(p string) bool {
	r := regexp.MustCompile(`^(.{5,})$`)
	return r.MatchString(p)
}

func FirstName(fn string) bool {
	r := regexp.MustCompile(`^\b[a-zA-Z]{3,}$`)
	return r.MatchString(fn)
}

func LastName(ln string) bool {
	r := regexp.MustCompile(`^\b[a-zA-Z]{3,}$`)
	return r.MatchString(ln)
}

func UUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
