export async function post<T,U>(url: string, data: T): Promise<U> {
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  console.log("Response", await response.json())

  return await response.json() as Promise<U>;
}