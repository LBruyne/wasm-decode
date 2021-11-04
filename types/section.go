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
		err = m.readSectionMemorie(r, ss)
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
