package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/rusldv/statefix/cryptolib"
)

func main() {
	fmt.Println("crtst")
	pubHex := "049255d872dd6d1ea8db0c0bfba7ba9c9f16314ca47713da4bc8101bc9107d94b801cd558c281fb70e7b9aae02eec9336a0cf5ae4da34a89d03d44bfe644f3fc83"
	sig := "3045022100b314ba99775468762ca04755975dfaf591b77112030e62c43fc3544f88e259800220797111331cd567b4af795e87903c673eb719b28c08d9f9bd38b1d98b84d04c73"
	data := "test"
	x, y, err := cryptolib.PublicKeyUnmarshalHex(pubHex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(x, y)
	pub := cryptolib.XYToPublicKey(x, y)
	fmt.Println(pub)
	data256 := sha256.Sum256([]byte(data))
	vf := cryptolib.VerifyDataASN1Hex(pub, sig, data256[:])
	fmt.Println(vf)
}
