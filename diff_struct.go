/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package diff

import (
	"reflect"
	"time"
)

func (d *Differ) diffStruct(path []string, pathTypes []interface{}, a, b reflect.Value, parent interface{}) error {
	if AreType(a, b, reflect.TypeOf(time.Time{})) {
		return d.diffTime(path, pathTypes, a, b)
	}

	if a.Kind() == reflect.Invalid {
		if d.DisableStructValues {
			d.cl.Add(CREATE, path, pathTypes, nil, exportInterface(b))
			return nil
		}
		return d.structValues(CREATE, path, pathTypes, b)
	}

	if b.Kind() == reflect.Invalid {
		if d.DisableStructValues {
			d.cl.Add(DELETE, path, pathTypes, exportInterface(a), nil)
			return nil
		}
		return d.structValues(DELETE, path, pathTypes, a)
	}

	for i := 0; i < a.NumField(); i++ {
		field := a.Type().Field(i)
		tname := tagName(d.TagName, field)

		if tname == "-" || hasTagOption(d.TagName, field, "immutable") {
			continue
		}

		if tname == "" {
			tname = field.Name
		}

		af := a.Field(i)
		bf := b.FieldByName(field.Name)

		fpath := path
		tpath := pathTypes
		if !(d.FlattenEmbeddedStructs && field.Anonymous) {
			fpath = copyAppend(fpath, tname)
			tpath = copyAppendType(tpath, exportInterface(bf))
		}

		if d.Filter != nil && !d.Filter(fpath, a.Type(), field) {
			continue
		}

		// skip private fields
		if !a.CanInterface() {
			continue
		}

		err := d.diff(fpath, tpath, af, bf, exportInterface(a))
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Differ) structValues(t string, path []string, pathTypes []interface{}, a reflect.Value) error {
	var nd Differ
	nd.Filter = d.Filter
	nd.customValueDiffers = d.customValueDiffers

	if t != CREATE && t != DELETE {
		return ErrInvalidChangeType
	}

	if a.Kind() == reflect.Ptr {
		a = reflect.Indirect(a)
	}

	if a.Kind() != reflect.Struct {
		return ErrTypeMismatch
	}

	x := reflect.New(a.Type()).Elem()

	for i := 0; i < a.NumField(); i++ {

		field := a.Type().Field(i)
		tname := tagName(d.TagName, field)

		if tname == "-" {
			continue
		}

		if tname == "" {
			tname = field.Name
		}

		af := a.Field(i)
		xf := x.FieldByName(field.Name)

		fpath := copyAppend(path, tname)

		if nd.Filter != nil && !nd.Filter(fpath, a.Type(), field) {
			continue
		}
		tPath := pathTypes
		err := nd.diff(fpath, tPath, xf, af, exportInterface(a))
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(nd.cl); i++ {
		(d.cl) = append(d.cl, swapChange(t, nd.cl[i]))
	}

	return nil
}
