// Code generated by "stringer -type=FeatureFlag"; DO NOT EDIT.

package features

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[unused-0]
	_ = x[PerformValidationRPC-1]
	_ = x[ACME13KeyRollover-2]
	_ = x[SimplifiedVAHTTP-3]
	_ = x[TLSSNIRevalidation-4]
	_ = x[AllowRenewalFirstRL-5]
	_ = x[SetIssuedNamesRenewalBit-6]
	_ = x[FasterRateLimit-7]
	_ = x[ProbeCTLogs-8]
	_ = x[RevokeAtRA-9]
	_ = x[CAAValidationMethods-10]
	_ = x[CAAAccountURI-11]
	_ = x[HeadNonceStatusOK-12]
	_ = x[NewAuthorizationSchema-13]
	_ = x[DisableAuthz2Orders-14]
	_ = x[EarlyOrderRateLimit-15]
	_ = x[EnforceMultiVA-16]
	_ = x[MultiVAFullResults-17]
	_ = x[RemoveWFE2AccountID-18]
	_ = x[CheckRenewalFirst-19]
	_ = x[MandatoryPOSTAsGET-20]
	_ = x[FasterGetOrderForNames-21]
	_ = x[AllowV1Registration-22]
	_ = x[ParallelCheckFailedValidation-23]
	_ = x[DeleteUnusedChallenges-24]
	_ = x[V1DisableNewValidations-25]
	_ = x[PrecertificateOCSP-26]
	_ = x[PrecertificateRevocation-27]
}

const _FeatureFlag_name = "unusedPerformValidationRPCACME13KeyRolloverSimplifiedVAHTTPTLSSNIRevalidationAllowRenewalFirstRLSetIssuedNamesRenewalBitFasterRateLimitProbeCTLogsRevokeAtRACAAValidationMethodsCAAAccountURIHeadNonceStatusOKNewAuthorizationSchemaDisableAuthz2OrdersEarlyOrderRateLimitEnforceMultiVAMultiVAFullResultsRemoveWFE2AccountIDCheckRenewalFirstMandatoryPOSTAsGETFasterGetOrderForNamesAllowV1RegistrationParallelCheckFailedValidationDeleteUnusedChallengesV1DisableNewValidationsPrecertificateOCSPPrecertificateRevocation"

var _FeatureFlag_index = [...]uint16{0, 6, 26, 43, 59, 77, 96, 120, 135, 146, 156, 176, 189, 206, 228, 247, 266, 280, 298, 317, 334, 352, 374, 393, 422, 444, 467, 485, 509}

func (i FeatureFlag) String() string {
	if i < 0 || i >= FeatureFlag(len(_FeatureFlag_index)-1) {
		return "FeatureFlag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FeatureFlag_name[_FeatureFlag_index[i]:_FeatureFlag_index[i+1]]
}
