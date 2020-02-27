package cwl

import (
	"fmt"
	"sort"
	"strings"
	// "encoding/json"

	"github.com/robertkrimen/otto"
)

// Input represents "CommandInputParameter".
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#CommandInputParameter
type Input struct {
	ID             string          `json:"id"`
	Label          string          `json:"label"`
	Doc            string          `json:"doc"`
	Format         string          `json:"format"`
	Binding        *Binding        `json:"inputBinding"`
	Default        *InputDefault   `json:"default"`
	Types          []Type          `json:"type"`
	SecondaryFiles []SecondaryFile `json:"secondary_files"`
	// Input.Provided is what provided by parameters.(json|yaml)
	Provided *Provided
	// Requirement ..
	RequiredType *Type
	Requirements Requirements
}

// New constructs "Input" struct from interface{}.
func (input Input) New(i interface{}) *Input {
	dest := &Input{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case "id":
				dest.ID = v.(string)
			case "type":
				dest.Types = Type{}.NewList(v)
			case "label":
				dest.Label = v.(string)
			case "doc":
				dest.Doc = v.(string)
			case "inputBinding":
				dest.Binding = Binding{}.New(v)
			case "default":
				dest.Default = InputDefault{}.New(v)
			case "format":
				dest.Format = v.(string)
			case "secondaryFiles":
				dest.SecondaryFiles = SecondaryFile{}.NewList(v)
			}
		}
		if dest.Default != nil {
			dest.Default.ID = dest.ID
		}
	case string:
		dest.Types = Type{}.NewList(x)
	case []interface{}:
		for _, v := range x {
			dest.Types = append(dest.Types, Type{}.New(v))
		}
	}
	return dest
}

// flatten
func (input *Input) flatten(typ Type, binding *Binding, prov interface{}) []string {
	flattened := []string{}
	switch typ.Type {
	case "int": // Array of Int
		tobejoined := []string{}
		for _, e := range input.Provided.Raw.([]interface{}) {
			tobejoined = append(tobejoined, fmt.Sprintf("%v", e))
		}
		if len(tobejoined) == 0 {
			return flattened
		}
		flattened = append(flattened, strings.Join(tobejoined, input.Binding.Separator))
	case "File": // Array of Files
		switch arr := input.Provided.Raw.(type) {
		case []string:
			// TODO:
		case []interface{}:
			separated := []string{}
			for _, e := range arr {
				switch v := e.(type) {
				case map[interface{}]interface{}:
					if binding != nil && binding.Prefix != "" {
						separated = append(separated, binding.Prefix)
					}
					separated = append(separated, fmt.Sprintf("%v", v["location"]))
				default:
					// TODO:
				}
			}
			// In case it's Array of Files, unlike array of int,
			// it's NOT gonna be joined with .Binding.Separator.
			flattened = append(flattened, separated...)
		}
	case "string": // Array of string
		switch arr := prov.(type) {
		case []interface{}:
			separated := []string{}
			for _, e := range arr {
				switch v := e.(type) {
				case string:
					if binding != nil && binding.Prefix != "" {
						separated = append(separated, binding.Prefix)
					}
					separated = append(separated, fmt.Sprintf("%v", v))
				default:
					// TODO
				}
			}
			flattened = append(flattened, separated...)
		default:
			// TODO
		}
	case "array":
		switch arr := prov.(type) {
		case []interface{}:
			flattened = append(flattened, input.flatten(typ.Items[0], typ.Binding, arr[0])...)
		default:
			// TODO
		}
	default:
		if input.RequiredType != nil {
			flattened = append(flattened, input.flattenWithRequiredType()...)
		} else {
			// TODO
		}
	}
	return flattened
}

func (input *Input) flattenWithRequiredType() []string {
	flattened := []string{}
	key, needed := input.Types[0].NeedRequirement()
	if !needed {
		return flattened
	}
	if input.RequiredType.Name != key {
		return flattened
	}
	switch provided := input.Provided.Raw.(type) {
	case []interface{}:
		for _, e := range provided {
			switch v := e.(type) {
			case map[interface{}]interface{}:
				for _, field := range input.RequiredType.Fields {
					val, ok := v[field.Name]
					if !ok {
						continue
					}
					if field.Binding == nil {
						// Without thinking anything, just append it!!!
						flattened = append(flattened, fmt.Sprintf("%v", val))
					} else {
						if field.Binding.Prefix != "" {
							if field.Binding.Separate {
								flattened = append(flattened, field.Binding.Prefix, fmt.Sprintf("%v", val))
							} else {
								// TODO: Join if .Separator is given
								flattened = append(flattened, fmt.Sprintf("%s%v", field.Binding.Prefix, val))
							}
						} else {
							switch v2 := val.(type) {
							case []interface{}:
								for _, val2 := range v2 {
									switch v3 := val2.(type) {
									case []interface{}:
									case map[interface{}]interface{}:
										for _, types := range input.Requirements[0].SchemaDefRequirement.Types {
											val3array := []string{}
											var val3count int = 0
											sort.Sort(types.Fields)
											for _, fields := range types.Fields {
												for key3, val3 := range v3 {
													if fields.Name == key3 {
														for _, val3type := range fields.Types {
															if val3type.Type == "" {
															} else {
																switch val3type.Type {
																case "enum":
																	for _, symbol := range val3type.Symbols {
																		if symbol == val3 {
																			val3array = append(val3array, fmt.Sprintf("%v", val3))
																			val3count = val3count + 1
																		}
																	}
																case "int":
																	if fields.Binding.Prefix != "" {
																		val3array = append(val3array, fields.Binding.Prefix, fmt.Sprintf("%v", val3))
																		val3count = val3count + 1
																	} else {
																		val3array = append(val3array, fmt.Sprintf("%v", val3))
																		val3count = val3count + 1
																	}
																}
															}
														}
													}
												}
											}
											if len(v3) == val3count {
												flattened = append(flattened, val3array...)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return flattened
}

// Flatten ...
func (input *Input) Flatten() []string {
	if input.Provided == nil {
		// In case "input.Default == nil" should be validated by usage layer.
		if input.Default != nil {
			return input.Default.Flatten(input.Binding)
		} else {
			return []string{}
		}
	}
	flattened := []string{}
	switch input.Types[0].Type {
	case "array":
		flattened = append(flattened, input.flatten(input.Types[0].Items[0], input.Types[0].Binding, input.Provided.Raw)...)
		if len(flattened) == 0 {
			return flattened
		}
	case "int":
		flattened = append(flattened, fmt.Sprintf("%v", input.Provided.Int))
	case "File":
		if input.Binding != nil && input.Binding.ValueFrom.String != "" {
			// ValueFrom case must be handled separately
			// value from ValueFrom gets stored in input.Provided.Raw in transformInput() in K8sEngine.go
			flattened = append(flattened, input.Provided.Raw.(string))
		} else {
			flattened = append(flattened, input.Provided.Entry.Location)
		}
	case "Any":
		switch v := input.Provided.Raw.(type) {
		case string:
			flattened = append(flattened, v)
		default:
			flattened = append(flattened, fmt.Sprintf("%v", v))
		}
	case "null":
		if len(input.Types) == 1 {
			return flattened
		}
		switch input.Types[1].Type {
		case "boolean":
			if input.Provided != nil && input.Provided.Raw == false {
				return flattened
			}
		default:
			// TODO other case
		}
	default:
		flattened = append(flattened, fmt.Sprintf("%v", input.Provided))
	}
	if input.Binding != nil && input.Binding.Prefix != "" {
		flattened = append([]string{input.Binding.Prefix}, flattened...)
	}

	return flattened
}

// Inputs represents "inputs" field in CWL.
type Inputs []*Input

// New constructs new "Inputs" struct.
func (ins Inputs) New(i interface{}) Inputs {
	dest := Inputs{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Input{}.New(v))
		}
	case map[string]interface{}:
		for key, v := range x {
			input := Input{}.New(v)
			input.ID = key
			dest = append(dest, input)
		}
	}
	return dest
}

// Len for sorting.
func (ins Inputs) Len() int {
	return len(ins)
}

// Less for sorting.
func (ins Inputs) Less(i, j int) bool {
	prev, next := ins[i].Binding, ins[j].Binding
	switch [2]bool{prev == nil, next == nil} {
	case [2]bool{true, true}:
		return true
	case [2]bool{false, true}:
		return prev.Position < 0
	case [2]bool{true, false}:
		return next.Position > 0
	default:
		return prev.Position <= next.Position
	}
}

// Swap for sorting.
func (ins Inputs) Swap(i, j int) {
	ins[i], ins[j] = ins[j], ins[i]
}

// ToJavaScriptVM ...
// load all context into js vm
// NOTE: NOT using this function in gen3cwl - we are using our own, better function for loading inputs context
func (ins Inputs) ToJavaScriptVM(prefix string) (*otto.Otto, error) {
	self := map[string]interface{}{}
	var (
		id        string
		path      string
		basename  string
		splitPath []string
	)
	for _, i := range ins {
		id = strings.TrimPrefix(i.ID, prefix)
		/*
			fmt.Printf("Input ID: %v\n", id)
			fmt.Println("Handling this input:")
			seeIn, _ := json.MarshalIndent(i, "", "    ")
			fmt.Println(string(seeIn))
			fmt.Println("Provided:")
			seeProvided, _ := json.MarshalIndent(i.Provided, "", "    ")
			fmt.Println(string(seeProvided))
		*/
		path, basename = "", ""
		if i.Provided != nil {
			if i.Provided.Entry != nil {
				path = i.Provided.Entry.Location
			}
			// every non-file handled here - populate with raw value
			// this is a janky way of loading context - needs to be made much more robust
			if i.Types[0].Type != "File" {
				self[id] = i.Provided.Raw
				continue
			}
		}
		if i.Default != nil && i.Default.Entry != nil {
			path = i.Default.Entry.Location
		}
		if path != "" {
			splitPath = strings.Split(path, "/")
			basename = splitPath[len(splitPath)-1]
			self[id] = map[string]interface{}{
				"path":     path,
				"location": path,
				"basename": basename,
			}
		}
	}

	// Create the js vm
	vm := otto.New()

	// Load any input variables into the vm
	if len(self) > 0 {
		if err := vm.Set("inputs", self); err != nil {
			return nil, err
		}
	}

	return vm, nil
}
