package models

const SettingKeyName = "name"

type Setting struct {
	Name  string `json:"name" bson:"_id"`
	Value any    `json:"value"`
}

func (s *Setting) Validate() []error {
	return []error{}
}

func (s *Setting) GetKey() string {
	return s.Name
}
