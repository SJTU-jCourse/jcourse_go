package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"jcourse_go/config"
	"jcourse_go/internal/domain/auth/repository"
)

type VerificationCodeService interface {
	SendCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email string, code string) error
	GenerateCode() string
}

var (
	ErrorCodeNotMatch = errors.New("code not match")
	ErrorRateLimit    = errors.New("rate limit")
)

type verificationCodeService struct {
	repo         repository.VerificationCodeRepository
	emailService EmailService
	codeTTL      time.Duration
	rateLimit    time.Duration
	codeLength   int
	codeCharSet  string
}

func (v *verificationCodeService) SendCode(ctx context.Context, email string) error {
	storeCode, err := v.repo.GetVerifyCode(ctx, email)
	if err != nil {
		return err
	}
	if storeCode.InRateLimit(v.rateLimit) {
		return ErrorRateLimit
	}

	code := v.GenerateCode()
	err = v.repo.StoreVerifyCode(ctx, email, code)
	if err != nil {
		return err
	}
	err = v.emailService.SendVerificationCode(ctx, email, code)
	if err != nil {
		return err
	}
	return nil
}

func (v *verificationCodeService) VerifyCode(ctx context.Context, email string, code string) error {
	storeCode, err := v.repo.GetVerifyCode(ctx, email)
	if err != nil {
		return err
	}
	err = v.repo.ClearHistory(ctx, email)
	if err != nil {
		return err
	}
	if !storeCode.IsValid(code) {
		return ErrorCodeNotMatch
	}
	return nil
}

func (v *verificationCodeService) GenerateCode() string {
	number := make([]byte, v.codeLength)
	maxIdx := int64(len(v.codeCharSet))
	for i := 0; i < v.codeLength; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(maxIdx))
		number[i] = v.codeCharSet[n.Int64()]
	}
	return string(number)
}

func NewVerificationCodeService(conf *config.Config, email EmailService, repo repository.VerificationCodeRepository) VerificationCodeService {
	return &verificationCodeService{
		repo:         repo,
		emailService: email,
		codeTTL:      time.Duration(int(time.Minute) * conf.VerifyCode.TTL),
		rateLimit:    time.Duration(int(time.Minute) * conf.VerifyCode.RateLimit),
		codeLength:   conf.VerifyCode.Length,
		codeCharSet:  conf.VerifyCode.Charset,
	}
}
