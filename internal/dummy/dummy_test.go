package dummy

import "testing"

func TestGetName(t *testing.T) {
	t.Run("Should return the name", func(t *testing.T) {
		d := Dummy{Name: "John"}

		if d.GetName() != "John" {
			t.Errorf("Expected %s, got %s", "John", d.GetName())
		}
	})
}
