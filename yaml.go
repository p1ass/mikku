package mikku

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-yaml/yaml"
)

var (
	errImageNotFoundInYAML = errors.New("image not found in yaml file")
	errInvalidYAML         = errors.New("invalid yaml file")
)

func getCurrentTag(yamlStr, imageName string) (string, error) {
	var obj interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &obj); err != nil {
		return "", fmt.Errorf("yaml unmarshal: %w", err)
	}

	if v, ok := obj.(map[interface{}]interface{}); ok {
		tag := traverse(v, imageName)
		if tag == "" {
			return "", fmt.Errorf("%s: %w", imageName, errImageNotFoundInYAML)
		}
		return tag, nil
	}
	return "", errInvalidYAML
}

func traverse(node map[interface{}]interface{}, imageName string) string {
	for key, value := range node {
		switch val := value.(type) {
		case string:
			if strings.Contains(val, imageName) && key == "image" {
				split := strings.Split(val, ":")
				return split[len(split)-1]
			}
		case map[interface{}]interface{}:
			if tag := traverse(val, imageName); tag != "" {
				return tag
			}
		case []interface{}:
			for _, v := range val {
				if convert, ok := v.(map[interface{}]interface{}); ok {
					if tag := traverse(convert, imageName); tag != "" {
						return tag
					}
				}
			}
		}
	}
	return ""
}
