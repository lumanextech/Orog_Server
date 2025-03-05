package util_test

import (
	"fmt"
	"testing"

	"github.com/simance-ai/smdx/rpcx/ws/internal/util"
)

func TestParseAccountAuthClaims(t *testing.T) {

}

func TestSignAccountAuthClaims(t *testing.T) {
	claims, err := util.SignAccountAuthClaims(10094, "Lbhi08lqB8k7bdKLKFsSyZwPygIOvwhX", 1000000)
	if err != nil {
		t.Fatal(err)
	}

	claims2, err := util.SignAccountAuthClaims(10094, "Lbhi08lqB8k7bdKLKFsSyZwPygIOvwhX", 1000000)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(claims)
	fmt.Println(claims2)
}
