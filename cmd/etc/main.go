package main

import (
	"fmt"

	"go.quinn.io/dataq/rpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	b, err := protojson.Marshal(&rpc.ExtractRequest{
		Oauth: &rpc.OAuth2{
			Config: &rpc.OAuth2_Config{
				Endpoint: &rpc.OAuth2_Endpoint{
					AuthStyle: rpc.OAuth2_AUTH_STYLE_UNSPECIFIED,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
