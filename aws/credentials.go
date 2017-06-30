package aws

import (
	"os"
	"path/filepath"
	"fmt"
	"errors"
)

func promptCredentials() (string, string, error) {
	var access, secret string

	fmt.Println()
	fmt.Print("AWS Access Key ID? ")
	if _, err := fmt.Scanln(&access);  err != nil {
		return access, secret, err
	}
	fmt.Print("AWS Secret Access Key? ")
	if _, err := fmt.Scanln(&secret);  err != nil {
		return access, secret, err
	}

	return access, secret, nil
}

func storeAWSCredentials(accessKey, secretAccesskey, profile string) (string, error) {
	dir :=  filepath.Join(os.Getenv("HOME"), ".aws")
	if _, err := os.Stat(dir); os.IsNotExist(err){
		if err := os.Mkdir(dir, 0700); err != nil {
			return "", fmt.Errorf("creating .aws dir: %s", err)
		}
	}

	filepath := filepath.Join(dir, "credentials")
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return filepath, fmt.Errorf("appending to '%s': %s", filepath, err)
	}

	if secretAccesskey == "" {
		return filepath, errors.New("given empty secret access key")
	}
	if accessKey == "" {
		return filepath, errors.New("given empty access key")
	}
	if profile == "" {
		 profile = "default"
	}

	if _, err := fmt.Fprintf(f, "[%s]\naws_access_key_id = %s\naws_secret_access_key = %s\n", profile, accessKey, secretAccesskey); err != nil {
		return filepath, err
	}

	return filepath, nil
}
