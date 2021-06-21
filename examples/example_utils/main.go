package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
)

const (
	//! WARNING: please, DO NOT use this key downloading in production apps, THIS IS ABSOLUTELY INSECURE!
	//! I mean, seriously, this way used just for examples, we can't create most secured app just for
	//! these examples
	publicKeysForExamplesURL = "https://git.io/JtImk"
)

func ReadWarningsToStdErr(err chan error) {
	go func() {
		for {
			pp.Fprintln(os.Stderr, <-err)
		}
	}()
}

func PrepareAppStorageForExamples() (appStoragePath string) {
	appStorage, err := GetAppStorage("mtproto-example", NamespaceUser)
	dry.PanicIfErr(err)

	if !dry.FileExists(appStorage) {
		if !dry.PathIsWritable(appStorage) {
			fmt.Printf("cant create app local storage at %v\n", appStorage)
			os.Exit(1)
		}
		err := os.MkdirAll(appStorage, 0755)
		dry.PanicIfErr(err)
	}
	publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")
	if !dry.FileExists(publicKeys) {
		fmt.Printf("Downloading public keys from %v, this can be insecure\n", publicKeysForExamplesURL)

		resp, err := http.Get(publicKeysForExamplesURL)
		dry.PanicIfErr(errors.Wrap(err, "can't download public keys"))
		defer resp.Body.Close()

		out, err := os.Create(publicKeys)
		dry.PanicIfErr(errors.Wrap(err, "can't download public keys"))

		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		dry.PanicIfErr(errors.Wrap(err, "can't download public keys"))
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
