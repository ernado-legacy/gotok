package gotok

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestGotok(t *testing.T) {
	Convey("Gotok", t, func() {
		sesson, err := mgo.Dial("localhost")
		So(err, ShouldBeNil)

		c := sesson.DB("test").C("tokens")

		Reset(func() {
			c.DropCollection()
		})

		storage := New(c)

		t, err := storage.Generate(bson.NewObjectId())
		So(err, ShouldBeNil)

		t2, err := storage.Get(t.Token)
		So(t.Id, ShouldEqual, t2.Id)

		Convey("From database", func() {
			s := New(c)

			t4, err := s.Get(t.Token)
			So(err, ShouldBeNil)
			So(t4.Id, ShouldEqual, t.Id)
		})

		Convey("Remove", func() {
			err = storage.Remove(t)
			So(err, ShouldBeNil)

			t3, err := storage.Get(t.Token)
			So(err, ShouldBeNil)
			So(t3, ShouldBeNil)
		})
	})
}

func BenchmarkGet(b *testing.B) {
	sesson, err := mgo.Dial("localhost")
	if err != nil {
		b.Fatal(err)
	}
	c := sesson.DB("test").C("tokens")
	storage := New(c)
	t, err := storage.Generate(bson.NewObjectId())
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		storage.Get(t.Token)
	}

	c.DropCollection()
}

func BenchmarkGenerate(b *testing.B) {
	sesson, err := mgo.Dial("localhost")
	if err != nil {
		b.Fatal(err)
	}
	c := sesson.DB("test").C("tokens")
	storage := New(c)

	for n := 0; n < b.N; n++ {
		_, err = storage.Generate(bson.NewObjectId())
		if err != nil {
			b.Fatal(err)
		}
	}

	c.DropCollection()
}
