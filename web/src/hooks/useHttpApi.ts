import { DefaultApi, Configuration } from "../api/http_client";

export function useHttpApi() {
  const configuration = new Configuration({ basePath: "http://localhost:8080" });
  const api = new DefaultApi(configuration);
  return api;
}
