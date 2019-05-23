package zend

/*
#include <zend_API.h>
#include <zend_ini.h>

#ifdef ZEND_ENGINE_2
typedef zend_ini_entry zend_ini_entry_def;
typedef char zend_string;  // 为了让gozend_ini_modifier7在php5下能够编译过
#endif

extern int gozend_ini_modifier7(zend_ini_entry *entry, zend_string *new_value, void *mh_arg1, void *mh_arg2, void *mh_arg3, int stage);
extern int gozend_ini_modifier5(zend_ini_entry *entry, char *new_value, uint new_value_length, void *mh_arg1, void *mh_arg2, void *mh_arg3, int stage);
extern void gozend_ini_displayer(zend_ini_entry *ini_entry, int type);
*/
import "C"
import "unsafe"

import (
	"fmt"
	// "reflect"
	"log"
	"runtime"
)

// 配置实体
type IniEntries struct {
	zies []C.zend_ini_entry_def
}

var zendIniEntryDefZero C.zend_ini_entry_def

func NewIniEntries() *IniEntries {
	this := &IniEntries{}
	this.zies = make([]C.zend_ini_entry_def, 1)
	this.zies[0] = zendIniEntryDefZero
	return this
}

func (ini *IniEntries) Register(moduleNumber int) int {
	r := C.zend_register_ini_entries(&ini.zies[0], C.int(moduleNumber))
	log.Println(r)
	return int(r)
}
func (ini *IniEntryDef) Unregister(moduleNumber int) {
	C.zend_unregister_ini_entries(C.int(moduleNumber))
}

func (ini *IniEntries) Add(ie *IniEntryDef) {
	ini.zies[len(ini.zies)-1] = ie.zie
	ini.zies = append(ini.zies, zendIniEntryDefZero)
}

type IniEntry struct {
	zie *C.zend_ini_entry
}

func newZendIniEntryFrom(ie *C.zend_ini_entry) *IniEntry {
	return &IniEntry{ie}
}
func (this *IniEntry) Name() string      { return fromZString(this.zie.name) }
func (this *IniEntry) Value() string     { return fromZString(this.zie.value) }
func (this *IniEntry) OrigValue() string { return fromZString(this.zie.orig_value) }

const (
	INI_USER   = int(C.ZEND_INI_USER)
	INI_PERDIR = int(C.ZEND_INI_PERDIR)
	INI_SYSTEM = int(C.ZEND_INI_SYSTEM)
)

type IniEntryDef struct {
	zie C.zend_ini_entry_def

	onModify  func(ie *IniEntry, newValue string, stage int) int
	onDisplay func(ie *IniEntry, itype int)
}

func NewIniEntryDef() *IniEntryDef {
	this := &IniEntryDef{}
	// this.zie = (*C.zend_ini_entry_def)(C.calloc(1, C.sizeof_zend_ini_entry_def))
	runtime.SetFinalizer(this, zendIniEntryDefFree)
	return this
}

func zendIniEntryDefFree(this *IniEntryDef) {
	if _, ok := iniNameEntries[C.GoString(this.zie.name)]; ok {
		delete(iniNameEntries, C.GoString(this.zie.name))
	}

	if this.zie.name != nil {
		C.free(unsafe.Pointer(this.zie.name))
	}
	if this.zie.value != nil {
		C.free(unsafe.Pointer(this.zie.value))
	}
}

func (ini *IniEntryDef) Fill3(name string, defaultValue interface{}, modifiable bool,
	onModify func(), arg1, arg2, arg3 interface{}, displayer func()) {
	ini.zie.name = C.CString(name)
	ini.zie.modifiable = C.uchar(go2cBool(modifiable)) //php7.3 cannot use go2cBool(modifiable) (type _Ctype_int) as type _Ctype_uchar in assignment

	// ini.zie.orig_modifiable = go2cBool(modifiable)
	if ZEND_ENGINE == ZEND_ENGINE_3 {
		ini.zie.on_modify = go2cfn(C.gozend_ini_modifier7)
	} else {
		ini.zie.on_modify = go2cfn(C.gozend_ini_modifier5)
	}
	ini.zie.displayer = go2cfn(C.gozend_ini_displayer)

	value := fmt.Sprintf("%v", defaultValue)
	ini.zie.value = C.CString(value)

	if arg1 == nil {
		ini.zie.mh_arg1 = nil
	}
	if arg2 == nil {
		ini.zie.mh_arg2 = nil
	}
	if arg3 == nil {
		ini.zie.mh_arg3 = nil
	}

	if ZEND_ENGINE == ZEND_ENGINE_3 {
		ini.zie.name_length = C.ushort(len(name))
		ini.zie.value_length = C.uint32_t(len(value))
	} else {
		// why need +1 for php5?
		// if not, zend_alter_ini_entry_ex:280行会出现zend_hash_find无结果失败
		ini.zie.name_length = C.ushort(len(name) + 1)
		ini.zie.value_length = C.uint32_t(len(value) + 1)
	}
	log.Println(name, len(name))

	iniNameEntries[name] = ini
}

func (ini *IniEntryDef) Fill2(name string, defaultValue interface{}, modifiable bool,
	onModify func(), arg1, arg2 interface{}, displayer func()) {
	ini.Fill3(name, defaultValue, modifiable, onModify, arg1, arg2, nil, displayer)
}

func (ini *IniEntryDef) Fill1(name string, defaultValue interface{}, modifiable bool,
	onModify func(), arg1 interface{}, displayer func()) {
	ini.Fill3(name, defaultValue, modifiable, onModify, arg1, nil, nil, displayer)
}

func (ini *IniEntryDef) Fill(name string, defaultValue interface{}, modifiable bool,
	onModify func(), displayer func()) {
	ini.Fill3(name, defaultValue, modifiable, onModify, nil, nil, nil, displayer)
}

func (ini *IniEntryDef) SetModifier(modifier func(ie *IniEntry, newValue string, state int) int) {
	ini.onModify = modifier
}

func (ini *IniEntryDef) SetDisplayer(displayer func(ie *IniEntry, itype int)) {
	ini.onDisplay = displayer
}

var iniNameEntries = make(map[string]*IniEntryDef, 0)

// the new_value is really not *C.char, it's *C.zend_string
//export gozend_ini_modifier7
func gozend_ini_modifier7(ie *C.zend_ini_entry, new_value *C.zend_string, mh_arg1 unsafe.Pointer, mh_arg2 unsafe.Pointer, mh_arg3 unsafe.Pointer, stage C.int) C.int {
	// log.Println(ie, "//", new_value, stage, ie.modifiable)
	// log.Println(ie.orig_modifiable, ie.modified, fromZString(ie.orig_value))
	// log.Println(fromZString(new_value), fromZString(ie.name), fromZString(ie.value))
	if iedef, ok := iniNameEntries[fromZString(ie.name)]; ok {
		iedef.onModify(newZendIniEntryFrom(ie), fromZString(new_value), int(stage))
	} else {
		log.Println("wtf", fromZString(ie.name))
	}
	return 0
}

//export gozend_ini_modifier5
func gozend_ini_modifier5(ie *C.zend_ini_entry, new_value *C.char, new_value_length C.uint, mh_arg1 unsafe.Pointer, mh_arg2 unsafe.Pointer, mh_arg3 unsafe.Pointer, stage C.int) C.int {
	// log.Println(ie, "//", new_value, new_value_length, stage, ie.modifiable)
	// log.Println(ie.orig_modifiable, ie.modified, fromZString(ie.orig_value))
	// log.Println(fromZString(new_value), fromZString(ie.name), fromZString(ie.value))
	if iedef, ok := iniNameEntries[fromZString(ie.name)]; ok {
		iedef.onModify(newZendIniEntryFrom(ie), fromZString(new_value), int(stage))
	} else {
		log.Println("wtf", fromZString(ie.name))
	}
	return 0
}

//export gozend_ini_displayer
func gozend_ini_displayer(ie *C.zend_ini_entry, itype C.int) {
	log.Println(ie, itype)
	if iedef, ok := iniNameEntries[fromZString(ie.name)]; ok {
		iedef.onDisplay(newZendIniEntryFrom(ie), int(itype))
	} else {
		log.Println("wtf", fromZString(ie.name))
	}
}
