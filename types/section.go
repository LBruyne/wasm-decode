package types

import (
	"errors"
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"io"
)

type SectionID byte

const (
	SectionIDCustom   SectionID = 0
	SectionIDType     SectionID = 1
	SectionIDImport   SectionID = 2
	SectionIDFunction SectionID = 3
	SectionIDTable    SectionID = 4
	SectionIDMemory   SectionID = 5
	SectionIDGlobal   SectionID = 6
	SectionIDExport   SectionID = 7
	SectionIDStart    SectionID = 8
	SectionIDElement  SectionID = 9
	SectionIDCode     SectionID = 10
	SectionIDData     SectionID = 11
)

// readSections read each section continuously until the end of file or meet an error
func (m *Module) readSections(r io.Reader) error {
	for {
		// read each section
		if err := m.readSection(r); errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
	}
}

// readSection read each section according to the section id
func (m *Module) readSection(r io.Reader) error {
	// read section id
	b := make([]byte, 1)
	if _, err := io.ReadFull(r, b); err != nil {
		return fmt.Errorf("read section id: %w", err)
	}

	// read section size
	ss, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of section for id=%d: %w", SectionID(b[0]), err)
	}

	// decode section according to its id
	switch SectionID(b[0]) {
	case SectionIDCustom:
		err = m.readSectionCustom(r, ss)
	case SectionIDType:
		err = m.readSectionType(r, ss)
	case SectionIDImport:
		err = m.readSectionImport(r, ss)
	case SectionIDFunction:
		err = m.readSectionFunction(r, ss)
	case SectionIDTable:
		err = m.readSectionTable(r, ss)
	case SectionIDMemory:
		err = m.readSectionMemory(r, ss)
	case SectionIDGlobal:
		err = m.readSectionGlobal(r, ss)
	case SectionIDExport:
		err = m.readSectionExport(r, ss)
	case SectionIDStart:
		err = m.readSectionStart(r, ss)
	case SectionIDElement:
		err = m.readSectionElement(r, ss)
	case SectionIDCode:
		err = m.readSectionCode(r, ss)
	case SectionIDData:
		err = m.readSectionData(r, ss)
	default:
		err = errors.New("invalid section id")
	}

	if err != nil {
		return fmt.Errorf("read section for %d: %w", SectionID(b[0]), err)
	}
	return nil
}

type CustomSec struct {
	Name  string
	Bytes []byte
}

func (m *Module) readSectionCustom(r io.Reader, ss uint32) error {
	// get name
	ns, n, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("read size of custom section name: %w", err)
	}

	buf := make([]byte, ns)
	if _, err := io.ReadFull(r, buf); err != nil {
		return fmt.Errorf("read bytes of custom section name: %w", err)
	}

	ss -= ns + uint32(n) // TODO has risk to lose precision

	bs := make([]byte, ss)
	if _, err := io.ReadFull(r, bs); err != nil {
		return fmt.Errorf("read custom section bytes: %w", err)
	}

	m.SecCustom = &CustomSec{
		Name:  string(buf),
		Bytes: bs,
	}
	return nil
}

func (m *Module) readSectionType(r io.Reader, size uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecType = make([]*FunctionType, vs)
	for i := range m.SecType {
		m.SecType[i], err = readFunctionType(r)
		if err != nil {
			return fmt.Errorf("read %d-th function type: %w", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionImport(r io.Reader, size uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecImport = make([]*ImportSegment, vs)
	for i := range m.SecImport {
		m.SecImport[i], err = readImportSegment(r)
		if err != nil {
			return fmt.Errorf("read %v-th import segment: %w", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionFunction(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecFunction = make([]uint32, vs)
	for i := range m.SecFunction {
		m.SecFunction[i], _, err = common.DecodeUint32(r)
		if err != nil {
			return fmt.Errorf("read %v-th function's type index: %w", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionTable(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecTable = make([]*TableType, vs)
	for i := range m.SecTable {
		m.SecTable[i], err = readTableType(r)
		if err != nil {
			return fmt.Errorf("read %v-th table type: %w", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionMemory(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecMemory = make([]*MemoryType, vs)
	for i := range m.SecMemory {
		m.SecMemory[i], err = readMemoryType(r)
		if err != nil {
			return fmt.Errorf("read %v-th memory type: %w", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionGlobal(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecGlobal = make([]*GlobalSegment, vs)
	for i := range m.SecGlobal {
		m.SecGlobal[i], err = readGlobalSegment(r)
		if err != nil {
			return fmt.Errorf("read %v-th global segment: %w ", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionExport(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecExport = make([]*ExportSegment, vs)
	for i := range m.SecExport {
		m.SecExport[i], err = readExportSegment(r)
		if err != nil {
			return fmt.Errorf("read %v-th export segment: %w ", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionStart(r io.Reader, ss uint32) error {
	idx, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get funcIdx of start section: %w", err)
	}

	m.SecStart = idx
	return nil
}

func (m *Module) readSectionElement(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecElement = make([]*ElementSegment, vs)
	for i := range m.SecElement {
		m.SecElement[i], err = readElementSegment(r)
		if err != nil {
			return fmt.Errorf("read %v-th element segment: %w ", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionCode(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecCode = make([]*CodeSegment, vs)
	for i := range m.SecCode {
		m.SecCode[i], err = readCodeSegment(r)
		if err != nil {
			return fmt.Errorf("read %v-th code segment: %w ", i, err)
		}
	}
	return nil
}

func (m *Module) readSectionData(r io.Reader, ss uint32) error {
	// get the vector size
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return fmt.Errorf("get size of vector: %w", err)
	}

	m.SecData = make([]*DataSegment, vs)
	for i := range m.SecData {
		m.SecData[i], err = readDataSegment(r)
		if err != nil {
			return fmt.Errorf("read %v-th data segment: %w ", i, err)
		}
	}
	return nil
}
