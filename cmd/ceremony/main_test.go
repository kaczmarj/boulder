package main

import "testing"

func TestKeyGenConfigValidate(t *testing.T) {
	cases := []struct {
		name          string
		config        keyGenConfig
		expectedError string
	}{
		{
			name:          "no key.type",
			config:        keyGenConfig{},
			expectedError: "key.type is required",
		},
		{
			name: "bad key.type",
			config: keyGenConfig{
				Type: "doop",
			},
			expectedError: "key.type can only be 'rsa' or 'ecdsa'",
		},
		{
			name: "bad key.rsa-mod-length",
			config: keyGenConfig{
				Type:         "rsa",
				RSAModLength: 1337,
			},
			expectedError: "key.rsa-mod-length can only be 2048 or 4096",
		},
		{
			name: "key.type is rsa but key.ecdsa-curve is present",
			config: keyGenConfig{
				Type:         "rsa",
				RSAModLength: 2048,
				ECDSACurve:   "bad",
			},
			expectedError: "if key.type = 'rsa' then key.ecdsa-curve is not used",
		},
		{
			name: "bad key.ecdsa-curve",
			config: keyGenConfig{
				Type:       "ecdsa",
				ECDSACurve: "bad",
			},
			expectedError: "key.ecdsa-curve can only be 'P-224', 'P-256', 'P-384', or 'P-521'",
		},
		{
			name: "key.type is ecdsa but key.rsa-mod-length is present",
			config: keyGenConfig{
				Type:         "ecdsa",
				RSAModLength: 2048,
				ECDSACurve:   "P-256",
			},
			expectedError: "if key.type = 'ecdsa' then key.rsa-mod-length is not used",
		},
		{
			name: "good rsa config",
			config: keyGenConfig{
				Type:         "rsa",
				RSAModLength: 2048,
			},
		},
		{
			name: "good ecdsa config",
			config: keyGenConfig{
				Type:       "ecdsa",
				ECDSACurve: "P-256",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.validate()
			if err != nil && err.Error() != tc.expectedError {
				t.Fatalf("Unexpected error, wanted: %q, got: %q", tc.expectedError, err)
			} else if err == nil && tc.expectedError != "" {
				t.Fatalf("validate didn't fail, wanted: %q", err)
			}
		})
	}
}

func TestRootConfigValidate(t *testing.T) {
	cases := []struct {
		name          string
		config        rootConfig
		expectedError string
	}{
		{
			name:          "no pkcs11.module",
			config:        rootConfig{},
			expectedError: "pkcs11.module is required",
		},
		{
			name: "no pkcs11.store-key-with-label",
			config: rootConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module: "module",
				},
			},
			expectedError: "pkcs11.store-key-with-label is required",
		},
		{
			name: "bad key fields",
			config: rootConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
			},
			expectedError: "key.type is required",
		},
		{
			name: "no outputs.public-key-path",
			config: rootConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
				Key: keyGenConfig{
					Type:         "rsa",
					RSAModLength: 2048,
				},
			},
			expectedError: "outputs.public-key-path is required",
		},
		{
			name: "no outputs.certificate-path",
			config: rootConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
				Key: keyGenConfig{
					Type:         "rsa",
					RSAModLength: 2048,
				},
				Outputs: struct {
					PublicKeyPath   string `yaml:"public-key-path"`
					CertificatePath string `yaml:"certificate-path"`
				}{
					PublicKeyPath: "path",
				},
			},
			expectedError: "outputs.certificate-path is required",
		},
		{
			name: "bad certificate-profile",
			config: rootConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
				Key: keyGenConfig{
					Type:         "rsa",
					RSAModLength: 2048,
				},
				Outputs: struct {
					PublicKeyPath   string `yaml:"public-key-path"`
					CertificatePath string `yaml:"certificate-path"`
				}{
					PublicKeyPath:   "path",
					CertificatePath: "path",
				},
			},
			expectedError: "not-before is required",
		},
		{
			name: "good config",
			config: rootConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
				Key: keyGenConfig{
					Type:         "rsa",
					RSAModLength: 2048,
				},
				Outputs: struct {
					PublicKeyPath   string `yaml:"public-key-path"`
					CertificatePath string `yaml:"certificate-path"`
				}{
					PublicKeyPath:   "path",
					CertificatePath: "path",
				},
				CertProfile: certProfile{
					NotBefore:          "a",
					NotAfter:           "b",
					SignatureAlgorithm: "c",
					CommonName:         "d",
					Organization:       "e",
					Country:            "f",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.validate()
			if err != nil && err.Error() != tc.expectedError {
				t.Fatalf("Unexpected error, wanted: %q, got: %q", tc.expectedError, err)
			} else if err == nil && tc.expectedError != "" {
				t.Fatalf("validate didn't fail, wanted: %q", err)
			}
		})
	}
}

func TestIntermediateConfigValidate(t *testing.T) {
	cases := []struct {
		name          string
		config        intermediateConfig
		expectedError string
	}{
		{
			name:          "no pkcs11.module",
			config:        intermediateConfig{},
			expectedError: "pkcs11.module is required",
		},
		{
			name: "no pkcs11.signing-key-label",
			config: intermediateConfig{
				PKCS11: struct {
					Module       string `yaml:"module"`
					PIN          string `yaml:"pin"`
					SigningSlot  uint   `yaml:"signing-key-slot"`
					SigningLabel string `yaml:"signing-key-label"`
					SigningKeyID string `yaml:"signing-key-id"`
				}{
					Module: "module",
				},
			},
			expectedError: "pkcs11.signing-key-label is required",
		},
		{
			name: "no pkcs11.key-id",
			config: intermediateConfig{
				PKCS11: struct {
					Module       string `yaml:"module"`
					PIN          string `yaml:"pin"`
					SigningSlot  uint   `yaml:"signing-key-slot"`
					SigningLabel string `yaml:"signing-key-label"`
					SigningKeyID string `yaml:"signing-key-id"`
				}{
					Module:       "module",
					SigningLabel: "label",
				},
			},
			expectedError: "pkcs11.signing-key-id is required",
		},
		{
			name: "no inputs.public-key-path",
			config: intermediateConfig{
				PKCS11: struct {
					Module       string `yaml:"module"`
					PIN          string `yaml:"pin"`
					SigningSlot  uint   `yaml:"signing-key-slot"`
					SigningLabel string `yaml:"signing-key-label"`
					SigningKeyID string `yaml:"signing-key-id"`
				}{
					Module:       "module",
					SigningLabel: "label",
					SigningKeyID: "id",
				},
			},
			expectedError: "inputs.public-key-path is required",
		},
		{
			name: "no inputs.issuer-certificate-path",
			config: intermediateConfig{
				PKCS11: struct {
					Module       string `yaml:"module"`
					PIN          string `yaml:"pin"`
					SigningSlot  uint   `yaml:"signing-key-slot"`
					SigningLabel string `yaml:"signing-key-label"`
					SigningKeyID string `yaml:"signing-key-id"`
				}{
					Module:       "module",
					SigningLabel: "label",
					SigningKeyID: "id",
				},
				Inputs: struct {
					PublicKeyPath         string `yaml:"public-key-path"`
					IssuerCertificatePath string `yaml:"issuer-certificate-path"`
				}{
					PublicKeyPath: "path",
				},
			},
			expectedError: "inputs.issuer-certificate is required",
		},
		{
			name: "no outputs.certificate-path",
			config: intermediateConfig{
				PKCS11: struct {
					Module       string `yaml:"module"`
					PIN          string `yaml:"pin"`
					SigningSlot  uint   `yaml:"signing-key-slot"`
					SigningLabel string `yaml:"signing-key-label"`
					SigningKeyID string `yaml:"signing-key-id"`
				}{
					Module:       "module",
					SigningLabel: "label",
					SigningKeyID: "id",
				},
				Inputs: struct {
					PublicKeyPath         string `yaml:"public-key-path"`
					IssuerCertificatePath string `yaml:"issuer-certificate-path"`
				}{
					PublicKeyPath:         "path",
					IssuerCertificatePath: "path",
				},
			},
			expectedError: "outputs.certificate-path is required",
		},
		{
			name: "bad certificate-profile",
			config: intermediateConfig{
				PKCS11: struct {
					Module       string `yaml:"module"`
					PIN          string `yaml:"pin"`
					SigningSlot  uint   `yaml:"signing-key-slot"`
					SigningLabel string `yaml:"signing-key-label"`
					SigningKeyID string `yaml:"signing-key-id"`
				}{
					Module:       "module",
					SigningLabel: "label",
					SigningKeyID: "id",
				},
				Inputs: struct {
					PublicKeyPath         string `yaml:"public-key-path"`
					IssuerCertificatePath string `yaml:"issuer-certificate-path"`
				}{
					PublicKeyPath:         "path",
					IssuerCertificatePath: "path",
				},
				Outputs: struct {
					CertificatePath string `yaml:"certificate-path"`
				}{
					CertificatePath: "path",
				},
			},
			expectedError: "not-before is required",
		},
		{
			name: "good config",
			config: intermediateConfig{
				PKCS11: struct {
					Module       string `yaml:"module"`
					PIN          string `yaml:"pin"`
					SigningSlot  uint   `yaml:"signing-key-slot"`
					SigningLabel string `yaml:"signing-key-label"`
					SigningKeyID string `yaml:"signing-key-id"`
				}{
					Module:       "module",
					SigningLabel: "label",
					SigningKeyID: "id",
				},
				Inputs: struct {
					PublicKeyPath         string `yaml:"public-key-path"`
					IssuerCertificatePath string `yaml:"issuer-certificate-path"`
				}{
					PublicKeyPath:         "path",
					IssuerCertificatePath: "path",
				},
				Outputs: struct {
					CertificatePath string `yaml:"certificate-path"`
				}{
					CertificatePath: "path",
				},
				CertProfile: certProfile{
					NotBefore:          "a",
					NotAfter:           "b",
					SignatureAlgorithm: "c",
					CommonName:         "d",
					Organization:       "e",
					Country:            "f",
					OCSPURL:            "g",
					CRLURL:             "h",
					IssuerURL:          "i",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.validate(intermediateCert)
			if err != nil && err.Error() != tc.expectedError {
				t.Fatalf("Unexpected error, wanted: %q, got: %q", tc.expectedError, err)
			} else if err == nil && tc.expectedError != "" {
				t.Fatalf("validate didn't fail, wanted: %q", err)
			}
		})
	}
}

func TestKeyConfigValidate(t *testing.T) {
	cases := []struct {
		name          string
		config        keyConfig
		expectedError string
	}{
		{
			name:          "no pkcs11.module",
			config:        keyConfig{},
			expectedError: "pkcs11.module is required",
		},
		{
			name: "no pkcs11.store-key-with-label",
			config: keyConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module: "module",
				},
			},
			expectedError: "pkcs11.store-key-with-label is required",
		},
		{
			name: "bad key fields",
			config: keyConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
			},
			expectedError: "key.type is required",
		},
		{
			name: "no outputs.public-key-path",
			config: keyConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
				Key: keyGenConfig{
					Type:         "rsa",
					RSAModLength: 2048,
				},
			},
			expectedError: "outputs.public-key-path is required",
		},
		{
			name: "good config",
			config: keyConfig{
				PKCS11: struct {
					Module     string `yaml:"module"`
					PIN        string `yaml:"pin"`
					StoreSlot  uint   `yaml:"store-key-in-slot"`
					StoreLabel string `yaml:"store-key-with-label"`
				}{
					Module:     "module",
					StoreLabel: "label",
				},
				Key: keyGenConfig{
					Type:         "rsa",
					RSAModLength: 2048,
				},
				Outputs: struct {
					PublicKeyPath string `yaml:"public-key-path"`
				}{
					PublicKeyPath: "path",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.validate()
			if err != nil && err.Error() != tc.expectedError {
				t.Fatalf("Unexpected error, wanted: %q, got: %q", tc.expectedError, err)
			} else if err == nil && tc.expectedError != "" {
				t.Fatalf("validate didn't fail, wanted: %q", err)
			}
		})
	}
}
