package lib

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	googlegen "github.com/zokypesch/protoc-gen-react/generator"

	"github.com/pkg/errors"
)

// Operations struct of list data
type Operations struct {
	Datas     Data
	Generator *generator
	List      []List
}

// NewMaster for new master
func NewMaster(list []List) *Operations {
	combineUnprec()

	return &Operations{Generator: newGenerator(), List: list}
}

// Generate for generating file
func (g *Operations) Generate() ([]fileAfterExecute, error) {
	res, err := g.Generator.generate(g.generateFile, g.List, g.GenSingleFile)
	if err != nil {
		return nil, err
	}

	log.Println("generate success")
	return res, nil
}

func (g *Operations) generateFile(protoFile *descriptor.FileDescriptorProto, targetRpc []string, goPkg []string) (Data, error) {
	if protoFile.Name == nil {
		return Data{}, errors.New("missing filename")
	}
	if protoFile.GetOptions().GetGoPackage() == "" {
		return Data{}, errors.New("missing go_package")
	}
	newGoPkg := strings.Split(protoFile.GetOptions().GetGoPackage(), "/")

	// initial message
	data := Data{
		FileName:  *protoFile.Name,
		GoPackage: newGoPkg[len(newGoPkg)-1],
		Package:   protoFile.GetPackage(),
		Services:  make([]Service, len(protoFile.Service)),
		Messages:  make([]Message, len(protoFile.MessageType)),
	}

	isTargetedRpc := false
	if len(targetRpc) > 0 {
		isTargetedRpc = true
	}

	var newMessage []Message
	var enums []*Enum

	// get message in proto
	for iMType, messageType := range protoFile.MessageType {
		var newField []*Field
		messageName := messageType.GetName()

		for kMessageField, messageField := range messageType.Field {
			dataTypeOrigin := messageField.GetType().String()
			var dataType = dataTypeOrigin
			var isFieldMessage, isRepeated, isOptional, isRequiredField, isAvailableDataExample bool

			if dataTypeOrigin == "TYPE_MESSAGE" || messageField.GetTypeName() != "" {
				onTypeComb := strings.Split(messageField.GetTypeName(), ".")
				dataType = ucFirst(onTypeComb[len(onTypeComb)-1:][0])
				dataTypeOrigin = ucFirst(onTypeComb[len(onTypeComb)-1:][0])

				if dataType == "Timestamp" {
					dataType = "string"
				}

				foundIndexChar := strings.Index(dataType, "_")
				if foundIndexChar == -1 {
					isFieldMessage = true
				}

				dataTypeOrigin = dataType
			} else {
				dataType = grpcTypeToTs(dataTypeOrigin)
			}

			if "LABEL_REPEATED" == messageField.GetLabel().String() {
				isRepeated = true
				dataType += "[]"
			}
			// else if "LABEL_OPTIONAL" == messageField.GetLabel().String() {
			// 	isOptional = true
			// }

			newFieldOptions := stringToOpt(messageField.GetOptions().String())
			for _, vOptField := range newFieldOptions {
				switch getStringFromOptCode(vOptField.Code) {
				case "requiredField":
					if res, err := strconv.Atoi(vOptField.Value); err == nil && res == 1 {
						isRequiredField = true
					}
				}
			}

			isOptional = true
			if dataType == "string" || dataType == "number" {
				isOptional = false
				isAvailableDataExample = true
			}

			dataExamp := ""
			switch dataType {
			case "string":
				dataExamp = fmt.Sprintf(`example %s`, messageField.GetJsonName())
			case "number":
				dataExamp = "99"
			}

			newField = append(newField, &Field{
				Index:                  kMessageField + 1,
				Name:                   messageField.GetJsonName(),
				TypeData:               dataType,
				TypeDataOrigin:         dataTypeOrigin,
				IsRepeated:             isRepeated,
				IsOptional:             isOptional,
				IsFieldMessage:         isFieldMessage,
				IsRequired:             isRequiredField,
				IsAvailableDataExample: isAvailableDataExample,
				DataExample:            dataExamp,
			})
		}

		// get enum declare in message protofile
		for _, typEnum := range messageType.GetEnumType() {
			listOptEnum := make([]*Option, len(typEnum.GetValue()))

			for kValEnum, valEnum := range typEnum.GetValue() {
				listOptEnum[kValEnum] = &Option{
					Name:  valEnum.GetName(),
					Code:  valEnum.GetName(),
					Value: strconv.Itoa(int(valEnum.GetNumber())),
				}
			}
			enums = append(enums, &Enum{
				Name:    ucDown(typEnum.GetName()),
				Options: listOptEnum,
			})
		}

		newMessage = append(newMessage, Message{
			Index:  iMType + 1,
			Name:   ucFirst(messageName),
			Fields: newField,
		})
	}

	// rebuild message
	for _, vNewMessage := range newMessage {
		for _, vNewFields := range vNewMessage.Fields {
			if vNewFields.IsFieldMessage == true {
				resNewField, foundNewField := findMessage(vNewFields.TypeDataOrigin, newMessage)
				if foundNewField {
					vNewFields.MessageTo = resNewField
					vNewFields.MessageToName = ucFirst(resNewField.Name)
				}
			}
		}
	}

	// get enum ini protofile
	for _, typEnum := range protoFile.GetEnumType() {
		listOptEnum := make([]*Option, len(typEnum.GetValue()))

		for kValEnum, valEnum := range typEnum.GetValue() {
			listOptEnum[kValEnum] = &Option{
				Name:  valEnum.GetName(),
				Code:  valEnum.GetName(),
				Value: strconv.Itoa(int(valEnum.GetNumber())),
			}
		}
		enums = append(enums, &Enum{
			Name:    ucDown(typEnum.GetName()),
			Options: listOptEnum,
		})
	}

	data.Messages = newMessage
	data.Enums = enums

	// get services
	for kSvc, svc := range protoFile.Service {
		var methods []*Method
		// methods := make([]*Method, len(svc.Method))
		for kMethod, method := range svc.Method {
			if isTargetedRpc {
				if !stringExist(method.GetName(), targetRpc) {
					// skip this process
					continue
				}
			}
			methOpt := method.GetOptions().String()
			methOpt = replAscii(methOpt)
			newOptions := stringToOpt(methOpt)
			var urlPath, httpMode string
			// get options in service
			for _, vOpt := range newOptions {
				switch getStringFromOptCode(vOpt.Code) {
				case "urlPath":
					urlPath = getHttpUrl(vOpt.Value)
				case "httpMode":
					httpMode = vOpt.Value
				}
			}

			onInputMethod := strings.Split(method.GetInputType(), ".")
			typeInputMethod := ucDown(onInputMethod[len(onInputMethod)-1:][0])
			onOutputMethod := strings.Split(method.GetOutputType(), ".")
			typeOutputMethod := ucDown(onOutputMethod[len(onOutputMethod)-1:][0])

			inputMessage := getMessageByName(typeInputMethod, data.Messages)
			outputMessage := getMessageByName(typeOutputMethod, data.Messages)

			// getting short url
			splURL := strings.Split(urlPath, "/")
			shortURL := splURL[0]
			if len(splURL) > 1 {
				shortURL = fmt.Sprintf("%s.%s", splURL[len(splURL)-2], splURL[len(splURL)-1])
			}

			checkAllMsg := make(map[string]Message)
			if typeInputMethod != "empty" {
				checkAllMsg[inputMessage.Name] = inputMessage
			}
			if typeOutputMethod != "empty" {
				checkAllMsg[outputMessage.Name] = outputMessage

			}

			for _, m := range checkAllMsg {
				for _, vNewFields := range m.Fields {
					if vNewFields.IsFieldMessage == true {
						resNewField, foundNewField := findMessage(vNewFields.TypeDataOrigin, newMessage)
						if foundNewField {
							_, ok := checkAllMsg[resNewField.Name]
							if !ok {
								checkAllMsg[resNewField.Name] = resNewField
							}
						}
					}
				}
			}
			var allMsg []Message
			for _, m := range checkAllMsg {
				allMsg = append(allMsg, m)
			}

			listOrP := []string{
				goPkg[len(goPkg)-1],
				fmt.Sprintf("%s.%s", httpMode, shortURL),
			}
			location := genLocation("./%s/%s/", listOrP...)
			location = strings.Replace(location, fmt.Sprintf("%s/", *svc.Name), "", -1)

			isEmptyInput, isEmptyOutput := false, false
			if typeInputMethod == "empty" {
				log.Println("empty input")
				isEmptyInput = true
			}

			if typeOutputMethod == "empty" {
				isEmptyOutput = true
			}

			ind := strings.Repeat(`h\t `, kMethod+1)
			// ind = `"` + ind + `"`

			methods = append(methods, &Method{
				Name:            ucFirst(*method.Name),
				OriginalName:    method.GetName(),
				Input:           typeInputMethod,
				Output:          typeOutputMethod,
				HttpMode:        httpMode,
				URLPath:         urlPath,
				MessageRequest:  inputMessage,
				MessageResponse: outputMessage,
				ShortURL:        shortURL,
				MessageAll:      allMsg,
				LocationPath:    location[:len(location)-1],
				IsEmptyRequest:  isEmptyInput,
				IsEmptyResponse: isEmptyOutput,
				Indent:          ind,
			})
		}
		// put service in datas
		data.Services[kSvc] = Service{
			Name:       googlegen.CamelCase(svc.GetName()),
			Methods:    methods,
			AllMessage: newMessage,
		}
	}

	return data, nil
}
