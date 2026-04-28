package resource

type Res struct {
	Name    string `json:"name"`
	Eatable bool   `json:"eatable"`
}

var (
	Food  = Res{"Food", true}
	Wood  = Res{"Wood", false}
	Water = Res{"Water", true}
)
