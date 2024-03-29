
// Code generated by sangkuriang protoc-gen-go. DO NOT EDIT.
// source: example.proto_zolos
// File Location: api/Zolos.ts

import AxiosHttpClient, { requestConfig } from '@prakerja/core-fe/build/core/axios.http.client'
import { BuildGetString } from '@prakerja/core-fe/build/core/utils.format'
import {
	GrantPermissionRequest,
	GrantPermissionResponse,
} from './GrantPermission.types'

const baseURL = process.env.REACT_APP_API_URL;

export const GrantPermission = (
	params: GrantPermissionRequest,
	requestConfig: requestConfig,
) => {
	let buildString = BuildGetString(params)
	return AxiosHttpClient.get<undefined, GrantPermissionResponse>(
		`${baseURL}/api/v1/facecapture/grant${buildString}`,
		undefined,
		requestConfig
	);
	
}
