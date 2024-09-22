package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {

	var jsonValue string = strconv.Quote(fmt.Sprintf("%d mins", r))
	
	return []byte(jsonValue), nil

}

