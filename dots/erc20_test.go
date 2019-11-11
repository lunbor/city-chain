package dots

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestNewErc20(t *testing.T) {
	_, err := NewErc20()
	if err != nil {
		t.Error(err)
	}
}

func TestErc20_Transfer(t *testing.T) {
	inputStr := `a9059cbb000000000000000000000000af8b171c65b1cd135e475f4c18826c28a45716c60000000000000000000000000000000000000000000000000f6b75914d0770006130663239346536313639623466356261336662313065326666663634376138`
	bs, err := hex.DecodeString(inputStr)

	erc, err := NewErc20()
	input, _ := erc.Transfer(bs)

	v := big.NewInt(0)
	v.SetString("1.111111", 10)
	if input.To.Hex() != "0xaf8b171c65b1cd135e475f4c18826c28a45716c6" && input.Id != "a0f294e6169b4f5ba3fb10e2fff647a8" && input.Value.Cmp(v) != 0 {
		t.Error("")
	}
	_ = err
}
