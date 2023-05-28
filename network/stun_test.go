package network

import "testing"

func TestGetUDPMapAddress(t *testing.T) {
	if _, err := GetUDPMapAddress(nil); err != nil {
		t.Error(err)
	}
}
