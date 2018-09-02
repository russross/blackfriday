package blackfriday

import "testing"

func TestEmtyAttributes(t *testing.T) {
	a := NewAttributes()
	r := a.String()
	e := ""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}

func TestAddOneAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper")
	r := a.String()
	e := "class=\"wrapper\""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}

func TestAddFewValuesToOneAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper").Add("class", "-with-image")
	r := a.String()
	e := "class=\"wrapper -with-image\""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}

func TestRemoveValueFromOneAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper").Add("class", "-with-image")
	a.RemoveValue("class", "wrapper")
	r := a.String()
	e := "class=\"-with-image\""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}

func TestRemoveWholeAttribute(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper")
	a.Remove("class")
	r := a.String()
	e := ""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}

func TestRemoveWholeAttributeByValue(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper")
	a.RemoveValue("class", "wrapper")
	r := a.String()
	e := ""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}

func TestAddFewAttributes(t *testing.T) {
	a := NewAttributes()
	a.Add("class", "wrapper").Add("id", "main-block")
	r := a.String()
	e := "class=\"wrapper\" id=\"main-block\""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}

func TestAddComplexAttributes(t *testing.T) {
	a := NewAttributes()
	a.
		Add("style", "background: #fff;").
		Add("style", "font-size: 14px;").
		Add("data-test-id", "block")
	r := a.String()
	e := "style=\"background: #fff; font-size: 14px;\" data-test-id=\"block\""
	if r != e {
		t.Errorf("Expected: %s\nActual: %s\n", e, r)
	}
}
