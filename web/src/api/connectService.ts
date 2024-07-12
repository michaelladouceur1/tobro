import { BaseResponse, ConnectPortRequest } from "./types";
import { post } from "./httpClient";

const BASE_URL = "http://localhost:8080"

export async function connectPort(request: ConnectPortRequest): Promise<BaseResponse> {
    return post<ConnectPortRequest, BaseResponse>(`${BASE_URL}/connect`, request)
}