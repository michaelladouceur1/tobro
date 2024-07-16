/* tslint:disable */
/* eslint-disable */
/**
 * tobro_http_server
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  ConnectRequest,
  ConnectResponse,
  DigitalWritePinResponse,
  ErrorResponse,
  Pong,
  SetupPinRequest,
  SetupPinResponse,
  WritePinRequest,
} from '../models/index';
import {
    ConnectRequestFromJSON,
    ConnectRequestToJSON,
    ConnectResponseFromJSON,
    ConnectResponseToJSON,
    DigitalWritePinResponseFromJSON,
    DigitalWritePinResponseToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    PongFromJSON,
    PongToJSON,
    SetupPinRequestFromJSON,
    SetupPinRequestToJSON,
    SetupPinResponseFromJSON,
    SetupPinResponseToJSON,
    WritePinRequestFromJSON,
    WritePinRequestToJSON,
} from '../models/index';

export interface ConnectPostRequest {
    connectRequest?: ConnectRequest;
}

export interface DigitalWritePinPostRequest {
    writePinRequest?: WritePinRequest;
}

export interface SetupPinPostRequest {
    setupPinRequest?: SetupPinRequest;
}

/**
 * 
 */
export class DefaultApi extends runtime.BaseAPI {

    /**
     */
    async connectPostRaw(requestParameters: ConnectPostRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ConnectResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/connect`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: ConnectRequestToJSON(requestParameters['connectRequest']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ConnectResponseFromJSON(jsonValue));
    }

    /**
     */
    async connectPost(requestParameters: ConnectPostRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ConnectResponse> {
        const response = await this.connectPostRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async digitalWritePinPostRaw(requestParameters: DigitalWritePinPostRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<DigitalWritePinResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/digital_write_pin`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: WritePinRequestToJSON(requestParameters['writePinRequest']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DigitalWritePinResponseFromJSON(jsonValue));
    }

    /**
     */
    async digitalWritePinPost(requestParameters: DigitalWritePinPostRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<DigitalWritePinResponse> {
        const response = await this.digitalWritePinPostRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async pingGetRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Pong>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/ping`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => PongFromJSON(jsonValue));
    }

    /**
     */
    async pingGet(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Pong> {
        const response = await this.pingGetRaw(initOverrides);
        return await response.value();
    }

    /**
     */
    async setupPinPostRaw(requestParameters: SetupPinPostRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<SetupPinResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/setup_pin`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: SetupPinRequestToJSON(requestParameters['setupPinRequest']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => SetupPinResponseFromJSON(jsonValue));
    }

    /**
     */
    async setupPinPost(requestParameters: SetupPinPostRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<SetupPinResponse> {
        const response = await this.setupPinPostRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
