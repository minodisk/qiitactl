package model

type Tags []Tag

func (tags Tags) MarshalYAML() (data interface{}, err error) {
	obj := make(map[string][]string)
	for _, tag := range tags {
		obj[tag.Name] = tag.Versions
	}
	data = interface{}(obj)
	return
}

func (tags *Tags) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var t map[string][]string
	err = unmarshal(&t)
	if err != nil {
		return
	}

	for name, versions := range t {
		tag := Tag{
			Name:     name,
			Versions: versions,
		}
		*tags = append(*tags, tag)
	}

	return
}
