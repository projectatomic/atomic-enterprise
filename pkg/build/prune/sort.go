package prune

import (
	buildapi "github.com/projectatomic/atomic-enterprise/pkg/build/api"
)

// sortableBuilds supports sorting Build items by most recently created Build
type sortableBuilds []*buildapi.Build

func (s sortableBuilds) Len() int {
	return len(s)
}

func (s sortableBuilds) Less(i, j int) bool {
	return !s[i].CreationTimestamp.Before(s[j].CreationTimestamp)
}

func (s sortableBuilds) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}
