import { BaseResponse, ConnectPortRequest } from "./types";
import { post } from "./httpClient";

const API_URL = "http://localhost:8081"

export async function connectPort(request: ConnectPortRequest): Promise<BaseResponse> {
    return post<ConnectPortRequest, BaseResponse>(`${API_URL}/connect`, request)
}