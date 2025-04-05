package testing

import "fmt"

type TestConfigType struct {
	Attr1 string  `json:"attr1"`
	Attr2 int     `json:"attr2"`
	Attr3 bool    `json:"attr3"`
	Attr4 float64 `json:"attr4"`
	Attr5 struct {
		SubAttr1 string `json:"subAttr1"`
		SubAttr2 int    `json:"subAttr2"`
	} `json:"attr5"`
	Attr6 []string `json:"attr6"`
	Port  int
}

func (c TestConfigType) Validate() error {
	if c.Attr1 == "" {
		return fmt.Errorf("Attr1 cannot be empty")
	}
	if c.Attr2 <= 0 {
		return fmt.Errorf("Attr2 must be greater than 0")
	}
	return nil
}

func (c TestConfigType) IsValid() bool {
	return c.Validate() == nil
}

func (c TestConfigType) GetLogLevel() string {
	return "debug"
}

func (c TestConfigType) GetLogFileName() string {
	return ""
}
