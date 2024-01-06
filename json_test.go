package magicjson_test

import (
	"encoding/json"
	"fmt"
	"github.com/parinpan/magicjson"
	"reflect"
	"testing"
	"time"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("lat:%f,lon:%f", l.Latitude, l.Longitude)
	return []byte(s), nil
}

type student struct {
	name     string
	age      uint8
	height   float32
	joinedAt *time.Time
}

type teacher struct {
	name     string
	age      uint8
	height   float32
	joinedAt *time.Time
}

type classRoom struct {
	name            string
	capacity        uint
	homeRoomTeacher teacher
	students        []student
}

type school struct {
	name            string
	establishedDate *time.Time
	classRooms      []classRoom

	Location Location
}

func TestMarshalJSON(t *testing.T) {
	t.Run("should generate null when data is nil", func(t *testing.T) {
		expected := []byte(`null`)
		actual, err := magicjson.MarshalJSON(nil)

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal on data type", func(t *testing.T) {
		expected, _ := json.Marshal([]byte(`1234`))
		actual, err := magicjson.MarshalJSON([]byte(`1234`))

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal on primitive slice", func(t *testing.T) {
		expected := []byte(`[1,2,3,4]`)
		actual, err := magicjson.MarshalJSON([]int{1, 2, 3, 4})

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal on primitive array", func(t *testing.T) {
		expected := []byte(`[1,2,3,4]`)
		actual, err := magicjson.MarshalJSON([4]int{1, 2, 3, 4})

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal given data", func(t *testing.T) {
		joinedAt := toPtr(time.Unix(1704537759, 0))
		establishedAt := toPtr(time.Unix(1504537759, 0))

		schools := []school{
			{
				name:            "School A",
				establishedDate: establishedAt,
				classRooms: []classRoom{
					{
						name:     "Class A",
						capacity: 100,
						homeRoomTeacher: teacher{
							name:     "Mrs. Diana",
							age:      30,
							height:   180.5,
							joinedAt: joinedAt,
						},
						students: []student{
							{
								name:     "Brian",
								age:      14,
								height:   170.5,
								joinedAt: joinedAt,
							},
						},
					},
				},
				Location: Location{
					Latitude:  1,
					Longitude: 2,
				},
			},
		}

		expected := []byte(`[{"name":"School A","establishedDate":"2017-09-04T17:09:19+02:00","classRooms":[{"name":"Class A","capacity":100,"homeRoomTeacher":{"name":"Mrs. Diana","age":30,"height":180.5,"joinedAt":"2024-01-06T11:42:39+01:00"},"students":[{"name":"Brian","age":14,"height":170.5,"joinedAt":"2024-01-06T11:42:39+01:00"}]}],"Location":lat:1.000000,lon:2.000000}]`)
		actual, err := magicjson.MarshalJSON(schools)

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})
}

func toPtr[T any](v T) *T {
	return &v
}
