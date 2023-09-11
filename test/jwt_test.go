package test

import (
	"testing"
	"time"

	"github.com/Milou666/Mitter/util"
)

func TestVerifyJWTWhenValid(t *testing.T) {
	payload := map[string]any{
		"val1": "val",
		"val2": float64(2), // int are parsed into float 64 by jwt
	}
	expiration := time.Now().Add(time.Hour * 24)
	jwt := util.GenJWT(expiration, payload)

	t.Logf("token  = %q", jwt)

	clear_payload, err := util.VerifyJWT(jwt)

	if err != nil {
		t.Errorf("shouldn't throw err %v", err)
	}

	if payload["val1"] != clear_payload["val1"] || payload["val2"] != clear_payload["val2"] {
		t.Errorf("different payload %q != %q or %d != %d", payload["val1"], clear_payload["val1"], payload["val2"], clear_payload["val2"])
	}
}

func TestVerifyJWTWhenExpired(t *testing.T) {
	payload := map[string]any{
		"val1": "val",
		"val2": float64(2), // int are parsed into float 64 by jwt
	}
	expiration := time.Now().Add(time.Second)
	jwt := util.GenJWT(expiration, payload)

	t.Logf("token  = %q", jwt)
	time.Sleep(2 * time.Second)

	_, err := util.VerifyJWT(jwt)

	if err == nil {
		t.Errorf("should throw an error")
	}

	if err.Error() != "Token is expired" {
		t.Errorf("err is %v but should be 'Token is expired'", err)
	}
}
