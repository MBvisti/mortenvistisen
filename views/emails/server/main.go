package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MBvisti/mortenvistisen/views/emails"
	"github.com/a-h/templ"
	"github.com/golang-module/carbon/v2"
	"github.com/riverqueue/river/rivertype"
)

// main is only in place to develop emails locally
func main() {
	http.Handle("/newsletter", templ.Handler(&emails.NewsletterMail{
		Title:           "Test",
		Content:         "<div><strong>dsadsa<br></strong><br></div><h2>dsadsadsad</h2><div><br><br></div><pre>dasdsadsa</pre><div><br></div><blockquote>dsadsadsadsadsaa</blockquote><div><br><br><br><a href=&quot;https://google.com&quot;>check it</a></div>",
		UnsubscribeLink: "http://localhost:8080/",
	}))

	http.Handle("/password-reset-mail", templ.Handler(&emails.PasswordReset{
		ResetPasswordLink: "https://mortenvistisen.com",
	}))

	http.Handle(
		"/books/how-start-freelancing",
		templ.Handler(&emails.BookWelcomeMail{
			ConfirmationLink: "https://mortenvistisen.com",
			UnsubscribeLink:  "https://mortenvistisen.com",
		}),
	)

	http.Handle(
		"/background-job-mail",
		templ.Handler(&emails.BackgroundJobErrorMail{
			JobID:       0,
			AttemptedAt: time.Now(),
			Kind:        "",
			MetaData:    "",
			Err:         errors.New("could not finish job"),
			AttemptedErrors: []rivertype.AttemptError{
				{
					At:      time.Now(),
					Attempt: 2,
					Error:   "bad connection",
					Trace:   "trace trace",
				},
				{
					At:      carbon.Now().SubDay().StdTime(),
					Attempt: 1,
					Error:   "bad token",
					Trace:   "trace",
				},
			},
		}),
	)

	http.Handle(
		"/user-signup-welcome-mail",
		templ.Handler(&emails.UserSignupWelcomeMail{
			ConfirmationLink: "https://mortenvistisen.com",
		}),
	)

	http.Handle(
		"/newsletter-welcome",
		templ.Handler(&emails.NewsletterWelcomeMail{
			ConfirmationLink: "https://mortenvistisen.com",
			UnsubscribeLink:  "https://mortenvistisen.com",
		}),
	)

	fmt.Println("Listening on :4444")
	if err := http.ListenAndServe(":4444", nil); err != nil {
		panic(err)
	}
}
