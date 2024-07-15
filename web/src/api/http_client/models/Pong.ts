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
 * @interface Pong
 */
export interface Pong {
    /**
     * 
     * @type {string}
     * @memberof Pong
     */
    ping: string;
}

/**
 * Check if a given object implements the Pong interface.
 */
export function instanceOfPong(value: object): value is Pong {
    if (!('ping' in value) || value['ping'] === undefined) return false;
    return true;
}

export function PongFromJSON(json: any): Pong {
    return PongFromJSONTyped(json, false);
}

export function PongFromJSONTyped(json: any, ignoreDiscriminator: boolean): Pong {
    if (json == null) {
        return json;
    }
    return {
        
        'ping': json['ping'],
    };
}

export function PongToJSON(value?: Pong | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'ping': value['ping'],
    };
}

