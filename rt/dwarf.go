package rt

import (
	"debug/dwarf"
	"debug/elf"
	"errors"
	"os"
	"sync"
)

// ReadBinaryDWARF read DWARF from input ELF binary file
func ReadBinaryDWARF(path string) (*dwarf.Data, error) {
	bin, err := elf.Open(path)
	if err != nil {
		return nil, err
	}

	dwraf, err := bin.DWARF()
	if err != nil {
		return nil, err
	}

	return dwraf, nil
}

var (
	processDWARF          *dwarf.Data
	readProcessDWARFOnce  sync.Once
	readProcessDWARFError error
)

// Offsetof is similar with unsafe.Offsetof, but can used with unexposed
// structure/filed like some runtime structure runtime.g/runtime.p.
// This function dependent on DWARF in current process ELF binary os.Args[0]
func Offsetof(structName, fieldName string) (uintptr, error) {
	readProcessDWARFOnce.Do(func() {
		data, err := ReadBinaryDWARF(os.Args[0])
		if err != nil {
			readProcessDWARFError = err
			return
		}
		processDWARF = data
	})

	if readProcessDWARFError != nil {
		return 0, readProcessDWARFError
	}

	return DwarfOffsetof(processDWARF, structName, fieldName)
}

// DwarfOffsetof get Offsetof from specified DWARF data
func DwarfOffsetof(dwraf *dwarf.Data, structName, fieldName string) (uintptr, error) {
	structFound := false
	for r := dwraf.Reader(); ; {
		entry, err := r.Next()
		if err != nil || entry == nil {
			break
		}

		if !structFound {
			if entry.Tag != dwarf.TagStructType {
				continue
			}

			nameAttr := entry.AttrField(dwarf.AttrName)
			name, ok := nameAttr.Val.(string)
			structFound = ok && name == structName
			continue
		}

		if entry.Tag != dwarf.TagMember {
			break
		}

		nameAttr := entry.AttrField(dwarf.AttrName)
		name, ok := nameAttr.Val.(string)
		if ok && name == fieldName {
			loc := entry.AttrField(dwarf.AttrDataMemberLoc)
			val, ok := loc.Val.(int64)
			if !ok {
				return 0, errors.New("offset value not int64")
			}
			return uintptr(val), nil
		}
	}

	if !structFound {
		return 0, errors.New("struct type not found in DWARF")
	}
	return 0, errors.New("member not found in DWARF")
}
