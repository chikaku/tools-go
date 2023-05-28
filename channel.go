package tools

// SendNonBlock try send the value v to channel s
// if the channel is full, return false
// otherwise the value will sent to channel and return true
func SendNonBlock[T any](s chan T, v T) bool {
	select {
	case s <- v:
		return true
	default:
		return false
	}
}
