// SGX Quote Verification Service
//
// SGX Quote Verification service (SQVS) is used to verify the sgx ecdsa quote provided by the Quote verifier.
// SQVS checks whether the quote signature and PCK Certificate are valid, checks other parameters in the quote and returns the verification result.
// SQVS contacts SGX Caching service (SCS) to make sure that PCKCRL, TCBInfo, and QEIdentity in the quote are correct.
// SQVS listening port is user-configurable.
//
//  License: Copyright (C) 2020 Intel Corporation. SPDX-License-Identifier: BSD-3-Clause
//
//  Version: 1.0
//  Host: svs.com:12000
//  BasePath: /svs
//
//  Schemes: https
//
//  SecurityDefinitions:
//   bearerAuth:
//     type: apiKey
//     in: header
//     name: Authorization
//     description: Enter your bearer token in the format **Bearer &lt;token&gt;**
//
// swagger:meta
package docs

import "intel/isecl/sqvs/v5/resource"

// QuoteData request payload
// swagger:parameters QuoteData
type QuoteDataInfo struct {
	// in:body
	Body resource.QuoteData
}

// QuoteDataWithChallenge request payload
// swagger:parameters QuoteDataWithChallenge
type QuoteDataAndChallengeInfo struct {
	// in:body
	Body resource.QuoteDataWithChallenge
}

// SGXResponse response payload
// swagger:response SGXResponse
type SGXResponseInfo struct {
	// in:body
	Body resource.SGXResponse
}

// SignedSGXResponse response payload
// swagger:response SignedSGXResponse
type SignedSGXResponseInfo struct {
	// in:body
	Body resource.SignedSGXResponse
}

// swagger:operation POST /v1/sgx_qv_verify_quote Quote sgxVerifyQuote
// ---
// description: |
//   Verifies the SGX ECDSA quote provided in the request body.
//   Quote verifier requests SGX Quote Verification Service (SQVS) to verify quote.
//   SQVS parses the quote, verifies all the parameters in the quote and returns the response.
//
// security:
//  - bearerAuth: []
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - name: request body
//   required: true
//   in: body
//   schema:
//     "$ref": "#/definitions/QuoteData"
// responses:
//   '200':
//     description: Successfully verified the quote and its parameters.
//     schema:
//       "$ref": "#/definitions/SGXResponse"
//
// x-sample-call-endpoint: https://svs.com:12000/svs/v1/sgx_qv_verify_quote
// x-sample-call-input: |
//  {
//    "quote":"AwACAAAAAAAFAAoAk5pyM/ecTKmUCg2zlX8GBykaeX/IbYBtMzyC5lkpA0IAAAAAAgIAAQEAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQAAAAAAAADnAAAAAAAAAJJwRC0b0ZYfo52+HyzfT4eVClT8r5ouUBOHXDNGV
//    C3KAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADUEqTwfvg4kqWRX7KrWEvjHhhuWk+Vq19pUP1OuGlNewAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1BAAAI2X3
//    h5lqazTamLWIQP8xEPDAwVi+oQDogM3xLv4bMtf4mZm7KkIHHI1x693ILPWxrIpaRiEggRGdheZ5NCdEwLyPIzUpScz6bN8oQo2
//    xaa8p/W1lIJhG7Yxd3c71Plz6kpi+bkaAo/ZbnDnh5PncIu8H00SCxdcZDOYAU/ualbsAgIAAQEAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFQAAAAAAAADnAAAAAAAAAGDYWvKL6NHECgjZiwCdX4rME4Sjhc9GCADkeHkdGpecAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACMT1d115ZQPpYTf3fGioKaAFasje1wFAsIGwlEkMV7/wAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAEABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    ceWekTM/vT9LqUAxTrRfBT6Eav49bou/MlhGY0GA4wwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYcmJZoOnc26Fwd
//    1goMVtx/40ZS+7077Zqdo1gcZfG2okXCrtC1vKkiR4bz6bmPe8r8OsBBpHfV83nandZ2kKRSAAAAECAwQFBgcICQoLDA0ODxARE
//    hMUFRYXGBkaGxwdHh8FAGwOAAAtLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJRTlEQ0NCSnFnQXdJQkFnSVVPeEl5dnQ5
//    cm5ybVhZWktQQjhmRzFEVFNpTDh3Q2dZSUtvWkl6ajBFQXdJd2NERWlNQ0FHQTFVRQpBd3daU1c1MFpXd2dVMGRZSUZCRFN5QlF
//    iR0YwWm05eWJTQkRRVEVhTUJnR0ExVUVDZ3dSU1c1MFpXd2dRMjl5Y0c5eVlYUnBiMjR4CkZEQVNCZ05WQkFjTUMxTmhiblJoSU
//    VOc1lYSmhNUXN3Q1FZRFZRUUlEQUpEUVRFTE1Ba0dBMVVFQmhNQ1ZWTXdIaGNOTWpBeE1qRTEKTVRRMU1UUXhXaGNOTWpjeE1qR
//    TFNVFExTVRReFdqQndNU0l3SUFZRFZRUUREQmxKYm5SbGJDQlRSMWdnVUVOTElFTmxjblJwWm1sagpZWFJsTVJvd0dBWURWUVFL
//    REJGSmJuUmxiQ0JEYjNKd2IzSmhkR2x2YmpFVU1CSUdBMVVFQnd3TFUyRnVkR0VnUTJ4aGNtRXhDekFKCkJnTlZCQWdNQWtOQk1
//    Rc3dDUVlEVlFRR0V3SlZVekJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCSDV5b1RTalhuS3AKMFd3Y1orYndTV0
//    9rV1VLMmRzbnIyb1pXV0I2S3oxRlBEdVk2ekFFenNEUGl2RWlqWWRxRGJPbGw2T1BUc28rMmYwVVlJNlJvL29TagpnZ01RTUlJR
//    EREQWZCZ05WSFNNRUdEQVdnQlJaSTlPblNxaGpWQzQ1Y0szZ0R3Y3JWeVFxdHpCdkJnTlZIUjhFYURCbU1HU2dZcUJnCmhsNW9k
//    SFJ3Y3pvdkwzTmllQzVoY0drdWRISjFjM1JsWkhObGNuWnBZMlZ6TG1sdWRHVnNMbU52YlM5elozZ3ZZMlZ5ZEdsbWFXTmgKZEd
//    sdmJpOTJNeTl3WTJ0amNtdy9ZMkU5Y0d4aGRHWnZjbTBtWlc1amIyUnBibWM5WkdWeU1CMEdBMVVkRGdRV0JCVGY1Z1JSWlkvVA
//    pSTXAxb0hKVmFQMEFvTWc2K2pBT0JnTlZIUThCQWY4RUJBTUNCc0F3REFZRFZSMFRBUUgvQkFJd0FEQ0NBamtHQ1NxR1NJYjRUU
//    UVOCkFRU0NBaW93Z2dJbU1CNEdDaXFHU0liNFRRRU5BUUVFRUQ1RkZvQlhLSCtaN0hFQlRmcGsxK1F3Z2dGakJnb3Foa2lHK0Uw
//    QkRRRUMKTUlJQlV6QVFCZ3NxaGtpRytFMEJEUUVDQVFJQkFqQVFCZ3NxaGtpRytFMEJEUUVDQWdJQkFqQVFCZ3NxaGtpRytFMEJ
//    EUUVDQXdJQgpBREFRQmdzcWhraUcrRTBCRFFFQ0JBSUJBREFRQmdzcWhraUcrRTBCRFFFQ0JRSUJBREFRQmdzcWhraUcrRTBCRF
//    FFQ0JnSUJBREFRCkJnc3Foa2lHK0UwQkRRRUNCd0lCQURBUUJnc3Foa2lHK0UwQkRRRUNDQUlCQURBUUJnc3Foa2lHK0UwQkRRR
//    UNDUUlCQURBUUJnc3EKaGtpRytFMEJEUUVDQ2dJQkFEQVFCZ3NxaGtpRytFMEJEUUVDQ3dJQkFEQVFCZ3NxaGtpRytFMEJEUUVD
//    REFJQkFEQVFCZ3NxaGtpRworRTBCRFFFQ0RRSUJBREFRQmdzcWhraUcrRTBCRFFFQ0RnSUJBREFRQmdzcWhraUcrRTBCRFFFQ0R
//    3SUJBREFRQmdzcWhraUcrRTBCCkRRRUNFQUlCQURBUUJnc3Foa2lHK0UwQkRRRUNFUUlCQ2pBZkJnc3Foa2lHK0UwQkRRRUNFZ1
//    FRQWdJQUFBQUFBQUFBQUFBQUFBQUEKQURBUUJnb3Foa2lHK0UwQkRRRURCQUlBQURBVUJnb3Foa2lHK0UwQkRRRUVCQVlnWUdvQ
//    UFBQXdEd1lLS29aSWh2aE5BUTBCQlFvQgpBVEFlQmdvcWhraUcrRTBCRFFFR0JCQlFsT3d5R3BWS0dqYVVMelFCM3N0UU1FUUdD
//    aXFHU0liNFRRRU5BUWN3TmpBUUJnc3Foa2lHCitFMEJEUUVIQVFFQi96QVFCZ3NxaGtpRytFMEJEUUVIQWdFQkFEQVFCZ3NxaGt
//    pRytFMEJEUUVIQXdFQi96QUtCZ2dxaGtqT1BRUUQKQWdOSUFEQkZBaUVBcEtJMDBYZG5GaU1jMkh6ZVpNT1FLQTFraHlBZTFKZH
//    hHWCtyRTBlUEVvSUNJR3lDcE1TM1ZmTjJTVWlPTjJIRwo3MndDU2o4WmNqZk93aTd1OHBRVWVadjkKLS0tLS1FTkQgQ0VSVElGS
//    UNBVEUtLS0tLS0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlDbWpDQ0FrQ2dBd0lCQWdJVVdTUFRwMHFvWTFRdU9YQ3Q0
//    QThISzFja0tyY3dDZ1lJS29aSXpqMEVBd0l3CmFERWFNQmdHQTFVRUF3d1JTVzUwWld3Z1UwZFlJRkp2YjNRZ1EwRXhHakFZQmd
//    OVkJBb01FVWx1ZEdWc0lFTnYKY25CdmNtRjBhVzl1TVJRd0VnWURWUVFIREF0VFlXNTBZU0JEYkdGeVlURUxNQWtHQTFVRUNBd0
//    NRMEV4Q3pBSgpCZ05WQkFZVEFsVlRNQjRYRFRFNU1UQXpNVEV5TXpNME4xb1hEVE0wTVRBek1URXlNek0wTjFvd2NERWlNQ0FHC
//    kExVUVBd3daU1c1MFpXd2dVMGRZSUZCRFN5QlFiR0YwWm05eWJTQkRRVEVhTUJnR0ExVUVDZ3dSU1c1MFpXd2cKUTI5eWNHOXlZ
//    WFJwYjI0eEZEQVNCZ05WQkFjTUMxTmhiblJoSUVOc1lYSmhNUXN3Q1FZRFZRUUlEQUpEUVRFTApNQWtHQTFVRUJoTUNWVk13V1R
//    BVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVF3cCtMYytUVUJ0ZzFICitVOEpJc01zYmpIakNrVHRYYjhqUE02cjJkaH
//    U5eklibGhEWjdJTmZxdDNJeDhYY0ZLRDhrME5FWHJrWjY2cUoKWGExS3pMSUtvNEcvTUlHOE1COEdBMVVkSXdRWU1CYUFGT25vU
//    kZKVE5seExHSm9SL0VNWUxLWGNJSUJJTUZZRwpBMVVkSHdSUE1FMHdTNkJKb0VlR1JXaDBkSEJ6T2k4dmMySjRMV05sY25ScFpt
//    bGpZWFJsY3k1MGNuVnpkR1ZrCmMyVnlkbWxqWlhNdWFXNTBaV3d1WTI5dEwwbHVkR1ZzVTBkWVVtOXZkRU5CTG1SbGNqQWRCZ05
//    WSFE0RUZnUVUKV1NQVHAwcW9ZMVF1T1hDdDRBOEhLMWNrS3Jjd0RnWURWUjBQQVFIL0JBUURBZ0VHTUJJR0ExVWRFd0VCL3dRSQ
//    pNQVlCQWY4Q0FRQXdDZ1lJS29aSXpqMEVBd0lEU0FBd1JRSWhBSjFxK0ZUeitnVXVWZkJRdUNnSnNGckwyVFRTCmUxYUJaNTNPN
//    TJUakZpZTZBaUFyaVBhUmFoVVg5T2E5a0dMbEFjaFdYS1Q2ajRSV1NSNTBCcWhyTjNVVDRBPT0KLS0tLS1FTkQgQ0VSVElGSUNB
//    VEUtLS0tLQotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJQ2xEQ0NBam1nQXdJQkFnSVZBT25vUkZKVE5seExHSm9SL0V
//    NWUxLWGNJSUJJTUFvR0NDcUdTTTQ5QkFNQwpNR2d4R2pBWUJnTlZCQU1NRVVsdWRHVnNJRk5IV0NCU2IyOTBJRU5CTVJvd0dBWU
//    RWUVFLREJGSmJuUmxiQ0JECmIzSndiM0poZEdsdmJqRVVNQklHQTFVRUJ3d0xVMkZ1ZEdFZ1EyeGhjbUV4Q3pBSkJnTlZCQWdNQ
//    WtOQk1Rc3cKQ1FZRFZRUUdFd0pWVXpBZUZ3MHhPVEV3TXpFd09UUTVNakZhRncwME9URXlNekV5TXpVNU5UbGFNR2d4R2pBWQpC
//    Z05WQkFNTUVVbHVkR1ZzSUZOSFdDQlNiMjkwSUVOQk1Sb3dHQVlEVlFRS0RCRkpiblJsYkNCRGIzSndiM0poCmRHbHZiakVVTUJ
//    JR0ExVUVCd3dMVTJGdWRHRWdRMnhoY21FeEN6QUpCZ05WQkFnTUFrTkJNUXN3Q1FZRFZRUUcKRXdKVlV6QlpNQk1HQnlxR1NNND
//    lBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJFLzZELzFXSE5yV3dQbU5NSXlCS01XNQpKNkp6TXNqbzZ4UDJ2a0sxY2RaR2IxUEdSUC9DL
//    zhFQ2dpRGtta2xtendMekxpKzAwMG03TExydEtKQTNvQzJqCmdiOHdnYnd3SHdZRFZSMGpCQmd3Rm9BVTZlaEVVbE0yWEVzWW1o
//    SDhReGdzcGR3Z2dFZ3dWZ1lEVlIwZkJFOHcKVFRCTG9FbWdSNFpGYUhSMGNITTZMeTl6WW5ndFkyVnlkR2xtYVdOaGRHVnpMblJ
//    5ZFhOMFpXUnpaWEoyYVdObApjeTVwYm5SbGJDNWpiMjB2U1c1MFpXeFRSMWhTYjI5MFEwRXVaR1Z5TUIwR0ExVWREZ1FXQkJUcD
//    ZFUlNVelpjClN4aWFFZnhER0N5bDNDQ0FTREFPQmdOVkhROEJBZjhFQkFNQ0FRWXdFZ1lEVlIwVEFRSC9CQWd3QmdFQi93SUIKQ
//    VRBS0JnZ3Foa2pPUFFRREFnTkpBREJHQWlFQXp3OXpkVWlVSFBNVWQwQzRteDQxamxGWmtyTTN5NWYxbGduVgpPN0Ziak9vQ0lR
//    Q29HdFVtVDRjWHQ3Vit5U0hiSjhIb2I5QWFucHZYTkgxRVIrL2daRitvcFE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
//    "userData": ""
//  }
// x-sample-call-output: |
//  {
//    "Message": "SGX_QL_QV_RESULT_OK",
//    "reportData": "0000000000000000000000000000000000000000000000000000000000000000",
//    "userDataMatch": "false",
//    "EnclaveIssuer": "d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b",
//    "EnclaveMeasurement": "9270442d1bd1961fa39dbe1f2cdf4f87950a54fcaf9a2e5013875c3346542dca",
//    "EnclaveIssuerProdID": "00",
//    "EnclaveIssuerExtProdID": "00000000000000000000000000000000",
//    "ConfigSvn": "00",
//    "IsvSvn": "01",
//    "ConfigId": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
//    "TcbLevel": "OutofDate"
//  }
// ---

// swagger:operation POST /v2/sgx_qv_verify_quote Quote sgxVerifyQuoteAndSign
// ---
// description: |
//   Verifies the SGX ECDSA quote provided in the request body.
//   Quote verifier requests SGX Quote Verification Service (SQVS) to verify quote.
//   SQVS parses the quote, verifies all the parameters in the quote and returns the response.
//   It signs the quote verification response in case it is configured to do so.
//
// security:
//  - bearerAuth: []
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - name: request body
//   required: true
//   in: body
//   schema:
//     "$ref": "#/definitions/QuoteDataWithChallenge"
// responses:
//   '200':
//     description: Successfully verified the quote and its parameters.
//     schema:
//       "$ref": "#/definitions/SignedSGXResponse"
//
// x-sample-call-endpoint: https://svs.com:12000/svs/v2/sgx_qv_verify_quote
// x-sample-call-input: |
//  {
//    "quote":"AwACAAAAAAAFAAoAk5pyM/ecTKmUCg2zlX8GBykaeX/IbYBtMzyC5lkpA0IAAAAAAgIAAQEAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQAAAAAAAADnAAAAAAAAAJJwRC0b0ZYfo52+HyzfT4eVClT8r5ouUBOHXDNGV
//    C3KAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADUEqTwfvg4kqWRX7KrWEvjHhhuWk+Vq19pUP1OuGlNewAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1BAAAI2X3
//    h5lqazTamLWIQP8xEPDAwVi+oQDogM3xLv4bMtf4mZm7KkIHHI1x693ILPWxrIpaRiEggRGdheZ5NCdEwLyPIzUpScz6bN8oQo2
//    xaa8p/W1lIJhG7Yxd3c71Plz6kpi+bkaAo/ZbnDnh5PncIu8H00SCxdcZDOYAU/ualbsAgIAAQEAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFQAAAAAAAADnAAAAAAAAAGDYWvKL6NHECgjZiwCdX4rME4Sjhc9GCADkeHkdGpecAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACMT1d115ZQPpYTf3fGioKaAFasje1wFAsIGwlEkMV7/wAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    AAAAAAAAAAAAAAEABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
//    ceWekTM/vT9LqUAxTrRfBT6Eav49bou/MlhGY0GA4wwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYcmJZoOnc26Fwd
//    1goMVtx/40ZS+7077Zqdo1gcZfG2okXCrtC1vKkiR4bz6bmPe8r8OsBBpHfV83nandZ2kKRSAAAAECAwQFBgcICQoLDA0ODxARE
//    hMUFRYXGBkaGxwdHh8FAGwOAAAtLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJRTlEQ0NCSnFnQXdJQkFnSVVPeEl5dnQ5
//    cm5ybVhZWktQQjhmRzFEVFNpTDh3Q2dZSUtvWkl6ajBFQXdJd2NERWlNQ0FHQTFVRQpBd3daU1c1MFpXd2dVMGRZSUZCRFN5QlF
//    iR0YwWm05eWJTQkRRVEVhTUJnR0ExVUVDZ3dSU1c1MFpXd2dRMjl5Y0c5eVlYUnBiMjR4CkZEQVNCZ05WQkFjTUMxTmhiblJoSU
//    VOc1lYSmhNUXN3Q1FZRFZRUUlEQUpEUVRFTE1Ba0dBMVVFQmhNQ1ZWTXdIaGNOTWpBeE1qRTEKTVRRMU1UUXhXaGNOTWpjeE1qR
//    TFNVFExTVRReFdqQndNU0l3SUFZRFZRUUREQmxKYm5SbGJDQlRSMWdnVUVOTElFTmxjblJwWm1sagpZWFJsTVJvd0dBWURWUVFL
//    REJGSmJuUmxiQ0JEYjNKd2IzSmhkR2x2YmpFVU1CSUdBMVVFQnd3TFUyRnVkR0VnUTJ4aGNtRXhDekFKCkJnTlZCQWdNQWtOQk1
//    Rc3dDUVlEVlFRR0V3SlZVekJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCSDV5b1RTalhuS3AKMFd3Y1orYndTV0
//    9rV1VLMmRzbnIyb1pXV0I2S3oxRlBEdVk2ekFFenNEUGl2RWlqWWRxRGJPbGw2T1BUc28rMmYwVVlJNlJvL29TagpnZ01RTUlJR
//    EREQWZCZ05WSFNNRUdEQVdnQlJaSTlPblNxaGpWQzQ1Y0szZ0R3Y3JWeVFxdHpCdkJnTlZIUjhFYURCbU1HU2dZcUJnCmhsNW9k
//    SFJ3Y3pvdkwzTmllQzVoY0drdWRISjFjM1JsWkhObGNuWnBZMlZ6TG1sdWRHVnNMbU52YlM5elozZ3ZZMlZ5ZEdsbWFXTmgKZEd
//    sdmJpOTJNeTl3WTJ0amNtdy9ZMkU5Y0d4aGRHWnZjbTBtWlc1amIyUnBibWM5WkdWeU1CMEdBMVVkRGdRV0JCVGY1Z1JSWlkvVA
//    pSTXAxb0hKVmFQMEFvTWc2K2pBT0JnTlZIUThCQWY4RUJBTUNCc0F3REFZRFZSMFRBUUgvQkFJd0FEQ0NBamtHQ1NxR1NJYjRUU
//    UVOCkFRU0NBaW93Z2dJbU1CNEdDaXFHU0liNFRRRU5BUUVFRUQ1RkZvQlhLSCtaN0hFQlRmcGsxK1F3Z2dGakJnb3Foa2lHK0Uw
//    QkRRRUMKTUlJQlV6QVFCZ3NxaGtpRytFMEJEUUVDQVFJQkFqQVFCZ3NxaGtpRytFMEJEUUVDQWdJQkFqQVFCZ3NxaGtpRytFMEJ
//    EUUVDQXdJQgpBREFRQmdzcWhraUcrRTBCRFFFQ0JBSUJBREFRQmdzcWhraUcrRTBCRFFFQ0JRSUJBREFRQmdzcWhraUcrRTBCRF
//    FFQ0JnSUJBREFRCkJnc3Foa2lHK0UwQkRRRUNCd0lCQURBUUJnc3Foa2lHK0UwQkRRRUNDQUlCQURBUUJnc3Foa2lHK0UwQkRRR
//    UNDUUlCQURBUUJnc3EKaGtpRytFMEJEUUVDQ2dJQkFEQVFCZ3NxaGtpRytFMEJEUUVDQ3dJQkFEQVFCZ3NxaGtpRytFMEJEUUVD
//    REFJQkFEQVFCZ3NxaGtpRworRTBCRFFFQ0RRSUJBREFRQmdzcWhraUcrRTBCRFFFQ0RnSUJBREFRQmdzcWhraUcrRTBCRFFFQ0R
//    3SUJBREFRQmdzcWhraUcrRTBCCkRRRUNFQUlCQURBUUJnc3Foa2lHK0UwQkRRRUNFUUlCQ2pBZkJnc3Foa2lHK0UwQkRRRUNFZ1
//    FRQWdJQUFBQUFBQUFBQUFBQUFBQUEKQURBUUJnb3Foa2lHK0UwQkRRRURCQUlBQURBVUJnb3Foa2lHK0UwQkRRRUVCQVlnWUdvQ
//    UFBQXdEd1lLS29aSWh2aE5BUTBCQlFvQgpBVEFlQmdvcWhraUcrRTBCRFFFR0JCQlFsT3d5R3BWS0dqYVVMelFCM3N0UU1FUUdD
//    aXFHU0liNFRRRU5BUWN3TmpBUUJnc3Foa2lHCitFMEJEUUVIQVFFQi96QVFCZ3NxaGtpRytFMEJEUUVIQWdFQkFEQVFCZ3NxaGt
//    pRytFMEJEUUVIQXdFQi96QUtCZ2dxaGtqT1BRUUQKQWdOSUFEQkZBaUVBcEtJMDBYZG5GaU1jMkh6ZVpNT1FLQTFraHlBZTFKZH
//    hHWCtyRTBlUEVvSUNJR3lDcE1TM1ZmTjJTVWlPTjJIRwo3MndDU2o4WmNqZk93aTd1OHBRVWVadjkKLS0tLS1FTkQgQ0VSVElGS
//    UNBVEUtLS0tLS0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlDbWpDQ0FrQ2dBd0lCQWdJVVdTUFRwMHFvWTFRdU9YQ3Q0
//    QThISzFja0tyY3dDZ1lJS29aSXpqMEVBd0l3CmFERWFNQmdHQTFVRUF3d1JTVzUwWld3Z1UwZFlJRkp2YjNRZ1EwRXhHakFZQmd
//    OVkJBb01FVWx1ZEdWc0lFTnYKY25CdmNtRjBhVzl1TVJRd0VnWURWUVFIREF0VFlXNTBZU0JEYkdGeVlURUxNQWtHQTFVRUNBd0
//    NRMEV4Q3pBSgpCZ05WQkFZVEFsVlRNQjRYRFRFNU1UQXpNVEV5TXpNME4xb1hEVE0wTVRBek1URXlNek0wTjFvd2NERWlNQ0FHC
//    kExVUVBd3daU1c1MFpXd2dVMGRZSUZCRFN5QlFiR0YwWm05eWJTQkRRVEVhTUJnR0ExVUVDZ3dSU1c1MFpXd2cKUTI5eWNHOXlZ
//    WFJwYjI0eEZEQVNCZ05WQkFjTUMxTmhiblJoSUVOc1lYSmhNUXN3Q1FZRFZRUUlEQUpEUVRFTApNQWtHQTFVRUJoTUNWVk13V1R
//    BVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVF3cCtMYytUVUJ0ZzFICitVOEpJc01zYmpIakNrVHRYYjhqUE02cjJkaH
//    U5eklibGhEWjdJTmZxdDNJeDhYY0ZLRDhrME5FWHJrWjY2cUoKWGExS3pMSUtvNEcvTUlHOE1COEdBMVVkSXdRWU1CYUFGT25vU
//    kZKVE5seExHSm9SL0VNWUxLWGNJSUJJTUZZRwpBMVVkSHdSUE1FMHdTNkJKb0VlR1JXaDBkSEJ6T2k4dmMySjRMV05sY25ScFpt
//    bGpZWFJsY3k1MGNuVnpkR1ZrCmMyVnlkbWxqWlhNdWFXNTBaV3d1WTI5dEwwbHVkR1ZzVTBkWVVtOXZkRU5CTG1SbGNqQWRCZ05
//    WSFE0RUZnUVUKV1NQVHAwcW9ZMVF1T1hDdDRBOEhLMWNrS3Jjd0RnWURWUjBQQVFIL0JBUURBZ0VHTUJJR0ExVWRFd0VCL3dRSQ
//    pNQVlCQWY4Q0FRQXdDZ1lJS29aSXpqMEVBd0lEU0FBd1JRSWhBSjFxK0ZUeitnVXVWZkJRdUNnSnNGckwyVFRTCmUxYUJaNTNPN
//    TJUakZpZTZBaUFyaVBhUmFoVVg5T2E5a0dMbEFjaFdYS1Q2ajRSV1NSNTBCcWhyTjNVVDRBPT0KLS0tLS1FTkQgQ0VSVElGSUNB
//    VEUtLS0tLQotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJQ2xEQ0NBam1nQXdJQkFnSVZBT25vUkZKVE5seExHSm9SL0V
//    NWUxLWGNJSUJJTUFvR0NDcUdTTTQ5QkFNQwpNR2d4R2pBWUJnTlZCQU1NRVVsdWRHVnNJRk5IV0NCU2IyOTBJRU5CTVJvd0dBWU
//    RWUVFLREJGSmJuUmxiQ0JECmIzSndiM0poZEdsdmJqRVVNQklHQTFVRUJ3d0xVMkZ1ZEdFZ1EyeGhjbUV4Q3pBSkJnTlZCQWdNQ
//    WtOQk1Rc3cKQ1FZRFZRUUdFd0pWVXpBZUZ3MHhPVEV3TXpFd09UUTVNakZhRncwME9URXlNekV5TXpVNU5UbGFNR2d4R2pBWQpC
//    Z05WQkFNTUVVbHVkR1ZzSUZOSFdDQlNiMjkwSUVOQk1Sb3dHQVlEVlFRS0RCRkpiblJsYkNCRGIzSndiM0poCmRHbHZiakVVTUJ
//    JR0ExVUVCd3dMVTJGdWRHRWdRMnhoY21FeEN6QUpCZ05WQkFnTUFrTkJNUXN3Q1FZRFZRUUcKRXdKVlV6QlpNQk1HQnlxR1NNND
//    lBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJFLzZELzFXSE5yV3dQbU5NSXlCS01XNQpKNkp6TXNqbzZ4UDJ2a0sxY2RaR2IxUEdSUC9DL
//    zhFQ2dpRGtta2xtendMekxpKzAwMG03TExydEtKQTNvQzJqCmdiOHdnYnd3SHdZRFZSMGpCQmd3Rm9BVTZlaEVVbE0yWEVzWW1o
//    SDhReGdzcGR3Z2dFZ3dWZ1lEVlIwZkJFOHcKVFRCTG9FbWdSNFpGYUhSMGNITTZMeTl6WW5ndFkyVnlkR2xtYVdOaGRHVnpMblJ
//    5ZFhOMFpXUnpaWEoyYVdObApjeTVwYm5SbGJDNWpiMjB2U1c1MFpXeFRSMWhTYjI5MFEwRXVaR1Z5TUIwR0ExVWREZ1FXQkJUcD
//    ZFUlNVelpjClN4aWFFZnhER0N5bDNDQ0FTREFPQmdOVkhROEJBZjhFQkFNQ0FRWXdFZ1lEVlIwVEFRSC9CQWd3QmdFQi93SUIKQ
//    VRBS0JnZ3Foa2pPUFFRREFnTkpBREJHQWlFQXp3OXpkVWlVSFBNVWQwQzRteDQxamxGWmtyTTN5NWYxbGduVgpPN0Ziak9vQ0lR
//    Q29HdFVtVDRjWHQ3Vit5U0hiSjhIb2I5QWFucHZYTkgxRVIrL2daRitvcFE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
//    "userData": "",
//    "challenge": "DJ4m0A9eBwTUuuiJOwi5ALgyMP5X99KH+afqF6qjn0ImiA2ej8LnNgV377sdsS17JRkHJWzJbucmHufcuRtpfA=="
//  }
// x-sample-call-output: |
//  {
//    "quoteData": {
//      "Message": "SGX_QL_QV_RESULT_OK",
//      "reportData": "0000000000000000000000000000000000000000000000000000000000000000",
//      "userDataMatch": "false",
//      "EnclaveIssuer": "d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b",
//      "EnclaveMeasurement": "9270442d1bd1961fa39dbe1f2cdf4f87950a54fcaf9a2e5013875c3346542dca",
//      "EnclaveIssuerProdID": "00",
//      "EnclaveIssuerExtProdID": "00000000000000000000000000000000",
//      "ConfigSvn": "00",
//      "IsvSvn": "01",
//      "ConfigId": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
//      "TcbLevel": "OutofDate",
//      "Quote": "<quote in request>",
//      "Challenge": "DJ4m0A9eBwTUuuiJOwi5ALgyMP5X99KH+afqF6qjn0ImiA2ej8LnNgV377sdsS17JRkHJWzJbucmHufcuRtpfA=="
//    },
//    "signature": "<Digital Signature of the content of response field>",
//	  "certificateChain": "<Chain of leaf and intermediate signing certs in PEM format>"
//  }

// ---
