package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getLastSplitStr(s, sep string) string {
	ss := strings.Split(s, sep)
	return ss[len(ss)-1]
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	checkErr(err)

	var req = new(pluginpb.CodeGeneratorRequest)
	checkErr(proto.Unmarshal(input, req))

	plugin, err := protogen.Options{}.New(req)
	checkErr(err)

	if len(plugin.Files) == 0 {
		return
	}

	var pkgName protogen.GoPackageName

	for _, file := range plugin.Files {
		if len(file.Proto.Service) == 0 {
			continue
		}

		var buf bytes.Buffer

		pkgName = file.GoPackageName

		// 生成包头信息
		buf.Write([]byte(fmt.Sprintf(packageTpl, pkgName)))

		// 生成svc
		for _, service := range file.Proto.Service {
			buf.Write([]byte(fmt.Sprintf(svcProcessTpl,
				service.GetName(),
				service.GetName(),
				service.GetName(),
			)))

			for _, method := range service.Method {
				buf.Write([]byte(fmt.Sprintf(svcMethodTpl,
					method.GetName(),
					getLastSplitStr(method.GetInputType(), "."),
					getLastSplitStr(method.GetOutputType(), "."),
					service.GetName(),
					service.GetName(),
					method.GetName())),
				)
			}
		}

		// 生成目标文件
		filename := file.GeneratedFilenamePrefix + "_client-pool.go"
		dstFile := plugin.NewGeneratedFile(filename, ".")
		dstFile.Write(buf.Bytes())
	}

	// 生成全局client map
	var buf bytes.Buffer
	buf.Write([]byte(fmt.Sprintf(globalTpl, pkgName)))
	filename := "client-pool.go"
	dstFile := plugin.NewGeneratedFile(filename, ".")
	dstFile.Write(buf.Bytes())

	// Generate a response from our plugin and marshall as protobuf
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	checkErr(err)

	// Write the response to stdout, to be picked up by protoc
	fmt.Fprintf(os.Stdout, string(out))
}
