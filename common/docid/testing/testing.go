package testing

import (
	. "github.com/pigeond-io/pigeond/common/docid"
	"math/rand"
	"strconv"
	"testing"
)

var (
	ShouldContainEdge = "Should contain %v -> %v"
	ShouldNotContainEdge = "Should not contain %v -> %v"
	CountShouldBeXButIsY = "Count should be %d but is %d"
	ShouldContain = "Should contain %v"
	ShouldNotContain = "Should not contain %v"
)

func ID(id int64) *IntId {
	return &IntId{Id: id}
}

func SID(id string) *StrId {
	return &StrId{Id: id}
}

func RandStrIds(count int) []DocId {
	max := count >> 1
	slice := make([]DocId, 0, count)
	for i := 0; i < count; i++ {
		slice = append(slice, SID(strconv.Itoa(rand.Intn(max))))
	}
	return slice
}

/* Default TestTuple */
type TestTuple struct {
	t *testing.T
}

func (tt *TestTuple) CountShouldBe(countExpected, countActual int) (* TestTuple) {
	if countExpected != countActual {
		tt.t.Errorf(CountShouldBeXButIsY, countExpected, countActual)
	}
	return tt
}

/* Set */
type TestTupleSet struct {
	TestTuple
	s Set
}

func TestSet(t *testing.T, s Set) (* TestTupleSet){
	tt := &TestTupleSet{s:s}
	tt.t = t
	return tt
}

func (tt *TestTupleSet) ShouldContain(a DocId) (* TestTupleSet){
	if !tt.s.Contains(a) {
		tt.t.Errorf(ShouldContain, a)
	}
	return tt
}

func (tt *TestTupleSet) ShouldNotContain(a DocId) (* TestTupleSet){
	if tt.s.Contains(a) {
		tt.t.Errorf(ShouldNotContain, a)
	}
	return tt
}

/* EdgeSet */
type TestTupleEdgeSet struct {
	TestTuple
	e EdgeSet
}

func TestEdgeSet(t *testing.T, e EdgeSet) (* TestTupleEdgeSet){
	tt := &TestTupleEdgeSet{e:e}
	tt.t = t
	return tt
}

func (tt *TestTupleEdgeSet) ShouldContain(a DocId, b DocId) (* TestTupleEdgeSet){
	if !tt.e.Contains(a, b) {
		tt.t.Errorf(ShouldContainEdge, a, b)
	}
	return tt
}

func (tt *TestTupleEdgeSet) ShouldNotContain(a DocId, b DocId) (* TestTupleEdgeSet){
	if tt.e.Contains(a, b) {
		tt.t.Errorf(ShouldNotContainEdge, a, b)
	}
	return tt
}