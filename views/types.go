package views

type InputElement struct {
	Invalid    bool
	InvalidMsg string
	OldValue   string
}

type InputNumberElement struct {
	Invalid    bool
	InvalidMsg string
	OldValue   int
}

type CheckboxElement struct {
	Invalid    bool
	InvalidMsg string
	WasChecked bool
}
