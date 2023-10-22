package main

import "strings"

type email struct { // not allowing it to blead out of package
	from, to, subject, body string
}

type EmailBuilder struct { // available type; build up an email but not going to expose parts of email directly
	email email
}

func (e *EmailBuilder) From(from string) *EmailBuilder {
	// validation layer
	if !strings.Contains(from, "@") {
		panic("email should contain @ character")
	}
	e.email.from = from
	return e
}

func (e *EmailBuilder) To(to string) *EmailBuilder {
	e.email.to = to
	return e
}

func (e *EmailBuilder) Subject(subject string) *EmailBuilder {
	e.email.subject = subject
	return e
}

func (e *EmailBuilder) Body(body string) *EmailBuilder {
	e.email.body = body
	return e
}

func sendMailImpl(email *email) {

}

// SendEmail is the function that end user is meant to be using
// whenever someone is calling SendEmail, they have to provide body of function of type build which doesn't return any value
type build func(builder *EmailBuilder)

func SendEmail(action build) {
	builder := EmailBuilder{}    // we create a builder
	action(&builder)             // invoke action on a builder pointer, which populates using EmailBuilder, actually email with all the values passed in. This already applies all the validations before putting in values to email object
	sendMailImpl(&builder.email) // then we finally pass it in some internal method which needs type email internal to system
}

func main() {
	SendEmail(func(b *EmailBuilder) {
		b.From("ishan123@gmail.com").To("someone1@eufn.com").Subject("Something").Body("Hi!")
	})
}
