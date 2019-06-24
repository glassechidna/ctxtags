package ctxerr

import (
	"context"
	"encoding/json"
	"github.com/glassechidna/ctxtags"
	"github.com/pkg/errors"
)

func Wrap(ctx context.Context, err error) error {
	msgbytes, _ := json.Marshal(ctxtags.Tags(ctx))
	return errors.Wrap(err, "ctx: " + string(msgbytes))
}
