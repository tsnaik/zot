package cli

import (
	"errors"
)

func getSearchers() []searcher {
	searchers := []searcher{
		new(cveByImageNameSearcher),
		new(imageByCveIDSearcher),
	}

	return searchers
}

var allowedCombinations = `
Only these combinations of flags(or their shorthands) are allowed:
  --image-name

URL of the zot repository with is required (--url)
`

type searcher interface {
	search(params map[string]*string, searchService CveSearchService) (string, error)
}

func canSearch(params map[string]*string, requiredParams *set) bool {
	for key, value := range params {
		if requiredParams.contains(key) && *value == "" {
			return false
		} else if !requiredParams.contains(key) && *value != "" {
			return false
		}
	}
	return true
}

type imageByCveIDSearcher struct{}

func (search imageByCveIDSearcher) search(params map[string]*string, searchService CveSearchService) (string, error) {
	if !canSearch(params, newSet("cveIDForImage")) {
		return "", cannotSearchError
	}
	if results, err := searchService.findImagesByCveId(*params["cveIDForImage"], servURL); err != nil {
		return "", err
	} else {
		return results.String(), nil
	}
}

type cveByImageNameSearcher struct{}

func (search cveByImageNameSearcher) search(params map[string]*string, searchService CveSearchService) (string, error) {
	if !canSearch(params, newSet("imageName")) {
		return "", cannotSearchError
	}
	if results, err := searchService.findCveByImageName(*params["imageName"], servURL); err != nil {
		return "", err
	} else {
		return results.String(), nil
	}
}

var exists = struct{}{}

type set struct {
	m map[string]struct{}
}

func newSet(initialValues ...string) *set {
	s := &set{}
	s.m = make(map[string]struct{})
	for _, val := range initialValues {
		s.m[val] = exists
	}
	return s
}

func (s *set) contains(value string) bool {
	_, c := s.m[value]
	return c
}

var cannotSearchError = errors.New("Cannot search with these parameters")
