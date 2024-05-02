package dummy

type Dummy struct {
	Name string
}

func (d *Dummy) GetName() string {
	return d.Name
}
