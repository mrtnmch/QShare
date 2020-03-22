package qr

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

const (
	secret = "random-key"
	defaultKeySize = 1024
	shortKeySize = 64
)

func TestNew(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, shortKeySize)

	t.Run("Create instance", func(t *testing.T) {
		if got := NewGenerator(secret, privateKey); got == nil {
			t.Errorf("NewGenerator() = %v", got)
		}
	})
}

func TestDefaultGenerator_Generate(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, defaultKeySize)
	privateKey2, _ := rsa.GenerateKey(rand.Reader, shortKeySize)

	type fields struct {
		privateKey   *rsa.PrivateKey
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Key too short", fields:fields{privateKey: privateKey2,}, wantErr: true},
		{name: "Correct", fields:fields{privateKey: privateKey,}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewGenerator(secret, tt.fields.privateKey)
			got, err := generator.Generate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v", err)
				return
			}

			if !tt.wantErr && (got == nil || got.Code == "" || got.Base64Image == "") {
				t.Errorf("Generate() got %v", got)
			}
		})
	}
}

func TestDefaultGenerator_Validate(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, defaultKeySize)
	privateKey2, _ := rsa.GenerateKey(rand.Reader, defaultKeySize)

	generated, _ := NewGenerator(secret, privateKey).Generate()
	wrongSecretGenerated, _ := NewGenerator(secret + "_postfix", privateKey).Generate()

	type fields struct {
		privateKey   *rsa.PrivateKey
	}

	type args struct {
		base64code string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Empty code", fields:fields{privateKey: privateKey,}, args: args{base64code: ""}, wantErr: true},
		{name: "Messed up code", fields:fields{privateKey: privateKey,}, args: args{base64code: "xxx"}, wantErr: true},
		{name: "Wrong secret", fields:fields{privateKey: privateKey,}, args: args{base64code: wrongSecretGenerated.Code}, wantErr: true},
		{name: "Wrong key", fields:fields{privateKey: privateKey2,}, args: args{base64code: generated.Code}, wantErr: true},
		{name: "Correct", fields:fields{privateKey: privateKey,}, args: args{base64code: generated.Code}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey = tt.fields.privateKey

			generator := &DefaultGenerator{
				privateKey:   privateKey,
				publicKey:    &privateKey.PublicKey,
				secret:       secret,
			}

			if err := generator.Validate(tt.args.base64code); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}