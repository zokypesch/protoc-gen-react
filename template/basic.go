package template

import lib "github.com/zokypesch/protoc-gen-react/lib"

var tmplBasic = `// TANGKUBAN PERAHU
// Code generated by dayang sumbi protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}_{{ .GoPackage }}
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
import React, { useState, createContext } from 'react';
import { {{ $method.MessageResponse.Name }} } from './{{ $method.Name }}.types';

export type ContextType = {
	get: () => {{ $method.MessageResponse.Name }} | null;
	{{- if eq $method.HttpMode "get" }}
	isAlreadyFetch: boolean;
	{{- end}}
	data: {{ $method.MessageResponse.Name }} | null;
	set: (response: {{ $method.MessageResponse.Name }}) => void;
}

const default{{ $method.Name }}Values: ContextType = {
	{{- if eq $method.HttpMode "get" }}
	isAlreadyFetch: false,
	{{- end}}
	data: null,
	set: () => {},
	get: () => null
}

type props = {
	children?: React.ReactNode
};

export const {{ $method.Name }}Context =
  createContext<ContextType>(default{{ $method.Name }}Values);

export const {{ $method.Name }}Provider: React.FC<props> = ({ children }) => {
	const [data, setData] = useState{{ unescape "<" }} {{ $method.MessageResponse.Name }} | null>(null);
	{{- if eq $method.HttpMode "get" }} 
	const [isAlreadyFetch, setAlreadyFetch] = useState<boolean>(false);
	{{- end}}
	const set = (response: {{ $method.MessageResponse.Name }}) => { 
		setData(response); 
		{{- if eq $method.HttpMode "get" }} 
		setAlreadyFetch(true);
		{{- end }}
	}
	const get = () => data;

	return (
		{{ unescape "<" }}{{ $method.Name }}Context.Provider
		  value= {{ unescape "{{" }}
			{{- if eq $method.HttpMode "get" }}
			isAlreadyFetch,
			{{- end}}
			data,
			set,
			{{- if eq $method.HttpMode "get" }}
			get,
			{{- end }}
		  {{ unescape "}}" }}
		>
		  {children}
		{{ unescape "<" }}/{{ $method.Name }}Context.Provider>
	  );
};

{{- end }}
{{- end }}
`

var Basic = lib.List{
	Name:     "GenerateBasic",
	FileType: ".context.tsx",
	Template: tmplBasic,
	Location: "./%s/%s/",
	Lang:     "tsx",
}
