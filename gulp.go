package gulp

import (
	"fmt"
	"github.com/chriso345/gulp/constants"
)

func Gulp() bool {
	fmt.Println(constants.LpConstraintLE)
	return constants.Runner()
}
