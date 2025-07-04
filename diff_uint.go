/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package diff

import (
	"reflect"
)

func (d *Differ) diffUint(path []string, pathTypes []interface{}, a, b reflect.Value, parent interface{}) error {
	if a.Kind() == reflect.Invalid {
		d.cl.Add(CREATE, path, pathTypes, nil, exportInterface(b))
		return nil
	}

	if b.Kind() == reflect.Invalid {
		d.cl.Add(DELETE, path, pathTypes, exportInterface(a), nil)
		return nil
	}

	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	if a.Uint() != b.Uint() {
		if a.CanInterface() {
			d.cl.Add(UPDATE, path, pathTypes, exportInterface(a), exportInterface(b), parent)
		} else {
			d.cl.Add(UPDATE, path, pathTypes, a.Uint(), b.Uint(), parent)
		}
	}

	return nil
}
