package form

type Validation struct {
	Name    string
	Args    []interface{}
	Message string
}

func (v Validation) Validate(value interface{}) bool {

	return true
}

type Validations []Validation

func (v Validations) Validate() bool {
	return false
}

func (v Validations) String() string {
	return ""
}
