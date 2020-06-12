package stater

// State is a generic store of state.
type State interface {
	// Whenever GetState is called it must return the entire represented state
	GetState() map[string]interface{}
}

// BasicState is a basic state storage just directly storing the map[string]interface{}
type BasicState struct {
	Fields map[string]interface{}
}

// GetState from basicstate
func (b *BasicState) GetState() map[string]interface{} {
	return b.Fields
}
