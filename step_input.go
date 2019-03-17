package cwl

// StepInput represents WorkflowStepInput.
// @see http://www.commonwl.org/v1.0/Workflow.html#WorkflowStepInput
type StepInput struct {
	ID        string
	Source    []string
	LinkMerge string
	Default   *InputDefault
	ValueFrom string
}

// New constructs a StepInput struct from any interface.
func (_ StepInput) New(i interface{}) StepInput {
	dest := StepInput{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case "id":
				dest.ID = v.(string)
			case "source":
				if list, ok := v.([]interface{}); ok {
					for _, s := range list {
						dest.Source = append(dest.Source, s.(string))
					}
				} else {
					dest.Source = append(dest.Source, v.(string))
				}
			case "linkMerge":
				dest.LinkMerge = v.(string)
			case "default":
				dest.Default = InputDefault{}.New(v)
			case "valueFrom":
				dest.ValueFrom = v.(string)
			}
		}
	}
	return dest
}

// StepInputs represents []StepInput
type StepInputs []StepInput

// NewList constructs a list of StepInput from interface.
func (_ StepInput) NewList(i interface{}) StepInputs {
	dest := StepInputs{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, StepInput{}.New(v))
		}
	case map[string]interface{}:
		for key, v := range x {
			item := make(map[string]interface{})
			item[key] = v
			dest = append(dest, StepInput{}.New(item))
		}
	default:
		dest = append(dest, StepInput{}.New(x))
	}
	return dest
}

// Len for sorting
func (s StepInputs) Len() int {
	return len(s)
}

// Less for sorting
func (s StepInputs) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// Swap for sorting
func (s StepInputs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
