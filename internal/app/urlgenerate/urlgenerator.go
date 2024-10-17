package urlgenerate

import (
	"math/rand/v2"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/utils"
)

type URLGenerator interface {
	GenerateURL(originalURL string) string
}

type URLGeneratorImpl struct {
	URLGenerator
}

func CreateURLGenerator() *URLGeneratorImpl {
	return &URLGeneratorImpl{}
}

func (g *URLGeneratorImpl) GenerateURL(originalURL string) string {
	newUrlLength := rand.IntN(len(originalURL))
	for newUrlLength != 0 {
		newUrlLength = rand.IntN(len(originalURL))
	}
	b := "/"
	for i := 0; i < newUrlLength; i++ {
		b += string(utils.Symbols[rand.IntN(len(utils.Symbols))])
	}
	return b
}