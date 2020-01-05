package utils

const (
	Success               = 200
	SuccessExisting       = 210
	ErrorParameter        = 400
	ErrorUnLogin          = 401
	ErrorForbidden        = 403
	ErrorNodata           = 404
	ErrorLogic            = 409
	ErrorExpire           = 410
	ErrorParseJson        = 422
	ErrorParseXML         = 423
	ErrorFrequency        = 490 // do something too frequency
	ErrorWebsocketRead    = 4100
	ErrorWebsocketWrite   = 4101
	ErrorService          = 500
	ErrorWebsocketUpgrade = 570 // error
	ErrorCache            = 800
	ErrorDB               = 900
	ErrorNetWork          = 990 // sending http request error
	ErrorHttp             = 991 // sending http request error

	ErrorPackage    = 1000
	ErrorThirdParty = 1100

	ErrorMWLogger = 2000
	ErrorMWEmail  = 2001

	BeegoNoData = "<QuerySeter> no row found"
)
