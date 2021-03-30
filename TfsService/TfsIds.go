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
func (i *TfsIds) JoinToString(sep string) string {
	if len(i.Ids) == 0 {
		return ""
	}

	b := make([]string, len(i.Ids))
	for i, v := range i.Ids {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
