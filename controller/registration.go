package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"account_service/dto"
	"account_service/service"

	errInternal "account_service/errors"
)

type registrationController struct {
	registrationService service.RegistrationService
}

func NewRegistrationController(registrationService service.RegistrationService) *registrationController {
	return &registrationController{
		registrationService: registrationService,
	}
}

func (c *registrationController) Register(ctx echo.Context) error {
	var registrationRequest dto.RegistrationRequest
	if err := ctx.Bind(&registrationRequest); err != nil {
		log.Errorf("Error when binding request: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.RegistrationResponse{
				Remark: errInternal.TextErrorRequestIsNotValid,
			},
		)
	}

	if err := ctx.Validate(registrationRequest); err != nil {
		log.Errorf("Error when validating registration request: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.RegistrationResponse{
				Remark: errInternal.TextErrorRequestIsNotValid,
			},
		)
	}

	isNIKOrPhoneNumberRegistered, err := c.registrationService.IsNIKOrPhoneNumberRegistered(
		registrationRequest.NIK,
		registrationRequest.PhoneNumber,
	)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			dto.RegistrationResponse{
				Remark: errInternal.TextErrorServer,
			},
		)
	}

	if isNIKOrPhoneNumberRegistered {
		log.Warnf("NIK or phone number %s is already registered", registrationRequest.NIK)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.RegistrationResponse{
				Remark: errInternal.TextErrorNIKOrPhoneNumberRegistered,
			},
		)
	}

	accountNumber, err := c.registrationService.Register(
		registrationRequest.Name,
		registrationRequest.NIK,
		registrationRequest.PhoneNumber,
	)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			dto.RegistrationResponse{
				Remark: errInternal.TextErrorServer,
			},
		)
	}

	log.Infof("Registration successful registered: %s", accountNumber)
	return ctx.JSON(
		http.StatusOK,
		dto.RegistrationResponse{
			AccountNumber: accountNumber,
		},
	)
}
