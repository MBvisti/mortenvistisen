package validation

type InputField struct {
	Invalid    bool
	InvalidMsg string
	OldValue   string
}

type CheckBox struct {
	Invalid    bool
	InvalidMsg string
	OldValue   bool
}
