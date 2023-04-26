package lib

// List of list and generate name
type List struct {
	Name         string
	FileType     string
	Template     string
	Location     string
	Lang         string
	ReplaceQuote bool
}

// Data for struct list of data
type Data struct {
	FileName  string
	Src       string
	GoPackage string
	Package   string
	Services  []Service
	Messages  []Message
	Enums     []*Enum
}

// Message for messages
type Message struct {
	Index    int
	Name     string
	NumField int
	IsEmpty  bool
	Fields   []*Field
	Options  []*Option
}

// Enum for messaging enum
type Enum struct {
	Name    string
	Options []*Option
}

// Service list of services
type Service struct {
	Name       string
	Methods    []*Method
	AllMessage []Message
	Elastic    bool
}

// Method list of method inside service
type Method struct {
	Name            string
	OriginalName    string
	Input           string
	Output          string
	Options         []*Option
	HttpMode        string
	URLPath         string
	ShortURL        string
	LocationPath    string
	MessageRequest  Message
	MessageResponse Message
	MessageAll      []Message
	Indent          string
	IsEmptyRequest  bool
	IsEmptyResponse bool
}

// Option for optional
type Option struct {
	Code  string
	Name  string
	Value string
}

// Field for field in messages
type Field struct {
	Index                  int
	Name                   string
	TypeData               string
	TypeDataOrigin         string
	IsRepeated             bool
	IsOptional             bool
	IsRequired             bool
	IsFieldMessage         bool
	IsAvailableDataExample bool
	DataExample            string
	MessageTo              Message `json:",omitempty"`
	MessageToName          string
}

type fileAfterExecute struct {
	Filename string
	PkgName  string
	Location string
}
