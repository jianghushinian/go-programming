package builder

import "testing"

func TestServerBuilder(t *testing.T) {
	builder := NewServerBuilder()
	director := NewDirector(builder)
	server := director.ConstructServer()
	t.Logf("server: %v", server)
}
