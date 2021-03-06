package acmetool_import_jwk_account

import (
	"io/ioutil"
	"os"

	"github.com/hlandau/acme/acmetool"
	"github.com/hlandau/acme/storage"
	jose "gopkg.in/square/go-jose.v1"
)

func Register(app *acmetool.App) {
	cmd := app.CommandLine.Command("import-jwk-account", "Import a JWK account key")
	url := cmd.Arg("provider-url", "Provider URL (e.g. https://acme-v01.api.letsencrypt.org/directory)").Required().String()
	path := cmd.Arg("private-key-file", "Path to private_key.json").Required().ExistingFile()
	app.Commands["import-jwk-account"] = func(ctx acmetool.Ctx) { Main(ctx, *url, *path) }
}

func Main(ctx acmetool.Ctx, argURL, argPath string) {
	s, err := storage.NewFDB(ctx.StateDir)
	ctx.Logger.Fatale(err, "storage")

	f, err := os.Open(argPath)
	ctx.Logger.Fatale(err, "cannot open private key file")
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	ctx.Logger.Fatale(err, "cannot read file")

	k := jose.JsonWebKey{}
	err = k.UnmarshalJSON(b)
	ctx.Logger.Fatale(err, "cannot unmarshal key")

	_, err = s.ImportAccount(argURL, k.Key)
	ctx.Logger.Fatale(err, "cannot import account key")
}
