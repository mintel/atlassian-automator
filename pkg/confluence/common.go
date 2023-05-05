package confluence

type BodyType struct {
	Representation string `json:"representation,omitempty" structs:"representation,omitempty"`
	Value          string `json:"value,omitempty" structs:"value,omitempty"`
}
