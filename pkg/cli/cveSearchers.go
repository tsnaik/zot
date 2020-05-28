package cli

import (
	"errors"
	"strings"
)

func getSearchers() []searcher {
	searchers := []searcher{
		new(cveByImageNameSearcher),
		new(imageByCveIDSearcher),
	}

	return searchers
}

type searcher interface {
	search(params map[string]*string, searchService CveSearchService, servURL, user *string) (string, error)
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

func (search imageByCveIDSearcher) search(params map[string]*string,
	searchService CveSearchService, servURL, user *string) (string, error) {
	if !canSearch(params, newSet("cveID")) {
		return "", ErrCannotSearch
	}

	username, password := getUsernameAndPassword(*user)
	results, err := searchService.findImagesByCveID(*params["cveID"], *servURL, username, password)

	if err != nil {
		return "", err
	}

	return results.String(), nil
}

type cveByImageNameSearcher struct{}

func (search cveByImageNameSearcher) search(params map[string]*string,
	searchService CveSearchService, servURL, user *string) (string, error) {
	if !canSearch(params, newSet("imageName")) {
		return "", ErrCannotSearch
	}

	username, password := getUsernameAndPassword(*user)
	results, err := searchService.findCveByImageName(*params["imageName"], *servURL, username, password)

	if err != nil {
		return "", err
	}

	return results.String(), nil
}

func getUsernameAndPassword(user string) (string, string) {
	if strings.Contains(user, ":") {
		split := strings.Split(user, ":")
		return split[0], split[1]
	}

	return "", ""
}

func getEmptyStruct() struct{} {
	return struct{}{}
}

type set struct {
	m map[string]struct{}
}

func newSet(initialValues ...string) *set {
	s := &set{}
	s.m = make(map[string]struct{})

	for _, val := range initialValues {
		s.m[val] = getEmptyStruct()
	}

	return s
}

func (s *set) contains(value string) bool {
	_, c := s.m[value]
	return c
}

var (
	ErrCannotSearch = errors.New("cannot search with these parameters")
)
