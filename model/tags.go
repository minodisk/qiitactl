package model

import "fmt"

// Tags is slice of Tag.
type Tags []Tag

// MarshalYAML encodes Tags to YAML encoded string.
func (tags Tags) MarshalYAML() (data interface{}, err error) {
	arr := make([]interface{}, len(tags))
	for i, tag := range tags {
		if len(tag.Versions) == 0 {
			arr[i] = tag.Name
		} else {
			obj := make(map[string][]string)
			obj[tag.Name] = tag.Versions
			arr[i] = obj
		}
	}
	data = interface{}(arr)
	return
}

// UnmarshalYAML decodes YAML encoded string to Tag struct.
func (tags *Tags) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var ts []interface{}
	err = unmarshal(&ts)
	if err != nil {
		return
	}

	for _, t := range ts {
		switch t := t.(type) {
		default:
			err = fmt.Errorf("unexpected type in tag: %s", t)
			return
		case string:
			tag := Tag{
				Name: t,
			}
			*tags = append(*tags, tag)
		case map[interface{}]interface{}:
			for n, v := range t {
				name := n.(string)
				vs := v.([]interface{})
				versions := make([]string, len(vs))
				for i, v := range vs {
					versions[i] = fmt.Sprint(v)
				}
				tag := Tag{
					Name:     name,
					Versions: versions,
				}
				*tags = append(*tags, tag)
			}
		}
	}

	return
}
