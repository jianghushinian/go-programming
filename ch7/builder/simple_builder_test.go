package builder

import "testing"

func TestSimpleBuilder(t *testing.T) {
	server := NewSimpleServerBuilder().
		Addr("localhost").
		Port(9000).
		Build()
	t.Logf("server: %v", server)
}
