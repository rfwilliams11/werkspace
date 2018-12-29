package mongo

import (
	"testing"
	"reflect"
	"github.com/globalsign/mgo/bson"
)

func TestUnwind_OneLayer(t *testing.T) {
	field := "test"

	result := Unwind(field)
	expected := bson.M{"$unwind":"$test"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestUnwind_MultipleLayers(t *testing.T) {
	field := "test.nested"

	result := Unwind(field)
	expected := bson.M{"$unwind":"$test"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestUnwind_Empty(t *testing.T) {
	field := ""

	result := Unwind(field)
	expected := bson.M{"$unwind":"$"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestDistinct(t *testing.T) {
	field := "test"

	result := Distinct(field)
	expected := bson.M{"$group":bson.M{"_id":"$test"}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestProjection_Single(t *testing.T) {
	fields := []string{"test"}

	result := Projection(fields)
	expected := bson.M{"$project":bson.M{"test":1}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestProjection_Multiple(t *testing.T) {
	fields := []string{"test","abc","def"}

	result := Projection(fields)
	expected := bson.M{"$project":bson.M{"test":1,"abc":1,"def":1}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestProjection_Empty(t *testing.T) {
	fields := []string{}

	result := Projection(fields)
	expected := bson.M{"$project":bson.M{}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestMatch_TwoFields(t *testing.T) {
	fields := map[string][]string{"field":[]string{"value"}, "aaa":[]string{"bbb",""}}
	result := Match(fields)
	expected := bson.M{"$match":
		bson.M{"$and":[]bson.M{
			bson.M{"$or":[]bson.M{bson.M{"field":"value"}}},
			bson.M{"$or":[]bson.M{bson.M{"aaa":"bbb"}}}}}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestMatch_TwoFieldsMultiple(t *testing.T) {
	fields := map[string][]string{"field":[]string{"value"}, "aaa":[]string{"bbb,111,222"}}
	result := Match(fields)
	expected := bson.M{"$match":
	bson.M{"$and":[]bson.M{
		bson.M{"$or":[]bson.M{bson.M{"field":"value"}}},
		bson.M{"$or":[]bson.M{
			bson.M{"aaa":"bbb"},
			bson.M{"aaa":"111"},
			bson.M{"aaa":"222"},
			}}}}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestMatch_Empty(t *testing.T) {
	fields := map[string][]string{}
	result := Match(fields)

	if result != nil{
		t.Errorf("Unexpected result -\n\texpected nothing\n\tgot : %v", result)
	}
}


func TestMatch_EmptyKey(t *testing.T) {
	fields := map[string][]string{"field":[]string{"value"}, "aaa":[]string{}}
	result := Match(fields)
	expected := bson.M{"$match":
	bson.M{"$and":[]bson.M{
		bson.M{"$or":[]bson.M{bson.M{"field":"value"}}}}}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}


func TestQuery(t *testing.T) {
	fields := map[string][]string{"field":[]string{"value"}}
	match := Match(fields)
	project := Projection([]string{"aaa"})
	unwind := Unwind("distinct")
	distinct := Distinct("distinct")

	expected := []bson.M{match,project,unwind,distinct}
	result := Query(match, []string{"aaa"}, "distinct")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestQuery_NoDistinct(t *testing.T) {
	fields := map[string][]string{"field":[]string{"value"}}
	match := Match(fields)
	project := Projection([]string{"aaa"})

	expected := []bson.M{match,project}
	result := Query(match, []string{"aaa"}, "")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}

func TestQuery_NoProjection(t *testing.T) {
	fields := map[string][]string{"field":[]string{"value"}}
	match := Match(fields)
	unwind := Unwind("distinct")
	distinct := Distinct("distinct")

	expected := []bson.M{match,unwind,distinct}
	result := Query(match, nil, "distinct")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result -\n\texpected: %v\n\tgot : %v", expected, result)
	}
}