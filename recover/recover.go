package recover

import (
	"context"
	"errors"
	"fmt"
	"github.com/reeceappling/goUtils/v2/logging"
	"runtime/debug"
)

var recoveringFromPanicErr = errors.New("recovering from panic")

func HandleRecoverAndLog(ctx context.Context, recoverResult any) (err error) {
	if recoverResult != nil {
		log := logging.GetSugaredLogger(ctx)
		err = recoveringFromPanicErr
		stacktrace := string(debug.Stack())
		log.Errorw("recovering from panic", "message", fmt.Sprint(recoverResult), "stacktrace", stacktrace)
	}
	return
}
