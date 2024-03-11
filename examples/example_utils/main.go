package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	// WARNING: please, DO NOT use this key downloading in production apps,
	// THIS IS ABSOLUTELY INSECURE! I mean, seriously, this way used just for
	// examples, we can't create most secured app just for these examples
	publicKeysForExamplesURL = "https://git.io/JtImk"
)

func PrepareAppStorageForExamples() (appStoragePath string) {
	appStorage, err := GetAppStorage("mtproto-example", NamespaceUser)
	check(err)

	if !fileExists(appStorage) {
		err := os.MkdirAll(appStorage, 0755)
		check(err)
	}
	publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")
	if !fileExists(publicKeys) {
		fmt.Printf("Downloading public keys from %v, this can be insecure\n", publicKeysForExamplesURL)

		resp, err := http.Get(publicKeysForExamplesURL)
		check(errors.Wrap(err, "can't download public keys"))
		defer resp.Body.Close()

		out, err := os.Create(publicKeys)
		check(errors.Wrap(err, "can't download public keys"))

		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		check(errors.Wrap(err, "can't download public keys"))
	}

	return appStorage
}

type Namespace uint8

const (
	NamespaceUnknown Namespace = iota
	NamespaceGlobal
	NamespaceUser
	NamespaceDirectory
)

func GetAppStorage(appName string, namespace Namespace) (string, error) {
	switch namespace {
	case NamespaceGlobal:
		return filepath.Join("var", "lib", appName), nil
	case NamespaceUser:
		p, err := GetAppStorage(appName, NamespaceGlobal)
		if err != nil {
			return "", err
		}
		u, _ := user.Current()
		userPath, err := GetUserNamespaceDir(u.Username)
		if err != nil {
			return "", err
		}
		return filepath.Join(userPath, p), nil
	default:
		return "", errors.New("Incompatible feature for this namespace")
	}
}

func GetUserNamespaceDir(username string) (string, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return "", errors.Wrapf(err, "looking up '%v'", username)
	}

	return filepath.Join(u.HomeDir, ".local"), nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
