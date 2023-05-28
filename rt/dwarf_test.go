package rt

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestDwarfOffsetof(t *testing.T) {
	output := path.Join(os.TempDir(), "testbin")
	t.Log(output)
	// defer os.Remove(output)

	cmd := exec.Command("go", "build", "-o", output, "./testbin")
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	dwarf, err := ReadBinaryDWARF(output)
	if err != nil {
		t.Error(err)
	}

	for _, tt := range [][2]string{
		{"runtime.p", "numTimers"},
		{"runtime.g", "goid"},
	} {
		offset, err := DwarfOffsetof(dwarf, tt[0], tt[1])
		if err != nil {
			t.Errorf("Offsetof %s.%s error: %v", tt[0], tt[1], err)
		} else {
			t.Logf("Offsetof %s.%s is %v", tt[0], tt[1], offset)
		}
	}
}
