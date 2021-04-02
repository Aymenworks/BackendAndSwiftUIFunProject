package requests

import (
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
)

type AppRequest interface {
	Validate() error
}

type CreateTipRequest struct {
	Name string `json:"name"`
}

func (r *CreateTipRequest) Validate() error {
	if utils.IsEmpty(r.Name) {
		return errors.TipNameInvalid.New()
	}
	return nil
}
