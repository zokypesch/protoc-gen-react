package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
	googlegen "github.com/zokypesch/protoc-gen-react/generator"
)

type generator struct {
	*googlegen.Generator
	reader io.Reader
	writer io.Writer
}

var checkExistBuf = make(map[string]*template.Template)

type fileGeneratorMulti func(protoFile *descriptor.FileDescriptorProto, targetRpc []string, goPkg []string) (Data, error)

type genSingle func(tpl List, data Data, protoName string, goPkg []string) (*plugin.CodeGeneratorResponse_File, string, error)

func newGenerator() *generator {
	return &generator{
		Generator: googlegen.New(),
		reader:    os.Stdin,
		writer:    os.Stdout,
	}
}

func (g *generator) generate(generateFile fileGeneratorMulti, listParam []List, genSingle genSingle) ([]fileAfterExecute, error) {
	var after []fileAfterExecute
	err := readRequest(g.reader, g.Request)
	if err != nil {
		return after, err
	}

	var rpcTarget []string
	targetRpc := false
	// targeting RPC
	param := strings.Split(g.Request.GetParameter(), ";")
	for _, v := range param {
		reg := strings.Split(v, "=")
		if len(reg) == 2 {
			if reg[0] == "rpctarget" {
				rpcTarget = strings.Split(reg[1], ",")
				break
			}
		}
	}
	if len(rpcTarget) > 0 {
		targetRpc = true
	}

	g.CommandLineParameters(g.Request.GetParameter())
	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GenerateAllFiles()

	for _, protoFile := range g.Request.ProtoFile {
		if len(protoFile.GetService()) < 1 {
			continue
		}
		newGoPkg := strings.Split(protoFile.GetOptions().GetGoPackage(), "/")

		datas, err := generateFile(protoFile, rpcTarget, newGoPkg)
		if err != nil {
			return after, err
		}

		// split into method
		var dataList []Data
		for _, v := range datas.Services {
			for _, vv := range v.Methods {
				if vv == nil {
					log.Println(rpcTarget)
					continue
				}
				if targetRpc && !stringExist(vv.OriginalName, rpcTarget) {
					continue
				}
				svc := Service(v)
				svc.Methods = []*Method{vv}
				newData := Data(datas)
				newData.Enums = genEmptyEnum()
				newData.Messages = genEmptyMessage()

				newData.Services = []Service{svc}
				dataList = append(dataList, newData)
			}
		}

		// end of split
		for _, genFile := range listParam {
			g.Reset()
			if genFile.Lang == "wrapper" {
				file, pkgName, err := genSingle(genFile, datas, *protoFile.Name, newGoPkg)
				if err != nil {
					return after, err
				}
				response := &plugin.CodeGeneratorResponse{}
				response.File = append(response.File, file)
				res, errWrite := writeResponseWithList(g.writer, response, genFile, pkgName, datas)
				if errWrite != nil {
					return after, errWrite
				}
				after = append(after, fileAfterExecute{
					Filename: res,
					PkgName:  pkgName,
					Location: genFile.Location,
				})

				continue
			}

			// gen based on datalist
			for _, v := range dataList {
				file, pkgName, err := genSingle(genFile, v, *protoFile.Name, newGoPkg)
				if err != nil {
					return after, err
				}

				response := &plugin.CodeGeneratorResponse{}
				response.File = append(response.File, file)
				res, errWrite := writeResponseWithList(g.writer, response, genFile, pkgName, v)
				if errWrite != nil {
					return after, errWrite
				}
				after = append(after, fileAfterExecute{
					Filename: res,
					PkgName:  pkgName,
					Location: genFile.Location,
				})
			}
		}
	}
	return after, nil
}

func (g *Operations) GenSingleFile(tpl List, data Data, protoName string, goPkg []string) (*plugin.CodeGeneratorResponse_File, string, error) {

	tplE, ok := checkExistBuf[tpl.Name]
	if !ok {
		// g.Generator = newGenerator()
		// mapping datas to template
		tplE = template.Must(template.New("").Funcs(
			template.FuncMap{
				"unescape":        unescape,
				"ucfirst":         ucFirst,
				"getFirstService": getFirstService,
				"upper":           strings.ToUpper,
				"toupper":         strings.ToUpper,
				"strReplaceParam": strReplaceParam,
				"ucDown":          ucDown,
			}).Parse(tpl.Template))
		checkExistBuf[tpl.Name] = tplE
	}

	buf := bytes.NewBuffer(nil)
	err := tplE.Execute(buf, data)

	if err != nil {
		return nil, "", err
	}

	// g.Generator.P(buf.String())
	// formatted := g.Generator.Bytes()
	// log.Println(string(formatted), data.Services[0].Methods[0].Name)

	finalGopkg := goPkg[len(goPkg)-1]
	// return code generator response
	return &plugin.CodeGeneratorResponse_File{
		Name: proto.String(data.Services[0].Methods[0].Name + tpl.FileType),
		// Name: proto.String(ucFirst(protoFileBaseName(protoName)) + "." + data.Services[0].Methods[0].Name + tpl.FileType),
		// Content: proto.String(string(formatted)),
		Content: proto.String(buf.String()),
	}, finalGopkg, nil
}

var generateIntegrationProto = false

func writeResponseWithList(w io.Writer, response *plugin.CodeGeneratorResponse, list List, pkgName string, datas Data) (string, error) {
	_, err := proto.Marshal(response)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal output proto")
	}

	listOrP := []string{
		pkgName,
		fmt.Sprintf("%s.%s", datas.Services[0].Methods[0].HttpMode, datas.Services[0].Methods[0].ShortURL),
	}
	fileName := response.GetFile()[0].GetName()
	location := genLocation(list.Location, listOrP...)

	// filename override
	switch list.Lang {
	case "wrapper":
		fileName = pkgName + list.FileType
	}
	// log.Fatalln(fileName, location)
	if list.Location == "" {
		return "", fmt.Errorf("location cannot be empty.")

	}
	content := response.GetFile()[0].GetContent()

	if list.ReplaceQuote {
		content = strings.Replace(content, "'", "`", -1)
	}

	// foundIndex := strings.Index(location, "%")
	// if foundIndex > -1 {
	// 	location = fmt.Sprintf(location, strings.ToLower(datas.Services[0].Name))
	// }

	if _, err := os.Stat(location); os.IsNotExist(err) {
		err := os.MkdirAll(location, 0755)
		if err != nil {
			log.Println("cannot create directory: ", err)
		}
	}

	fileDestPrefix, err := filepath.Abs(location + fileName)
	if err != nil {
		log.Println("cannot get filepath Abs: ", err)
	}

	// check do not replace
	fileDestPrefix = doOperateRename(fileDestPrefix, 0)

	d1 := []byte(content)
	err = ioutil.WriteFile(fileDestPrefix, d1, 0644)

	if err != nil {
		return "", errors.Wrap(err, "failed to write output proto")
	}

	return fileName, nil
}

func doOperateRename(fileDest string, i int) string {
	// check do not replace
	dat, errDoNot := ioutil.ReadFile(fileDest)

	if errDoNot == nil {
		foundIndexDoNot := strings.Index(string(dat), "DO_NOT_REPLACE")
		if foundIndexDoNot > -1 {
			log.Println("found DO_NOT_REPLACE skipping generate " + fileDest)
			return doOperateRename(fmt.Sprintf("v%d.%s", i+1, fileDest), i+1)
		}
	}
	return fileDest
}

func readRequest(r io.Reader, request *plugin.CodeGeneratorRequest) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "error while reading input")
	}

	if err = proto.Unmarshal(input, request); err != nil {
		return errors.Wrap(err, "error while parsing input proto")
	}

	if len(request.FileToGenerate) == 0 {
		return errors.New("no files to generate")
	}

	return nil
}
