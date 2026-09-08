package core

import (
	"encoding/json"
	"reflect"
	"testing"
)

// A campaign is a save file. Every field added to the world or to the player
// has to survive being written down and read back, and the one time that failed
// it was a json tag collision that only `go vet` noticed — five fields sharing
// two names. This checks the whole shape rather than the fields somebody
// remembered to test.

// distinctive fills every field of a struct with a value that is not the zero
// value, so a field that fails to round-trip is a field that comes back zero.
func distinctive(v reflect.Value, seed *int) {
	*seed++
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(*seed + 7))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(*seed + 11))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(float64(*seed) + .5)
	case reflect.String:
		v.SetString("value-" + itoa(*seed))
	case reflect.Slice:
		item := reflect.New(v.Type().Elem()).Elem()
		distinctive(item, seed)
		v.Set(reflect.Append(reflect.MakeSlice(v.Type(), 0, 1), item))
	case reflect.Map:
		key := reflect.New(v.Type().Key()).Elem()
		distinctive(key, seed)
		val := reflect.New(v.Type().Elem()).Elem()
		distinctive(val, seed)
		m := reflect.MakeMap(v.Type())
		m.SetMapIndex(key, val)
		v.Set(m)
	case reflect.Ptr:
		p := reflect.New(v.Type().Elem())
		distinctive(p.Elem(), seed)
		v.Set(p)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).PkgPath != "" {
				continue // unexported
			}
			distinctive(v.Field(i), seed)
		}
	}
}

// omitted lists fields that deliberately do not survive a save, with the reason.
var omitted = map[string]string{
	"VisualCues": "presentation for one response, never stored",
}

func TestEveryFieldOfTheWorldSurvivesBeingWrittenDown(t *testing.T) {
	var before World
	seed := 0
	distinctive(reflect.ValueOf(&before).Elem(), &seed)

	data, err := json.Marshal(&before)
	if err != nil {
		t.Fatal(err)
	}
	var after World
	if err := json.Unmarshal(data, &after); err != nil {
		t.Fatal(err)
	}

	kind := reflect.TypeOf(before)
	for i := 0; i < kind.NumField(); i++ {
		field := kind.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if why, ok := omitted[field.Name]; ok {
			if !reflect.ValueOf(&after).Elem().Field(i).IsZero() {
				t.Fatalf("%s was stored after all, though it is meant to be %s", field.Name, why)
			}
			continue
		}
		got := reflect.ValueOf(&after).Elem().Field(i).Interface()
		want := reflect.ValueOf(&before).Elem().Field(i).Interface()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("World.%s did not survive: wrote %#v, read %#v", field.Name, want, got)
		}
	}
}

func TestEveryFieldOfThePlayerSurvivesBeingWrittenDown(t *testing.T) {
	var before Person
	seed := 100
	distinctive(reflect.ValueOf(&before).Elem(), &seed)

	data, err := json.Marshal(&before)
	if err != nil {
		t.Fatal(err)
	}
	var after Person
	if err := json.Unmarshal(data, &after); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) {
		kind := reflect.TypeOf(before)
		for i := 0; i < kind.NumField(); i++ {
			got := reflect.ValueOf(&after).Elem().Field(i).Interface()
			want := reflect.ValueOf(&before).Elem().Field(i).Interface()
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Person.%s did not survive: wrote %#v, read %#v", kind.Field(i).Name, want, got)
			}
		}
	}
}

// And the same for every other thing the world stores a list of, because a
// collision inside one of those is just as fatal and just as invisible.
func TestEveryStoredRecordSurvivesBeingWrittenDown(t *testing.T) {
	for _, sample := range []any{
		&Property{}, &Faction{}, &NPC{}, &Crew{}, &Plot{}, &Task{}, &Record{},
		&Death{}, &Story{}, &Good{}, &Contract{}, &Commission{}, &Grudge{},
		&Pact{}, &Conflict{}, &TableHand{}, &Choice{}, &Effect{}, &Scene{},
		&Offer{}, &Director{}, &Result{}, &ArrangementMemory{}, &SuspendedJob{},
	} {
		value := reflect.ValueOf(sample).Elem()
		seed := 1
		distinctive(value, &seed)
		data, err := json.Marshal(sample)
		if err != nil {
			t.Fatalf("%T: %v", sample, err)
		}
		back := reflect.New(value.Type())
		if err := json.Unmarshal(data, back.Interface()); err != nil {
			t.Fatalf("%T: %v", sample, err)
		}
		kind := value.Type()
		for i := 0; i < kind.NumField(); i++ {
			field := kind.Field(i)
			if field.PkgPath != "" || field.Tag.Get("json") == "-" {
				continue
			}
			got := back.Elem().Field(i).Interface()
			want := value.Field(i).Interface()
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("%s.%s did not survive: wrote %#v, read %#v", kind.Name(), field.Name, want, got)
			}
		}
	}
}
