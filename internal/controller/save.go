package controller

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
)

func SaveFile(code [][]string) {
	codeString := bytes.NewBufferString("")
	for index, v := range code {
		if index == 0 {
			continue
		}
		codeString.WriteString(v[0])
		codeString.WriteString("\r\n")

	}
	f, err := os.OpenFile("code.txt", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		logrus.Error("open file error :", err)
		return
	}
	// 关闭文件
	defer f.Close()
	_, err = f.WriteString(codeString.String())
	if err != nil {
		logrus.Error("write file error :", err)
		return
	}
}
