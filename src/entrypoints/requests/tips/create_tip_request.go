package requests

import (
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
)

type AppRequest interface {
	Validate() error
}

type CreateTipRequest struct {
	Name           string `json:"name"`
	Image          string `json:"image"`
	SignedImageURL string `json:"signed_image_url"`
}

func (r *CreateTipRequest) Validate() error {
	if utils.IsEmpty(r.Name) {
		return errors.TipNameInvalid
	}
	if utils.IsEmpty(r.SignedImageURL) && utils.IsEmpty(r.Image) {
		return errors.InvalidParameter
	}
	return nil
}
