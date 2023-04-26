# install plugins
go install .

# show after installation
ls $GOPATH/bin/protoc-gen-react

# klao mau uninstall
go clean -i github.com/zokypesch/protoc-gen-react

# convention
we follow convention protoc, so the project name is protoc-gen-react
"protoc-gen" is convention, then "react" is our plugin name.
if you want use this plugin make sure --react_out=$GOPATH/src is declare in your statement.

# how to run
go install . && protoc -I proto_example example.proto --react_out=rpctarget="GrantPermission,CheckResult":$GOPATH/src -I=$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

# docs
- multiple context
    https://dev.to/alexandprivate/smoothly-test-react-components-with-multiple-contexts-453f
- axios mock
    https://blog.bitsrc.io/axios-mocking-with-reactjs-85d83d51704f

import axios, { AxiosRequestConfig } from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';
const axiosMockInstance = axios.create();
const axiosLiveInstance = axios.create();
export const axiosMockAdapterInstance= new AxiosMockAdapter(axiosMockInstance, { delayResponse: 0 });
export default process.env.isAxioMock? axiosMockInstance : axiosLiveInstance;

go install . && protoc -I proto_example example.proto --react_out=rpctarget="GrantPermission,CheckResult,GetAllTransHistoryByUser":$GOPATH/src -I=$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

go install . && protoc -I proto_example course.proto --react_out=rpctarget="GetCourseDetail":$GOPATH/src -I=$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

{{- if eq $method.HttpMode "get"}}
{{- if eq $method.Input "empty"}}
	{{ $method.Name }} = () => Api(this.props).{{ $method.HttpMode }}("{{ $method.URLPath }}")
{{- else}}
	{{ $method.Name }} = (params: {{ ucfirst $method.Input }}) => Api(this.props).{{ $method.HttpMode }}(buildParams("{{ $method.URLPath }}", params))
{{- end}}
{{- else}}
	{{ $method.Name }} = (params: {{ ucfirst $method.Input }}) => Api(this.props).{{ $method.HttpMode }}("{{ $method.URLPath }}", JSON.stringify(params))
{{- end}}
{{- end}}
{{- end}}

callApi: {{ unescape "<" }}TypeRequest, TypeResponse>() => Promise{{ unescape "<" }}AxiosResponse{{ unescape "<" }}TypeResponse>>,

# PR 
- jika dia get belakangnya sudah ada id jangan di kasih base params
- check parametr -