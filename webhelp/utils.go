package webhelp

import (
	"math/rand"
	"strconv"
)

// BuildRandomNumber is set at build time
var BuildRandomNumber string = strconv.Itoa(rand.Intn(1000000))
