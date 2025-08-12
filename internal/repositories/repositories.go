package stores

import (
	"time"

	"github.com/tuantran1810/go-di-template/libs/logger"
)

var log = logger.MustNamedLogger("repositories")

const defaultTimeout = 20 * time.Second
