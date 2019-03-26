package iradix

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleNode_Seek() {
	r := New()
	keys := []string{
		"aaa",
		"foo",
		"foo/bar",
		"zzz",
	}
	for _, k := range keys {
		r, _, _ = r.Insert([]byte(k), nil)
	}

	iter := r.Root().Seek([]byte("foo"))
	for {
		k, _, ok := iter.Next()
		if !ok {
			break
		}
		fmt.Printf("%s, ", k)
	}
	// Output: foo/bar, zzz,
}

func TestSeekNext(t *testing.T) {
	r := New()

	keys := []string{
		"foo/bar/baz",
		"foo/baz/bar",
		"foo/zip/zap",
		"foobar",
		"zipzap",
	}
	for _, k := range keys {
		r, _, _ = r.Insert([]byte(k), nil)
	}
	if r.Len() != len(keys) {
		t.Fatalf("bad len: %v %v", r.Len(), len(keys))
	}

	type exp struct {
		inp string
		out []string
	}
	cases := []exp{
		exp{
			"",
			keys,
		},
		exp{
			"f",
			keys[4:],
		},
		exp{
			"foo",
			keys,
		},
		exp{
			"foob",
			keys[4:],
		},
		exp{
			"foo/",
			keys,
		},
		exp{
			"foo/b",
			keys[2:],
		},
		exp{
			"foo/ba",
			keys,
		},
		exp{
			"foo/bar",
			keys[1:],
		},
		exp{
			"foo/bar/baz",
			keys[1:],
		},
		exp{
			"foo/bar/bazoo",
			keys[1:],
		},
		exp{
			"foobar",
			keys[4:],
		},
		exp{
			"z",
			[]string{},
		},
		exp{
			"nosuch",
			[]string{},
		},
	}

	root := r.Root()
	for _, test := range cases {
		iter := root.Seek([]byte(test.inp))

		// Consume all the keys
		out := []string{}
		for {
			key, _, ok := iter.Next()
			if !ok {
				break
			}
			out = append(out, string(key))
		}
		t.Logf("seek %q: got %v", test.inp, out)
		if !reflect.DeepEqual(out, test.out) {
			t.Errorf("mismatch: seek %q, expected %v, got %v", test.inp, test.out, out)
		}
	}
}
