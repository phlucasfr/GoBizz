package util

import (
	"fmt"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func SendVerificationSms(phone *string) error {
	env := GetConfig(".")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioUsername,
		Password: env.TwilioPassword,
	})

	params := &verify.CreateVerificationParams{}
	params.SetTo("+55" + *phone)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(env.TwilioVerificationService, params)
	if err != nil {
		return fmt.Errorf("erro ao enviar o código de verificação: %w", err)
	}

	if resp.Sid != nil {
		fmt.Printf("Código de verificação enviado com sucesso. SID: %s\n", *resp.Sid)
	} else {
		return fmt.Errorf("erro desconhecido ao enviar o código de verificação")
	}

	return nil
}

func CheckVerificationCode(phone *string, code *string) (bool, error) {
	env := GetConfig(".")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioUsername,
		Password: env.TwilioPassword,
	})

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo("+55" + *phone)

	params.SetCode(*code)

	resp, err := client.VerifyV2.CreateVerificationCheck(env.TwilioVerificationService, params)
	if err != nil {
		return false, err
	}

	if resp.Status != nil && *resp.Status == "approved" {
		return true, nil
	}

	return false, nil
}
