package config

import "testing"

func TestNewAWSSession(t *testing.T) {
    _, err := NewAWSSession("development")
    if err != nil {
       t.Errorf("New AWS Session could not be created: %s", err)
	}
}
