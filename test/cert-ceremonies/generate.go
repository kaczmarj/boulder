package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

func genKey(path string) (string, error) {
	out, err := exec.Command("bin/ceremony", "-config", path).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return "", err
	}
	re := regexp.MustCompile(`and ID ([a-z0-9]{8})`)
	matches := re.FindSubmatch(out)
	if len(matches) != 2 {
		return "", errors.New("unexpected number of key ID matches")
	}
	fmt.Println(string(out))
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
	out, err := exec.Command("bin/ceremony", "-config", path).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func main() {
	// Generate keys
	rsaRootKeyID, err := genKey("test/cert-ceremonies/root-ceremony-rsa.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(rsaRootKeyID)

	rsaIntermediateKeyID, err := genKey("test/cert-ceremonies/intermediate-key-ceremony-rsa.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(rsaIntermediateKeyID)

	tmpRSAIntermediate, err := rewriteIntermediate("test/cert-ceremonies/intermediate-ceremony-rsa.yaml", rsaRootKeyID)
	if err != nil {
		panic(err)
	}

	err = genCert(tmpRSAIntermediate)
	if err != nil {
		panic(err)
	}
}
