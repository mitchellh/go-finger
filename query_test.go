package finger

import (
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	cases := []struct {
		Input    string
		Expected *Query
	}{
		{
			"",
			&Query{Username: ""},
		},

		{
			"foo",
			&Query{Username: "foo"},
		},

		{
			"foo@host",
			&Query{
				Username: "foo",
				Hostname: []string{"host"},
			},
		},

		{
			"foo@jump@host",
			&Query{
				Username: "foo",
				Hostname: []string{"jump", "host"},
			},
		},

		{
			"@host",
			&Query{
				Username: "",
				Hostname: []string{"host"},
			},
		},

		{
			"foo@host\n",
			&Query{
				Username: "foo",
				Hostname: []string{"host"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			actual, err := ParseQuery(tc.Input)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
