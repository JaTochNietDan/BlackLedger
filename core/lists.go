package core

import "reflect"

// A list the core has never put anything in is not a list at all by the time it
// reaches the view. Go writes a nil slice as `null`, and `omitempty` leaves the
// key out entirely, so `cards.board.length` is a blank screen rather than an
// empty table — which is exactly how a hand of cards blanked the screen.
//
// Nothing is gained by the saving. An empty list costs two characters and says
// the true thing: there is nothing here yet. So no list in the core omits
// itself any more, and this walks whatever is about to be sent and gives every
// nil one a body before it goes.
//
// It is called where the world becomes a payload rather than at every place a
// slice could be left alone, because there are hundreds of those and one of
// this.

// FillLists replaces every nil slice reachable from the world with an empty one.
func (w *World) FillLists() { fillLists(reflect.ValueOf(w)) }

func fillLists(v reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			fillLists(v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() {
				fillLists(v.Field(i))
			}
		}
	case reflect.Slice:
		if v.IsNil() {
			if v.CanSet() {
				v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			}
			return
		}
		for i := 0; i < v.Len(); i++ {
			fillLists(v.Index(i))
		}
	case reflect.Map:
		// A map nothing has been put in has no entries to walk and cannot be
		// written to at all — a save from before a table existed has several.
		if v.IsNil() {
			return
		}
		// A map's values cannot be set in place, so each one is taken out,
		// filled and put back.
		for _, key := range v.MapKeys() {
			held := v.MapIndex(key)
			if held.Kind() == reflect.Pointer || held.Kind() == reflect.Interface {
				fillLists(held)
				continue
			}
			copied := reflect.New(held.Type()).Elem()
			copied.Set(held)
			fillLists(copied)
			v.SetMapIndex(key, copied)
		}
	}
}
