package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	defaultAlipayClient "unibee/internal/logic/gateway/api/alipay/api"
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request/marketplace"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

func main() {
	const alipayGatewayUrl = ""
	const alipayClientId = ""
	const alipayMerchantPrivateKey = ""
	const alipayAlipayPublicKey = ""

	client := defaultAlipayClient.NewDefaultAlipayClient(
		alipayGatewayUrl,
		alipayClientId,
		alipayMerchantPrivateKey,
		alipayAlipayPublicKey, false)

	register(client)
	//update(client, "mid_zhangtianren_ztr_20230807_180716_981")
	//queryBalance(client, "mid_zhangtianren_ztr_20230807_180716_981")
}

func register(client *defaultAlipayClient.DefaultAlipayClient) {
	request, registerRequest := marketplace.NewAlipayRegisterRequest()
	registerRequest.RegistrationRequestId = uuid.NewString()
	registerRequest.SettlementInfos = []*model.SettlementInfo{
		{
			SettlementCurrency: "BRL",
			SettlementBankAccount: &model.SettlementBankAccount{
				BranchCode:        "1231",
				RoutingNumber:     "12",
				BankRegion:        "BR",
				AccountType:       model.AccountType_CHECKING,
				AccountHolderTIN:  "12345678901",
				AccountHolderName: "Timi",
				BankAccountNo:     "12312412421",
				AccountHolderType: model.AccountHolderType_ENTERPRISE,
			},
		},
	}

	registerRequest.MerchantInfo = &model.MerchantInfo{
		LoginId:         uuid.NewString()[:8] + "wangzunj3ao.wzj@digital-engine.com",
		LegalEntityType: model.LegalEntityType_Company,
		Company: &model.Company{
			LegalName: "legalName",
			RegisteredAddress: &model.Address{
				Region: "BR",
			},
			CompanyType: model.CompanyType_LTDA,
			OperatingAddress: &model.Address{
				Region:   "BR",
				Address1: "address1",
			},
			Attachments: []*model.Attachment{
				{
					AttachmentName: "1.jpg",
					FileKey:        "test",
					AttachmentType: model.AttachmentType_ARTICLES_OF_ASSOCIATION,
				},
				{
					AttachmentName: "2.jpg",
					FileKey:        "test",
					AttachmentType: model.AttachmentType_ASSOCIATION_ARTICLE,
				},
			},
			Certificates: &model.Certificate{
				CertificateNo:   "12312412",
				CertificateType: model.CertificateType_ENTERPRISE_REGISTRATION,
			},
		},
		BusinessInfo: &model.BusinessInfo{
			DoingBusinessAs: "businessName_DBA",
			Websites: []*model.WebSite{
				{
					Url: "www.global.alipay.com",
				},
			},
		},
		ReferenceMerchantId: uuid.NewString(),
		EntityAssociations: []*model.EntityAssociations{
			{
				Individual: &model.Individual{
					Certificates: &model.Certificate{
						CertificateNo:   "123214321",
						CertificateType: model.CertificateType_CPF,
					},
					Name: &model.UserName{
						FirstName:  "firstName",
						LastName:   "lastName",
						MiddleName: "middleName",
						FullName:   "fullName",
					},
					DateOfBirth: "1990-01-01",
				},
				LegalEntityType: model.LegalEntityType_INDIVIDUAL,
				AssociationType: model.AssociationType_LEGAL_REPRESENTATIVE,
			},
			{
				Individual: &model.Individual{
					Certificates: &model.Certificate{
						CertificateNo:   "123214321",
						CertificateType: model.CertificateType_CPF,
						FileKeys:        []string{"wetrewqratewtewgewgewg"},
					},
					Name: &model.UserName{
						FirstName:  "firstName",
						LastName:   "lastName",
						MiddleName: "middleName",
						FullName:   "fullName",
					},
					DateOfBirth: "1990-01-01",
				},
				LegalEntityType: model.LegalEntityType_INDIVIDUAL,
				AssociationType: model.AssociationType_UBO,
			},
		},
	}

	response, err := client.Execute(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}

}

func update(client *defaultAlipayClient.DefaultAlipayClient, referenceMerchantId string) {
	request, updateRequest := marketplace.NewAlipaySettlementInfoUpdateRequest()
	updateRequest.ReferenceMerchantId = referenceMerchantId
	updateRequest.UpdateRequestId = uuid.NewString()
	updateRequest.SettlementCurrency = "BRL"
	updateRequest.SettlementBankAccount = &model.SettlementBankAccount{
		BranchCode:        "1231",
		RoutingNumber:     "12",
		BankRegion:        "BR",
		AccountType:       model.AccountType_CHECKING,
		AccountHolderTIN:  "12345678901",
		AccountHolderName: "Timi",
		BankAccountNo:     "12312412421",
		AccountHolderType: model.AccountHolderType_ENTERPRISE,
	}

	response, err := client.Execute(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
}

func queryBalance(client *defaultAlipayClient.DefaultAlipayClient, referenceMerchantId string) {
	request, queryBalanceRequest := marketplace.NewAlipayInquireBalanceRequest()
	queryBalanceRequest.ReferenceMerchantId = referenceMerchantId
	response, err := client.Execute(request)
	if err != nil {
		fmt.Println(err)
	} else {
		indent, _ := json.MarshalIndent(response.(*responseMarketplace.AlipayInquireBalanceResponse), "", "  ")
		fmt.Println(string(indent))
	}
}

func settleRequest(client *defaultAlipayClient.DefaultAlipayClient, paymentId string) {
	request, settleRequest := marketplace.NewAlipaySettleRequest()
	settleRequest.PaymentId = paymentId
	settleRequest.SettlementRequestId = uuid.NewString()
	settleRequest.SettlementDetails = []*model.SettlementDetail{
		{
			SettleTo: model.SettleToType_SELLER,
			SettlementAmount: &model.Amount{
				Currency: "BRL",
				Value:    "90",
			},
		},
		{
			SettleTo: model.SettleToType_MARKETPLACE,
			SettlementAmount: &model.Amount{
				Currency: "BRL",
				Value:    "10",
			},
		},
	}

	response, err := client.Execute(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
}

func createPayout(client *defaultAlipayClient.DefaultAlipayClient) {
	request, createPayoutRequest := marketplace.NewAlipayCreatePayoutRequest()
	createPayoutRequest.TransferRequestId = uuid.NewString()
	createPayoutRequest.TransferFromDetail = &model.TransferFromDetail{
		TransferFromMethod: &model.PaymentMethod{
			PaymentMethodId:   uuid.NewString(),
			PaymentMethodType: model.BALANCE_ACCOUNT,
		},
	}
	createPayoutRequest.TransferToDetail = &model.TransferToDetail{
		TransferToMethod: &model.PaymentMethod{
			PaymentMethodId:   uuid.NewString(),
			PaymentMethodType: model.SETTLEMENT_CARD,
		},
		TransferToCurrency: "BRL",
		PurposeCode:        "GSD",
		TransferNotifyUrl:  "www.global.alipay.com",
	}

	response, err := client.Execute(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}

}

func createTransfer(client *defaultAlipayClient.DefaultAlipayClient) {
	request, createTransferRequest := marketplace.NewAlipayCreateTransferRequest()
	createTransferRequest.TransferRequestId = uuid.NewString()

	createTransferRequest.TransferFromDetail = &model.TransferFromDetail{
		TransferFromMethod: &model.PaymentMethod{
			PaymentMethodId:   uuid.NewString(),
			PaymentMethodType: model.BALANCE_ACCOUNT,
		},
		TransferFromAmount: model.NewAmount("100", "BRL"),
	}

	createTransferRequest.TransferToDetail = &model.TransferToDetail{
		TransferToMethod: &model.PaymentMethod{
			PaymentMethodId:   uuid.NewString(),
			PaymentMethodType: model.BALANCE_ACCOUNT,
		},
		TransferToCurrency: "BRL",
		PurposeCode:        "GSD",
		TransferNotifyUrl:  "www.global.alipay.com",
	}

	response, err := client.Execute(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
}
