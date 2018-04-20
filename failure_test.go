package failure

import (
	"errors"
	"testing"
)

func Test001(t *testing.T) {

	err := New("something went wrong")

	if err == nil {
		t.Fatal("err == nil")
	}

	err2 := err

	if err != err2 {
		t.Fatal("err != err2")
	}

	err3 := New("something went wrong")

	if err == err3 || err2 == err3 {
		t.Fatal("err == err3 || err2 == err3")
	}
}

func Test002(t *testing.T) {

	err := Newf("%s went %s", "something", "wrong")

	if err == nil {
		t.Fatal("err == nil")
	}

	err2 := err

	if err != err2 {
		t.Fatal("err != err2")
	}

	err3 := Errorf("%s went %s", "something", "wrong")

	if err == err3 || err2 == err3 {
		t.Fatal("err == err3 || err2 == err3")
	}
}

func Test003(t *testing.T) {

	err := Build("something went wrong").WithField("id", 5).Done()

	if !TestField(err, "id", 5) {
		t.Fatal("!TestField(err, \"id\", 5)")
	}

	if TestField(err, "id", "5") {
		t.Fatal("TestField(err, \"id\", \"5\")")
	}

	id, e := Field(err, "id")

	if e != nil {
		t.Fatal(e)
	}

	if id != interface{}(5) {
		t.Fatal("id != interface{}(5)")
	}
}

func Test004(t *testing.T) {

	err := Build("something went wrong").
		WithField("a", 1).
		WithField("b", "2").
		Done()

	if !TestField(err, "a", 1) {
		t.Fatal("!TestField(err, \"a\", 1)")
	}

	if !TestField(err, "b", "2") {
		t.Fatal("!TestField(err, \"b\", \"2\")")
	}

	if TestField(err, "c", nil) {
		t.Fatal("!TestField(err, \"c\", nil)")
	}

	if TestField(err, "", nil) {
		t.Fatal("!TestField(err, \"\", nil)")
	}
}

func Test005(t *testing.T) {

	err := Build("something went wrong").
		WithFields(Fields{"a": 1, "b": "2"}).
		Done()

	if !TestField(err, "a", 1) {
		t.Fatal("!TestField(err, \"a\", 1)")
	}

	if !TestField(err, "b", "2") {
		t.Fatal("!TestField(err, \"b\", \"2\")")
	}

	if TestField(err, "c", nil) {
		t.Fatal("!TestField(err, \"c\", nil)")
	}

	if TestField(err, "", nil) {
		t.Fatal("!TestField(err, \"\", nil)")
	}
}

func Test006(t *testing.T) {

	err := New("something went wrong")

	parent := Build("parent of something went wrong").
		ParentOf(err).
		Done()

	if !IsParentOf(parent, err) {
		t.Fatal("!Contains(parent, err)")
	}

	grandParent := Build("parent of parent something went wrong").
		ParentOf(parent).
		Done()

	if !IsParentOf(grandParent, parent) {
		t.Fatal("!Contains(grandParent, parent)")
	}

	if !IsParentOf(grandParent, err) {
		t.Fatal("!Contains(grandParent, parent)")
	}

	if IsParentOf(parent, grandParent) {
		t.Fatal("Contains(parent, grandParent)")
	}

	if IsParentOf(err, parent) {
		t.Fatal("Contains(err, parent)")
	}

	if IsParentOf(err, grandParent) {
		t.Fatal("Contains(err, grandParent)")
	}
}

func Test007(t *testing.T) {

	if IsParentOf(nil, nil) {
		t.Fatal("IsParentOf(nil, nil)")
	}

	err := Build("something went wrong").
		ParentOf(nil).Done()

	if !IsParentOf(err, nil) {
		t.Fatal("IsParentOf(err, nil)")
	}

	if IsParentOf(nil, err) {
		t.Fatal("IsParentOf(nil, err)")
	}
}

func Test008(t *testing.T) {

	err := New("something went terribly wrong")

	parent := Build("something went wrong").
		ParentOf(err).
		ParentOf(nil).
		Done()

	if IsParentOf(parent, err) {
		t.Fatal("IsParentOf(parent, err)")
	}

	if !IsParentOf(parent, nil) {
		t.Fatal("!IsParentOf(parent, nil)")
	}
}

func Test009(t *testing.T) {

	err3 := Buildf("something went wrong %d", 3).
		WithField("id", 3).
		Done()

	err2 := Buildf("something went wrong %d", 2).
		WithField("id", 2).
		ParentOf(err3).
		Done()

	err1 := Buildf("something went wrong %d", 1).
		WithField("id", 1).
		ParentOf(err2).
		Done()

	err0 := Buildf("something went wrong %d", 0).
		WithField("id", 0).
		ParentOf(err1).
		Done()

	for err, i := err0, 0; err != nil; i++ {

		if !TestField(err, "id", i) {
			t.Fatalf("!TestField(err, \"id\", %d)", i)
		}

		err = Inner(err)
	}

	if Origin(err0) != err3 {
		t.Fatal("Origin(err0) != err3")
	}
}

func Test010(t *testing.T) {

	err := New("something went wrong")

	if Message(err) != "something went wrong" {
		t.Fatalf("Message(err) != \"something went wrong\"")
	}

	err = Build("something went wrong").
		WithField("message", "something went terribly wrong").
		Done()

	if Message(err) != "something went terribly wrong" {
		t.Fatalf("Message(err) != \"something went terribly wrong\"")
	}

	err = errors.New("something went wrong")

	if Message(err) != "something went wrong" {
		t.Fatalf("Message(err) != \"something went wrong\"")
	}
}

func Test011(t *testing.T) {

	err0 := errors.New("something went wrong")
	err1 := errors.New("something went wrong")
	err2 := New("something went wrong")

	err3 := Build("something went wrong").
		ParentOf(New("something went terribly wrong")).
		Done()

	err4 := Build("something went wrong").
		ParentOf(Newf("something went %s wrong", "terribly")).
		Done()

	err5 := Build("something went wrong").
		ParentOf(New("many things went wrong")).
		Done()

	err6 := New("everything is wrong")

	err7 := Build("something went wrong").
		WithField("id", 10).
		ParentOf(Newf("something went %s wrong", "terribly")).
		Done()

	err8 := Build("something went wrong").
		WithField("id", 20).
		ParentOf(Newf("something went %s wrong", "terribly")).
		Done()

	err9 := Build("something went wrong").
		ParentOf(Newf("something went %s wrong", "terribly")).
		WithFields(Fields{"id": 20}).
		Done()

	if !Like(err0, err1) {
		t.Fatalf("!Like(err0, err1)")
	}

	if !Like(err0, err2) {
		t.Fatalf("!Like(err0, err2)")
	}

	if !Like(err0, err3) {
		t.Fatalf("!Like(err0, err3)")
	}

	if Like(err0, err6) {
		t.Fatalf("Like(err0, err6)")
	}

	if !Like(err0, err7) {
		t.Fatalf("!Like(err0, err7)")
	}

	if !Same(err0, err1) {
		t.Fatalf("!Same(err0, err1)")
	}

	if !Same(err0, err2) {
		t.Fatalf("!Same(err0, err2)")
	}

	if !Same(err3, err4) {
		t.Fatalf("!Same(err3, err4)")
	}

	if !Same(err0, err0) {
		t.Fatalf("!Same(err0, err0)")
	}

	if !Same(err5, err5) {
		t.Fatalf("!Same(err5, err5)")
	}

	if !Same(err8, err9) {
		t.Fatalf("!Same(err8, err9)")
	}

	if Same(err0, err6) {
		t.Fatalf("Same(err0, err6)")
	}

	if Same(err3, err5) {
		t.Fatalf("Same(err3, err5)")
	}

	if Same(err4, err7) {
		t.Fatalf("Same(err4, err7)")
	}

	if Same(err7, err8) {
		t.Fatalf("Same(err7, err8)")
	}
}

func Test012(t *testing.T) {

	err2 := Build("something went wrong 2").
		WithField("id", 1).
		WithField("type", "fatal").
		Done()

	err1 := Build("something went wrong 1").
		WithField("type", "normal").
		ParentOf(err2).
		Done()

	err0 := Build("something went wrong 0").
		WithField("id", 5).
		ParentOf(err1).
		Done()

	if !TestFieldRecursively(err0, "id", 5) {
		t.Fatal("!TestFieldRecursively(err0, \"id\", 5)")
	}

	if !TestFieldRecursively(err0, "id", 1) {
		t.Fatal("!TestFieldRecursively(err0, \"id\", 1)")
	}

	if !TestFieldRecursively(err0, "type", "normal") {
		t.Fatal("!TestFieldRecursively(err0, \"id\", \"normal\")")
	}

	if TestFieldRecursively(err0, "id", 2) {
		t.Fatal("TestFieldRecursively(err0, \"id\", 2)")
	}

	if TestFieldRecursively(err0, "typo", "fatal") {
		t.Fatal("TestFieldRecursively(err0, \"typo\", \"fatal\")")
	}
}

func Test013(t *testing.T) {
	err := Build("something went wrong").
		WithField("id", 1).
		Done()

	if FieldOrDefault(err, "message", "") != "something went wrong" {
		t.Fatal("FieldOrDefault(err, \"message\", \"\") != \"something went wrong\"")
	}

	if FieldOrDefault(err, "id", 0) != 1 {
		t.Fatal("FieldOrDefault(err, \"id\", 0) != 1")
	}

	if FieldOrDefault(err, "type", "normal") != "normal" {
		t.Fatal("FieldOrDefault(err, \"type\", \"normal\") != \"normal\"")
	}
}

func Test014(t *testing.T) {

	err := errors.New("something went wrong")

	errc := Buildc(err).
		Done()

	if !Same(err, errc) {
		t.Fatal("!Same(err, errc)")
	}

	err = errors.New("something went wrong")

	errc = Buildc(err).
		WithField("id", 1).
		Done()

	if Same(err, errc) {
		t.Fatal("Same(err, errc)")
	}

	err2 := Build("something went wrong 2").
		WithField("id", 1).
		WithField("type", "fatal").
		Done()

	err1 := Build("something went wrong 1").
		WithField("type", "normal").
		ParentOf(err2).
		Done()

	err0 := Build("something went wrong 0").
		WithField("id", 5).
		ParentOf(err1).
		Done()

	errc = Buildc(err0).
		Done()

	if !Same(err0, errc) {
		t.Fatal("!Same(err0, errc)")
	}
}

type extErr string

func (e *extErr) Error() string {
	return string(*e)
}

func (e *extErr) Impersonate(b Builder) {
	b.WithField(MessageField, "everything is wrong")
}

func Test015(t *testing.T) {

	ext := extErr("something went terribly wrong")

	err := Build("something went wrong").
		ParentOf(&ext).
		Done()

	if msg := Message(Origin(err)); msg != "everything is wrong" {
		t.Fatalf("%s != everything is wrong", msg)
	}
}

func Test016(t *testing.T) {

	err2 := Build("something went wrong 2").
		WithField("id", 1).
		WithField("type", "fatal").
		Done()

	err1 := Build("something went wrong 1").
		WithField("type", "normal").
		ParentOf(err2).
		Done()

	err0 := Build("something went wrong 0").
		WithField("id", 5).
		ParentOf(err1).
		Done()

	if Depth(err0) != 3 {
		t.Fatal("Depth(err0) != 3")
	}
}
