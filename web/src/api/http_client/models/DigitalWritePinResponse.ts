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

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface DigitalWritePinResponse
 */
export interface DigitalWritePinResponse {
    /**
     * 
     * @type {number}
     * @memberof DigitalWritePinResponse
     */
    pinNumber?: number;
    /**
     * 
     * @type {number}
     * @memberof DigitalWritePinResponse
     */
    value?: number;
}

/**
 * Check if a given object implements the DigitalWritePinResponse interface.
 */
export function instanceOfDigitalWritePinResponse(value: object): value is DigitalWritePinResponse {
    return true;
}

export function DigitalWritePinResponseFromJSON(json: any): DigitalWritePinResponse {
    return DigitalWritePinResponseFromJSONTyped(json, false);
}

export function DigitalWritePinResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): DigitalWritePinResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'pinNumber': json['pinNumber'] == null ? undefined : json['pinNumber'],
        'value': json['value'] == null ? undefined : json['value'],
    };
}

export function DigitalWritePinResponseToJSON(value?: DigitalWritePinResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'pinNumber': value['pinNumber'],
        'value': value['value'],
    };
}

