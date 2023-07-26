package rt

import (
	"os"
	"os/exec"
	"path"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func replaceProcessDWARF(t *testing.T) {
	output := path.Join(os.TempDir(), "testbin")
	defer os.Remove(output)

	cmd := exec.Command("go", "build", "-o", output, "./testbin")
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	dwarf, err := ReadBinaryDWARF(output)
	if err != nil {
		t.Error(err)
	}

	processDWARF = dwarf
	readProcessDWARFOnce.Do(func() {})
}

func TestTimerCount(t *testing.T) {
	ast := assert.New(t)
	replaceProcessDWARF(t)

	// disable GC to avoid hidden timer
	debug.SetGCPercent(-1)

	count0, _ := NumTimers()
	t.Log("initial timer count", count0)

	const n2 = 128
	wg := new(sync.WaitGroup)
	for i := 0; i < n2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.AfterFunc(time.Minute, func() {})
		}()
	}
	wg.Wait()

	// make sure timers write to P
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		runtime.Gosched()
	}

	count1, _ := NumTimers()
	t.Log("timer count", count1)

	ast.Equal(n2, int(count1-count0))
}

func TestGoid(t *testing.T) {
	replaceProcessDWARF(t)

	id, err := Goid()
	if err != nil {
		t.Error(err)
	}

	t.Logf("current goroutine id is %d", id)
}

func TestGoStack(t *testing.T) {
	replaceProcessDWARF(t)

	st, err := GoStack()
	if err != nil {
		t.Error(err)
	}

	var stackVar int
	stackPtr := uintptr(unsafe.Pointer(&stackVar))

	if !(st.Lo <= stackPtr && stackPtr < st.Hi) {
		t.Errorf("variable at %x, stack at [%x, %x)",
			stackPtr, st.Lo, st.Lo)
	}
	t.Logf("current goroutine stack [%x, %x)", st.Lo, st.Hi)
}
