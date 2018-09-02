package blackfriday

import "testing"

func TestEmtyAttributes(t *testing.T) {
	a := NewAttributes()
	r := a.String()
	e := ""
	if r != e {
		t.Errorf("Emmpty attributes must return empty string\nExpected: %s\nActual: %s", e, r)
	}
}

func TestAddOneAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper")
	if s := a.String(); s != "class=\"wrapper\"" {
		t.Errorf("Unexprected output: %s", s)
	}
}

func TestAddFewValuesToOneAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper").Add("class", "-with-image")
	if s := a.String(); s != "class=\"wrapper -with-image\"" {
		t.Errorf("Unexpected output: %s", s)
	}
}

func TestRemoveValueFromOneAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper").Add("class", "-with-image")
	if s := a.String(); s != "class=\"wrapper -with-image\"" {
		t.Errorf("Unexpected output: %s", s)
	}
	a.RemoveValue("class", "wrapper")
	if s := a.String(); s != "class=\"-with-image\"" {
		t.Errorf("Unexpected output: %s", s)
	}
}

func TestRemoveWholeAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper")
	if s := a.String(); s != "class=\"wrapper\"" {
		t.Errorf("Unexprected output: %s", s)
	}
	a.Remove("class")
	if a.String() != "" {
		t.Errorf("Emmpty attributes must return empty string")
	}
}

func TestRemoveWholeAttributeByValue(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper")
	if s := a.String(); s != "class=\"wrapper\"" {
		t.Errorf("Unexprected output: %s", s)
	}
	a.RemoveValue("class", "wrapper")
	r := a.String()
	e := ""
	if r != e {
		t.Errorf("Emmpty attributes must return empty string\nExpected: %s\nActual: %s", e, r)
	}
}

func TestAddFewAttributes(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper").Add("id", "main-block")
	if s := a.String(); s != "class=\"wrapper\" id=\"main-block\"" {
		t.Errorf("Unexprected output: %s", s)
	}
}

func TestAddComplexAttributes(t *testing.T) {
	a := NewAttributes()
	a.
		Add("style", "background: #fff;").
		Add("style", "font-size: 14px;").
		Add("data-test-id", "block")
	e := "style=\"background: #fff; font-size: 14px;\" data-test-id=\"block\""
	if s := a.String(); s != e {
		t.Errorf("Unexpected output: %s", s)
	}
}
