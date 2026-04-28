package resource

type Type struct {
	Name    string `json:"name"`
	Eatable bool   `json:"eatable"`
}

var (
	Food  = Type{"Food", true}
	Wood  = Type{"Wood", false}
	Water = Type{"Water", true}
)

type Resource struct {
	Type Type `json:"type"`
	Cnt  uint `json:"count"`
}
