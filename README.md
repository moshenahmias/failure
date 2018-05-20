# **failure**

**failure** is an error handling package for [Go](https://golang.org/).



![](failure.png)



With **failure** you can construct fielded errors:

```go
err := failure.Build("something went wrong").
		WithField("id", 5).
		WithField("severity", "fatal").
		Done()

// you can test for fields value the followng way:
b0 := failure.TestField(err, "id", 5) // b0 == true
b1 := failure.TestField(err, "severity", "normal") // b1 == false
b2 := failure.TestField(err, "message", "something went wrong") // b2 == true
b3 := failure.TestField(err, "e", "mc^2") // b3 == false

// there's also a fields getter:
v0, e0 := failure.Field(err, "severity") // v0 == "fatal", e0 == nil
v1 := failure.FieldOrDefault(err, "e", "mc^2") // v1 == "mc^2"
```



**failure** is "compatible" with the *[errors.New](https://golang.org/pkg/errors/#example_New)* and *[fmt.Errorf](https://golang.org/pkg/errors/#example_New_errorf)* functions signature:

```go
err0 := failure.New("something went wrong")
err1 := failure.Errorf("something went %s wrong", "terribly") // or failure.Newf
```



**failure** errors string representation are JSONs which are easy to read and parse:

```go
err := failure.Buildf("something went %s wrong", "terribly").
		WithField("id", 5).
		WithField("severity", "fatal").
		Done()

fmt.Println(err) // {"message":"something went terribly wrong","fields":{"id":5,"severity":"fatal"}}
```



Use **failure** to construct recursive errors:

```go
err3 := failure.Build("something went wrong").
		WithField("level", 3).
		Done()

err2 := failure.Build("something went wrong").
		WithField("level", 2).
		ParentOf(err3).
		Done()

err1 := failure.Build("something went wrong").
		WithField("level", 1).
		ParentOf(err2).
		Done()

err0 := failure.Build("something went wrong").
		WithField("level", 0).
		ParentOf(err1).
		Done()

// you can find the error's immediate descendant:
inner := failure.Inner(err0) // inner == err1

// or the error's origin error:
origin := failure.Origin(err0) // origin == err3

// you can verify a parent-descendant relationship between two errors:
b0 := failure.IsParentOf(err0, err1) // b0 == true
b1 := failure.IsParentOf(err0, err2) // b1 == true
b2 := failure.IsParentOf(err0, err3) // b2 == true
b3 := failure.IsParentOf(err0, err0) // b3 == false
b4 := failure.IsParentOf(err1, err0) // b4 == false

// there's a recursive TestField version:
b5 := failure.TestFieldRecursively(err0, "level", 3) // b5 == true
b6 := failure.TestFieldRecursively(err0, "level", 4) // b6 == false
```



You can enrich existing errors with fields and inner error:

```go
if _, err := ioutil.ReadFile("/dev/null"); err != nil {

    // use Buildc if you need to add fields or inner error
    // to an existing error
    return failure.Buildc(err).
    		WithField("id", 3).
    		ParentOf(errors.New("something went wrong")).
    		Done()
}
```



Errors comparison with **failure** is simple:

```go
err0 := failure.New("something went wrong")
err1 := failure.New("something went wrong")
err2 := errors.New("something went wrong")

err3 := failure.Build("something went wrong").
		WithField("id", 2).
		Done()

err4 := failure.Build("something went wrong").
		WithField("id", 2).
		ParentOf(err0).
		Done()

err5 := failure.New("something went terribly wrong")

b0 := err0 == err1             // b0 == false
b1 := err0 == err2             // b1 == false
b2 := err0 == err3             // b2 == false
b3 := err0 == err0             // b3 == true

// Same compares the message and fields for the error and every descendant:
b4 := failure.Same(err0, err1) // b4 == true
b5 := failure.Same(err0, err2) // b5 == true
b6 := failure.Same(err0, err3) // b6 == false
b7 := failure.Same(err3, err4) // b7 == false

// Like compares only the message of both errors (without comparing the descendants):
b8 := failure.Like(err0, err1) // b8 == true
b9 := failure.Like(err0, err3) // b9 == true
b10 := failure.Like(err0, err5) // b10 == false
```



Using Like + Buildc is very common with package-level errors:

```go
package mypkg

var (
	ErrInvalidValue = failure.New("mypkg: invalid value")
)

func isValid(val string) bool {
    ...
}

func Foo(val string) error {
    
    if isValid(val) {
        return nil
    }
    
    return failure.Buildc(ErrInvalidValue).
    				WithField("value", val).
    				Done()
}

package main

func main() {

    err := mypkg.Foo("all your base are belong to us");
    
    // don't do:
    if err == mypkg.ErrInvalidValue {
        ...
    }
    
    // do:
    if failure.Like(err, mypkg.ErrInvalidValue) {
        ...
    }
}
```



**failure** can work with any type that implements the *error* interface:

```go
err1 := fmt.Errorf("something went %s", "wrong")

err0 := failure.Build("something went terribly wrong").
		ParentOf(err1).
		Done()
```



When working with external errors (created without **failure**), the error's `Error() string` result will be used as the *message* field for a newly created error with no other fields or inner error.

For a more "accurate" conversion, implement the *failure.Impersonator* interface:

```go
type extErr string

func (e *extErr) Error() string {
	return string(*e)
}

func (e *extErr) Impersonate(b failure.Builder) {
    b.WithField(failure.MessageField, "everything is wrong")
    b.WithField("id", 1)
}

func main() {

	ext := extErr("something went terribly wrong")

	err := failure.Build("something went wrong").
    			ParentOf(&ext).
    			Done()

    msg := failure.Message(failure.Origin(err)) // msg == "everything is wrong"
}
```

