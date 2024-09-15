package forms

// Form struct
type Form struct {
	IsFormData  bool           `json:"form_data"`
	WithoutBody bool           `json:"without_body"`
	Body        map[string]any `json:"body"`
	Headers     map[string]any `json:"headers"`
	Data        map[string]any `json:"data"`
	Endpoint    string         `json:"endpoint"`
	Timeout     string         `json:"timeout"`
	Method      string         `json:"method"`

	Error      error
	name       string
	bodyString string
}

// GetName returns form name
func (form *Form) GetName() string {
	return form.name
}

// DataItem returns data item value
func (form *Form) DataItem(key string, defaultValue ...any) any {
	var value any
	if defaultValue != nil && len(defaultValue) > 0 {
		value = defaultValue[0]
	} else {
		value = nil
	}

	if form.Data != nil {
		val, ok := form.Data[key]
		if ok {
			value = val
		}
	}

	return value
}

// SetHeader add form header
func (form *Form) SetHeader(key string, value any) {
	if form.Headers == nil {
		form.Headers = make(map[string]any, 0)
	}

	form.Headers[key] = value
}

// IsError returns is form has error
func (form *Form) IsError() bool {
	return form.Error != nil
}
