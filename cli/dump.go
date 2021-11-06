package cli

import (
	"fmt"
	"github.com/LBruyne/wasm-decode/types"
)

type Dumper struct {
	module              *types.Module
	importedFuncCount   int
	importedTableCount  int
	importedMemCount    int
	importedGlobalCount int
}

func dump(module *types.Module) {
	d := &Dumper{module: module}

	fmt.Printf("Version: 0x%02x\n", d.module.Version)
	d.dumpTypeSection()
	d.dumpImportSection()
	d.dumpFuncSection()
	d.dumpTableSection()
	d.dumpMemSection()
	d.dumpGlobalSection()
	d.dumpExportSection()
	d.dumpStartSection()
	d.dumpElemSection()
	d.dumpCodeSection()
	d.dumpDataSection()
	d.dumpCustomSection()
}

func (d *Dumper) dumpTypeSection() {
	fmt.Printf("Type[%d]:\n", len(d.module.SecType))
	for i, ft := range d.module.SecType {
		fmt.Printf("  type[%d]: %s\n", i, ft)
	}
}

func (d *Dumper) dumpImportSection() {
	fmt.Printf("Import[%d]:\n", len(d.module.SecImport))
	for _, imp := range d.module.SecImport {
		switch imp.Desc.Kind {
		case types.ImportTypeFunc:
			fmt.Printf("  func[%d]: %s.%s, sig=%d\n",
				d.importedFuncCount, imp.Module, imp.Name, imp.Desc.TypeIndex)
			d.importedFuncCount++
		case types.ImportTypeTable:
			fmt.Printf("  table[%d]: %s.%s, %s\n",
				d.importedTableCount, imp.Module, imp.Name, imp.Desc.TableType.Limit)
			d.importedTableCount++
		case types.ImportTypeMem:
			fmt.Printf("  memory[%d]: %s.%s, %s\n",
				d.importedMemCount, imp.Module, imp.Name, imp.Desc.MemType)
			d.importedMemCount++
		case types.ImportTypeGlobal:
			fmt.Printf("  global[%d]: %s.%s, %s\n",
				d.importedGlobalCount, imp.Module, imp.Name, imp.Desc.GlobalType)
			d.importedGlobalCount++
		}
	}
	return
}

func (d *Dumper) dumpFuncSection() {
	fmt.Printf("Function[%d]:\n", len(d.module.SecFunction))
	for i, sig := range d.module.SecFunction {
		fmt.Printf("  func[%d]: sig=%d\n",
			d.importedFuncCount+i, sig)
	}
}

func (d *Dumper) dumpTableSection() {
	fmt.Printf("Table[%d]:\n", len(d.module.SecTable))
	for i, t := range d.module.SecTable {
		fmt.Printf("  table[%d]: %s\n",
			d.importedTableCount+i, t.Limit)
	}
}

func (d *Dumper) dumpMemSection() {
	fmt.Printf("Memory[%d]:\n", len(d.module.SecMemory))
	for i, limits := range d.module.SecMemory {
		fmt.Printf("  memory[%d]: %s\n",
			d.importedMemCount+i, limits)
	}
}

func (d *Dumper) dumpGlobalSection() {
	fmt.Printf("Global[%d]:\n", len(d.module.SecGlobal))
	for i, g := range d.module.SecGlobal {
		fmt.Printf("  global[%d]: %s\n",
			d.importedGlobalCount+i, g.Type)
	}
}

func (d *Dumper) dumpExportSection() {
	fmt.Printf("Export[%d]:\n", len(d.module.SecExport))
	for _, exp := range d.module.SecExport {
		switch exp.Desc.Kind {
		case types.ExportTypeFunc:
			fmt.Printf("  func[%d]: name=%s\n", int(exp.Desc.Index), exp.Name)
		case types.ExportTypeTable:
			fmt.Printf("  table[%d]: name=%s\n", int(exp.Desc.Index), exp.Name)
		case types.ExportTypeMem:
			fmt.Printf("  memory[%d]: name=%s\n", int(exp.Desc.Index), exp.Name)
		case types.ExportTypeGlobal:
			fmt.Printf("  global[%d]: name=%s\n", int(exp.Desc.Index), exp.Name)
		}
	}
}

func (d *Dumper) dumpStartSection() {
	fmt.Printf("Start:\n")
	if d.module.SecStart != 0 {
		fmt.Printf("  func=%d\n", d.module.SecStart.(uint32))
	}
}

func (d *Dumper) dumpElemSection() {
	fmt.Printf("Element[%d]:\n", len(d.module.SecElement))
	for i, elem := range d.module.SecElement {
		fmt.Printf("  elem[%d]: table=%d\n", i, elem.TableIdx)
	}
}

func (d *Dumper) dumpCodeSection() {
	fmt.Printf("Code[%d]:\n", len(d.module.SecCode))
	for i, code := range d.module.SecCode {
		fmt.Printf("  func[%d]: locals=[", d.importedFuncCount+i)
		if len(code.Locals) > 0 {
			for i, locals := range code.Locals {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s x %d",
					locals.Type, locals.Count)
			}
		}
		fmt.Println("]")
	}
}

func (d *Dumper) dumpDataSection() {
	fmt.Printf("Data[%d]:\n", len(d.module.SecData))
	for i, data := range d.module.SecData {
		fmt.Printf("  data[%d]: mem=%d\n", i, data.MemIdx)
	}
}

func (d *Dumper) dumpCustomSection() {
	fmt.Printf("Custom:\n")
	fmt.Printf("  custom: name=%s\n", d.module.SecCustom.Name)
	fmt.Printf("          byte=%s\n", string(d.module.SecCustom.Bytes))
}
