package TfsService

import (
	"strconv"
	"strings"
)

type TfsIds struct {
	Ids []int
}

// Test ids contains element
func (i *TfsIds) Contains(element int) bool {
	for _, id := range i.Ids {
		if id == element {
			return true
		}
	}
	return false
}

// Add distinct targets ids
func (i *TfsIds) AddTargets(relations TfsRelations) {
	for _, relation := range relations.WorkItemRelations {
		if !i.Contains(relation.Target.Id) {
			i.Ids = append(i.Ids, relation.Target.Id)
		}
	}
}

// Convert to string with separator
func (i *TfsIds) JoinToString(separator string) string {
	if len(i.Ids) == 0 {
		return ""
	}

	ids := make([]string, len(i.Ids))
	for i, id := range i.Ids {
		ids[i] = strconv.Itoa(id)
	}
	return strings.Join(ids, separator)
}

// Take Part
func (i *TfsIds) TakePart(size int, skip int) TfsIds {
	if skip > len(i.Ids) {
		skip = len(i.Ids)
	}
	end := skip + size
	if end > len(i.Ids) {
		end = len(i.Ids)
	}
	return TfsIds{Ids: i.Ids[skip:end]}
}
