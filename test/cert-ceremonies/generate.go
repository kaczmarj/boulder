package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/letsencrypt/boulder/cmd"
)

func genKey(path string) (string, error) {
	output, err := exec.Command("bin/ceremony", "-config", path).Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`and ID ([a-z0-9]{8})`)
	matches := re.FindSubmatch(output)
	if len(matches) != 2 {
		return "", errors.New("unexpected number of key ID matches")
	}
	return string(matches[1]), nil
}

func rewriteIntermediate(path, keyID string) (string, error) {
	config, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	tmp, err := ioutil.TempFile(os.TempDir(), "intermediate")
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	config = bytes.Replace(config, []byte("$keyID$"), []byte(keyID), 1)
	_, err = tmp.Write(config)
	if err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func genCert(path string) error {
	return err := exec.Command("bin/ceremony", "-config", path).Run()
}

func main() {
	// Generate keys
	rsaRootKeyID, err := genKey("test/cert-ceremonies/root-ceremony-rsa.yaml")
	cmd.FailOnError(err, "failed to generate root key + root cert")

	rsaIntermediateKeyID, err := genKey("test/cert-ceremonies/intermediate-key-ceremony-rsa.yaml")
	cmd.FailOnError(err, "failed to generate intermediate key")

	tmpRSAIntermediateA, err := rewriteIntermediate("test/cert-ceremonies/intermediate-ceremony-rsa-a.yaml", rsaRootKeyID)
	cmd.FailOnError(err, "failed to rewrite intermediate cert config with key ID")
	err = genCert(tmpRSAIntermediateA)
	cmd.FailOnError(err, "failed to generate intermediate cert")

	tmpRSAIntermediateB, err := rewriteIntermediate("test/cert-ceremonies/intermediate-ceremony-rsa-b.yaml", rsaRootKeyID)
	cmd.FailOnError(err, "failed to rewrite intermediate cert config with key ID")
	err = genCert(tmpRSAIntermediateA)
	cmd.FailOnError(err, "failed to generate intermediate cert")
}
