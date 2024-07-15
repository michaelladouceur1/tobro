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
 * @interface WritePinRequest
 */
export interface WritePinRequest {
    /**
     * 
     * @type {number}
     * @memberof WritePinRequest
     */
    pin: number;
    /**
     * 
     * @type {number}
     * @memberof WritePinRequest
     */
    value: number;
}

/**
 * Check if a given object implements the WritePinRequest interface.
 */
export function instanceOfWritePinRequest(value: object): value is WritePinRequest {
    if (!('pin' in value) || value['pin'] === undefined) return false;
    if (!('value' in value) || value['value'] === undefined) return false;
    return true;
}

export function WritePinRequestFromJSON(json: any): WritePinRequest {
    return WritePinRequestFromJSONTyped(json, false);
}

export function WritePinRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): WritePinRequest {
    if (json == null) {
        return json;
    }
    return {
        
        'pin': json['pin'],
        'value': json['value'],
    };
}

export function WritePinRequestToJSON(value?: WritePinRequest | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'pin': value['pin'],
        'value': value['value'],
    };
}

