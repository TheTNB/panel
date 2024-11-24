package bootstrap

import (
	"log"

	"github.com/go-rat/utils/crypt"

	"github.com/TheTNB/panel/internal/app"
)

func bootCrypter() {
	crypter, err := crypt.NewXChacha20Poly1305([]byte(app.Key))
	if err != nil {
		log.Fatalf("failed to create crypter: %v", err)
	}

	app.Crypter = crypter
}
