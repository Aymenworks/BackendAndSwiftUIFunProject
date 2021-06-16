package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/entities"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/services/tips"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	requests "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/http/requests/tips"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/infra/s3"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
	"github.com/go-chi/chi"
)

type TipsController struct {
	Controller
	service       tips.Service
	imageUploader *s3.S3ImageUploader
}

func NewTipsController(service tips.Service, imageUploader *s3.S3ImageUploader) *TipsController {
	return &TipsController{
		service:       service,
		imageUploader: imageUploader,
	}
}

func (c *TipsController) GetAll(w http.ResponseWriter, r *http.Request) {
	tips, err := c.service.GetAll(r.Context())
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	c.JsonResponse(w, tips)
}

func (c *TipsController) Get(w http.ResponseWriter, r *http.Request) {
	tipID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	tip, err := c.service.MustGetByID(r.Context(), uint(tipID))
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	c.JsonResponse(w, tip)
}

func (c *TipsController) Create(w http.ResponseWriter, r *http.Request) {
	ct, err := c.ParseContentTypeHeader(r)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	var tip *entities.Tip
	switch ct {
	case "application/json":
		var request requests.CreateTipRequest
		if err := c.ParseBody(r, &request); err != nil {
			c.ErrorResponse(w, err)
			return
		}
		if utils.IsNotEmpty(request.SignedImageURL) {
			tip, err = c.createWithAWSPreSignedURL(r.Context(), request.Name, request.SignedImageURL)
		} else if utils.IsNotEmpty(request.Image) {
			tip, err = c.createWithBase64(r.Context(), request.Name, request.Image)
		} else {
			err = errors.Stack(errors.InvalidParameter)
		}
	case "multipart/form-data":
		tip, err = c.createWithMultipart(r.Context(), r)
	default:
		err = errors.Stack(errors.InvalidContentTypeHeader)
	}

	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	c.JsonResponse(w, tip)
}

func (c *TipsController) createWithMultipart(ctx context.Context, r *http.Request) (*entities.Tip, error) {
	name := r.FormValue("name")
	if utils.IsEmpty(name) {
		return nil, errors.Stack(errors.TipNameInvalid)
	}
	file, handler, err := r.FormFile("myFile")
	defer file.Close()
	if err != nil {
		switch err {
		case http.ErrMissingFile:
			return nil, errors.Wrap(errors.InvalidParameter, fmt.Sprintf("%v", err))
		default:
			return nil, errors.Stack(err)
		}
	}

	ct := handler.Header.Get("content-type")
	if ct != "image/jpeg" {
		return nil, errors.Wrap(errors.InvalidParameter, "img not jpeg")
	}

	path, err := c.imageUploader.UploadWithFile(file)
	if err != nil {
		return nil, errors.Stack(err)
	}
	if utils.IsEmpty(path) {
		return nil, errors.Wrap(errors.UnknownError, "filepath is empty")
	}

	tip, err := c.service.Create(ctx, name, path)
	if err != nil {
		// TODO: delete file
		return nil, errors.Stack(err)
	}

	return tip, nil
}

func (c *TipsController) createWithBase64(ctx context.Context, name, imageB64 string) (*entities.Tip, error) {
	data, err := base64.StdEncoding.DecodeString(imageB64)
	if err != nil {
		return nil, errors.Stack(fmt.Errorf("error  decoding base64 image string err = %w", err))
	}

	ct := http.DetectContentType(data)
	if ct != "image/jpeg" {
		return nil, errors.Wrap(errors.InvalidParameter, "img not jpeg")
	}

	path, err := c.imageUploader.UploadWithBase64(data)
	if err != nil {
		return nil, errors.Stack(err)
	}
	if utils.IsEmpty(path) {
		return nil, errors.Wrap(errors.UnknownError, "filepath is empty")
	}

	tip, err := c.service.Create(ctx, name, path)
	if err != nil {
		// TODO: delete file
		return nil, errors.Stack(err)
	}

	return tip, nil
}

func (c *TipsController) createWithAWSPreSignedURL(ctx context.Context, name, url string) (*entities.Tip, error) {
	// TODO: Temporary bucket with X hours auto delete expiration
	// TODO: Copy from temporary bucket to new
	return nil, nil
}

func (c *TipsController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := c.PathParameterUint(r, "id")
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	err = c.service.DeleteByID(r.Context(), id)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	c.NoContentResponse(w)
}
