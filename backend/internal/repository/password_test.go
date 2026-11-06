package repository

import "testing"

func TestPasswords(t *testing.T) {
	a, e := passwordHash("demo-long-password")
	if e != nil {
		t.Fatal(e)
	}
	b, _ := passwordHash("demo-long-password")
	if a == b {
		t.Fatal("salt reused")
	}
	if !passwordMatches(a, "demo-long-password") || passwordMatches(a, "wrong") {
		t.Fatal("password verification")
	}
	for _, s := range []string{"", "plain", "$argon2id$v=19$m=999999999,t=2,p=1$a$b"} {
		if passwordMatches(s, "x") {
			t.Fatal("invalid hash accepted")
		}
	}
}
