package errors

import "errors"

var (
	ErrAccountNotFound     = errors.New("Unable to find account by account number")
	ErrInsufficientBalance = errors.New("Insufficient balance")

	TextErrorServer                     = "terjadi kesalahan pada server"
	TextErrorInvalidAmount              = "nominal tidak valid"
	TextErrorAccountNotFound            = "rekening tidak ditemukan"
	TextErrorRequestIsNotValid          = "request tidak valid"
	TextErrorInsufficientBalance        = "saldo tidak mencukupi"
	TextErrorNIKOrPhoneNumberRegistered = "nik atau nomor hp sudah terdaftar"
)
