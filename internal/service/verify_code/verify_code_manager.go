package verify_code

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/unify-z/go-surl/internal/config"
	"github.com/unify-z/go-surl/internal/helpers"
	"github.com/unify-z/go-surl/internal/utils"
)

type VerifyCodeManager struct {
	codes      map[string]string
	rateLimit  map[string]int
	smtpHelper *helpers.SMTPHelper
	mu         sync.Mutex
}

func NewVerifyCodeManager(smtpHelper helpers.SMTPHelper) *VerifyCodeManager {
	v := &VerifyCodeManager{
		codes:      map[string]string{},
		rateLimit:  make(map[string]int),
		smtpHelper: &smtpHelper,
	}
	go v.clearRateLimit()
	return v
}

func (v *VerifyCodeManager) clearRateLimit() {
	ticker := time.NewTicker(12 * time.Hour)
	for {
		<-ticker.C
		v.mu.Lock()
		v.rateLimit = make(map[string]int)
		v.mu.Unlock()
	}
}

func (v *VerifyCodeManager) IsRateLimited(email string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	if count, exists := v.rateLimit[email]; exists && count >= 5 {
		return true
	}
	return false
}

func (v *VerifyCodeManager) GenerateCode(email string) string {
	v.mu.Lock()
	defer v.mu.Unlock()
	code := utils.RandStr(6)
	v.codes[email] = code
	return code
}

func (v *VerifyCodeManager) VerifyCode(email, code string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	fmt.Print(email, code, v.codes)
	if storedCode, exists := v.codes[email]; exists && storedCode == code {
		delete(v.codes, email)
		return true
	}
	return false
}

func (v *VerifyCodeManager) SendCode(email string) error {
	if v.IsRateLimited(email) {
		return errors.New("rate limit exceeded")
	}

	code := v.GenerateCode(email)
	subject := config.ConfigManagerInstance.Config.Site.SiteName + "验证码"
	body := "你的验证码为：" + code
	err := v.smtpHelper.SendEmail(email, subject, body)
	if err != nil {
		v.mu.Lock()
		delete(v.codes, email)
		v.mu.Unlock()
	}
	v.mu.Lock()
	v.rateLimit[email]++
	v.mu.Unlock()

	return err
}
