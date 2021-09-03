package consts

const (
	// Endpoint is YooMoney's API endpoint
	//
	// The API supports POST and GET requests.
	// POST requests use JSON arguments, GET requests use query strings.
	// The API always returns the response in JSON format, regardless of the type of request.
	//
	// URL: https://yookassa.ru/en/developers/using-api/basics?lang=bash#interaction-format
	Endpoint = "https://api.yookassa.ru/v3/"

	// IdempotentHeader is used for POST requests
	//
	// Learn more at: https://yookassa.ru/en/developers/using-api/basics#idempotence
	IdempotentHeader = "Idempotence-Key"

	TestingCard3DSecureFailedMastercard = "5555555555554592"
	TestingCard3DSecureFailedVisa       = "4839665499603842"
	TestingCard3DSecureFailedMir        = "2200000000000012"

	TestingCardCallIssuerMastercard = "5555555555554535"
	TestingCardCallIssuerVisa       = "4926946416239025"
	TestingCardCallIssuerMir        = "2200000000000020"

	TestingCardExpiredMastercard = "5555555555554543"
	TestingCardExpiredVisa       = "4141435412630840"
	TestingCardExpiredMir        = "2200000000000038"

	TestingCardFraudSuspectedMastercard = "5555555555554568"
	TestingCardFraudSuspectedVisa       = "4483274282299972"
	TestingCardFraudSuspectedMir        = "2200000000000046"

	TestingCardGeneralDeclineMastercard = "5555555555554527"
	TestingCardGeneralDeclineVisa       = "4889971706588753"
	TestingCardGeneralDeclineMir        = "2202202212312379"

	TestingCardInsufficientFundsMastercard = "5555555555554600"
	TestingCardInsufficientFundsVisa       = "4562265587712390"
	TestingCardInsufficientFundsMir        = "2200000000000053"

	TestingCardInvalidCardNumberMastercard = "5555555555554618"
	TestingCardInvalidCardNumberVisa       = "4951017853630544"
	TestingCardInvalidCardNumberMir        = "2201382000000013"

	TestingCardInvalidCSCMastercard = "5555555555554626"
	TestingCardInvalidCSCVisa       = "4194180666146368"
	TestingCardInvalidCSCMir        = "2200770212727079"

	TestingCardIssuerUnavailableMastercard = "5555555555554501"
	TestingCardIssuerUnavailableVisa       = "4654130848359150"
	TestingCardIssuerUnavailableMir        = "2201382000000021"

	TestingCardPaymentMethodLimitExceededMastercard = "5555555555554576"
	TestingCardPaymentMethodLimitExceededVisa       = "4565231022577548"
	TestingCardPaymentMethodLimitExceededMir        = "2201382000000039"

	TestingCardPaymentMethodRestrictedMastercard = "5555555555554550"
	TestingCardPaymentMethodRestrictedVisa       = "4233961169071671"
	TestingCardPaymentMethodRestrictedMir        = "2201382000000047"

	TestingCardCancelledByYooMoneyCountryForbiddenMastercard = "5555555555554584"
	TestingCardCancelledByYooMoneyCountryForbiddenVisa       = "4969751510013864"
	TestingCardCancelledByYooMoneyCountryForbiddenMir        = "2201382000000054"

	TestingCardCancelledByYooMoneyFraudSuspectedMastercard = "5555555555554634"
	TestingCardCancelledByYooMoneyFraudSuspectedVisa       = "4119098878796485"
	TestingCardCancelledByYooMoneyFraudSuspectedMir        = "2201696981989955"

	TestingCardSuccessfulMastercard3DSecure = "5555555555554477"
	TestingCardSuccessfulMastercard         = "5555555555554444"
	TestingCardSuccessfulMaestro            = "6759649826438453"
	TestingCardSuccessfulVisa3DSecure       = "4793128161644804"
	TestingCardSuccessfulVisa               = "4111111111111111"
	TestingCardSuccessfulVisaElectron       = "4175001000000017"
	TestingCardSuccessfulMir3DSecure        = "2200000000000004"
	TestingCardSuccessfulMir                = "2202474301322987"
	TestingCardSuccessfulAmericanExpress    = "370000000000002"
	TestingCardSuccessfulJCB                = "3528000700000000"
	TestingCardSuccessfulDinersClub         = "36700102000000"
)
