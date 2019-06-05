package util

import (
	"fmt"

	"github.com/pdropaiva/poc-diff-csv/domain"
)

// SplitDiff ...
func SplitDiff(diff map[string]*domain.ExportDiff) (add []domain.UserAudience, remove []domain.UserAudience) {
	for _, u := range diff {
		if !u.IsOld && u.IsNew {
			add = append(add, u.Data)
		}

		if u.IsOld && !u.IsNew {
			remove = append(remove, u.Data)
		}
	}
	return add, remove
}

// PrintDiff ...
func PrintDiff(add, remove []domain.UserAudience) {
	fmt.Println("************* Count add *************")
	fmt.Println(len(add))
	fmt.Println("************* Array add *************")
	fmt.Println(add)
	fmt.Println("************ Count remove ***********")
	fmt.Println(len(remove))
	fmt.Println("************ Array remove ***********")
	fmt.Println(remove)
}
