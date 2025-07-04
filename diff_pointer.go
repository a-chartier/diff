/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package diff

import (
	"reflect"
	"unsafe"
)

var isExportFlag uintptr = (1 << 5) | (1 << 6)

func (d *Differ) diffPtr(path []string, pathTypes []interface{}, a, b reflect.Value, parent interface{}) error {
	if a.Kind() != b.Kind() {
		if a.Kind() == reflect.Invalid {
			if !b.IsNil() {
				return d.diff(path, pathTypes, reflect.ValueOf(nil), reflect.Indirect(b), parent)
			}

			d.cl.Add(CREATE, path, pathTypes, nil, exportInterface(b), parent)
			return nil
		}

		if b.Kind() == reflect.Invalid {
			if !a.IsNil() {
				return d.diff(path, pathTypes, reflect.Indirect(a), reflect.ValueOf(nil), parent)
			}

			d.cl.Add(DELETE, path, pathTypes, exportInterface(a), nil, parent)
			return nil
		}

		return ErrTypeMismatch
	}

	if a.IsNil() && b.IsNil() {
		return nil
	}

	if a.IsNil() {
		d.cl.Add(UPDATE, path, pathTypes, nil, exportInterface(b), parent)
		return nil
	}

	if b.IsNil() {
		d.cl.Add(UPDATE, path, pathTypes, exportInterface(a), nil, parent)
		return nil
	}

	return d.diff(path, pathTypes, reflect.Indirect(a), reflect.Indirect(b), parent)
}

func exportInterface(v reflect.Value) interface{} {
	if !v.CanInterface() {
		flagTmp := (*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + 2*unsafe.Sizeof(uintptr(0))))
		*flagTmp = (*flagTmp) & (^isExportFlag)
	}
	return v.Interface()
}
