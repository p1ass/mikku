package mikku

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	semVerPrefix = "v"
)

const (
	majorIdx = iota
	minorIdx
	patchIdx
)

type bumpType int

const (
	major bumpType = iota + 1
	minor
	patch
	version
)

var semVerReg = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)`)

// determineNewTag bump version if the given version is major, minor, or patch
// otherwise, use the given version without change
func determineNewTag(version string, currentTag string) (string, error) {
	if version == "major" || version == "minor" || version == "patch" {
		if !validSemver(currentTag) {
			return "", errInvalidSemanticVersioningTag
		}
		bumpTyp := strToBumpType(version)
		newTag, err := bumpVersion(currentTag, bumpTyp)
		if err != nil {
			return "", fmt.Errorf("bump version: %w", err)
		}
		return newTag, nil
	}

	if !validSemver(version) {
		return "", errInvalidSemanticVersioningTag
	}

	return version, nil
}

func validSemver(ver string) bool {
	return semVerReg.Match([]byte(ver))
}

func strToBumpType(str string) bumpType {
	switch str {
	case "major":
		return major
	case "minor":
		return minor
	case "patch":
		return patch
	default:
		return 0
	}
}

func bumpVersion(tag string, typ bumpType) (string, error) {
	tag = strings.TrimPrefix(tag, semVerPrefix)
	splitTag := strings.Split(tag, ".")

	versions, err := strsToInts(splitTag)
	if err != nil {
		return "", fmt.Errorf("strsToInts: %w", err)
	}

	switch typ {
	case major:
		versions[majorIdx]++
		versions[minorIdx] = 0
		versions[patchIdx] = 0
	case minor:
		versions[minorIdx]++
		versions[patchIdx] = 0
	case patch:
		versions[patchIdx]++
	default:
		return "", fmt.Errorf("invalid bump type: %w", err)
	}
	return createSemanticVersion(versions), nil
}

func strsToInts(strs []string) ([]int, error) {
	ints := make([]int, len(strs))

	for idx, s := range strs {
		converted, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("convert string to int: %w", err)
		}
		ints[idx] = converted
	}
	return ints, nil
}

func createSemanticVersion(versions []int) string {
	sm := semVerPrefix

	sm += strconv.Itoa(versions[0])

	for _, v := range versions[1:] {
		sm += "."
		sm += strconv.Itoa(v)
	}

	return sm
}
