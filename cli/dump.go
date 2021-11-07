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

func Dump(module *types.Module) {
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
		fmt.Printf("  type[%d]: ", i)
		d.dumpFunctionType(ft)
		fmt.Printf("\n")
	}
}

func (d *Dumper) dumpImportSection() {
	fmt.Printf("Import[%d]:\n", len(d.module.SecImport))
	for _, imp := range d.module.SecImport {
		switch imp.Desc.Kind {
		case types.ImportTypeFunc:
			fmt.Printf("  func[%d]: <%s.%s>, sig=%d\n",
				d.importedFuncCount, imp.Module, imp.Name, imp.Desc.TypeIndex)
			d.importedFuncCount++
		case types.ImportTypeTable:
			fmt.Printf("  table[%d]: <%s.%s>, %v\n",
				d.importedTableCount, imp.Module, imp.Name, imp.Desc.TableType.Limit)
			d.importedTableCount++
		case types.ImportTypeMem:
			fmt.Printf("  memory[%d]: <%s.%s>, %v\n",
				d.importedMemCount, imp.Module, imp.Name, imp.Desc.MemType)
			d.importedMemCount++
		case types.ImportTypeGlobal:
			fmt.Printf("  global[%d]: <%s.%s>, %v\n",
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
		fmt.Printf("  table[%d]: ", d.importedTableCount+i)
		d.dumpElemType(t.ElemType)
		fmt.Printf(" ")
		d.dumpLimitType(t.Limit)
		fmt.Printf("\n")
	}
}

func (d *Dumper) dumpMemSection() {
	fmt.Printf("Memory[%d]:\n", len(d.module.SecMemory))
	for i, l := range d.module.SecMemory {
		fmt.Printf("  memory[%d]: pages ", d.importedMemCount+i)
		d.dumpLimitType(l)
		fmt.Printf("\n")
	}
}

func (d *Dumper) dumpGlobalSection() {
	fmt.Printf("Global[%d]:\n", len(d.module.SecGlobal))
	for i, g := range d.module.SecGlobal {
		fmt.Printf("  global[%d]: ", d.importedGlobalCount+i)
		dumpGlobalType(g.Type)
		fmt.Printf(" - ")
		// dumpInitExpression(g.Init)
		fmt.Printf("\n")
	}
}

func dumpGlobalType(gt *types.GlobalType) {
	fmt.Printf("%v", gt.Value.Type)

	fmt.Printf(" ")

	if gt.Mutable {
		fmt.Printf("mutable=true")
	} else {
		fmt.Printf("mutable=false")
	}
}

func (d *Dumper) dumpExportSection() {
	fmt.Printf("Export[%d]:\n", len(d.module.SecExport))
	for _, exp := range d.module.SecExport {
		switch exp.Desc.Kind {
		case types.ExportTypeFunc:
			fmt.Printf("  func[%d]: name=<%s>\n", int(exp.Desc.Index), exp.Name)
		case types.ExportTypeTable:
			fmt.Printf("  table[%d]: name=<%s>\n", int(exp.Desc.Index), exp.Name)
		case types.ExportTypeMem:
			fmt.Printf("  memory[%d]: name=<%s>\n", int(exp.Desc.Index), exp.Name)
		case types.ExportTypeGlobal:
			fmt.Printf("  global[%d]: name=<%s>\n", int(exp.Desc.Index), exp.Name)
		}
	}
}

func (d *Dumper) dumpStartSection() {
	fmt.Printf("Start:\n")
	if d.module.SecStart != nil {
		fmt.Printf("  func=%d\n", d.module.SecStart.(uint32))
	} else {
		fmt.Printf("  No start function.\n")
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
	for i, _ := range d.module.SecCode {
		fmt.Printf("  func[%d]:\n", d.importedFuncCount+i)
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
	fmt.Printf("  name=%s\n", d.module.SecCustom.Name)
	fmt.Printf("  %s\n", string(d.module.SecCustom.Bytes))
}

func (d *Dumper) dumpElemType(et byte) {
	if et == types.ElemTypeFuncRef {
		fmt.Printf("type=funcref")
	}
}

func (d *Dumper) dumpLimitType(limit *types.LimitType) {
	if limit.Tag == types.LimitTypeOnlyMin {
		fmt.Printf("initial=%v", limit.Min)
	} else if limit.Tag == types.LimitTypeBothMinAndMax {
		fmt.Printf("initial=%v max=%v", limit.Min, limit.Max)
	}
}

func (d *Dumper) dumpFunctionType(ft *types.FunctionType) {
	d.dumpValueTypes(ft.InputType)
	fmt.Printf(" -> ")
	d.dumpValueTypes(ft.ReturnType)
}

func (d *Dumper) dumpValueTypes(vt []types.ValueType) {
	// output
	if len(vt) == 0 {
		fmt.Printf("nil")
	} else {
		for i, op := range vt {
			fmt.Printf("%v", op.Type)
			if i != len(vt)-1 {
				fmt.Printf(", ")
			}
		}
	}
}
