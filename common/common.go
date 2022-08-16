package common

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

func Encode(imgByte *bytes.Buffer) string {

	return base64.StdEncoding.EncodeToString(imgByte.Bytes())
}

//デコード
func Decode(str string) *bytes.Buffer {
	data, _ := base64.StdEncoding.DecodeString(str) //[]byte

	wb := new(bytes.Buffer)
	wb.Write(data)
	return wb
}

func NewUuid() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
	}
	uu := u.String()
	fmt.Println(uu)

	return uu
}
