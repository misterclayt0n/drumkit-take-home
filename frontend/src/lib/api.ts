import type { CreateLoadResponse, ListLoadsParams, Load, LoadsResponse } from './types'

const DEFAULT_API_BASE_URL = 'http://localhost:6969'

export function getApiBaseUrl() {
  return import.meta.env.VITE_API_BASE_URL ?? DEFAULT_API_BASE_URL
}

async function apiRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${getApiBaseUrl()}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
  })

  if (!response.ok) {
    const message = await readErrorMessage(response)
    throw new Error(message)
  }

  return response.json() as Promise<T>
}

export async function listLoads(params: ListLoadsParams): Promise<LoadsResponse> {
  const searchParams = new URLSearchParams({
    page: String(params.page),
    limit: String(params.limit),
  })

  if (params.status) {
    searchParams.set('status', params.status)
  }

  return apiRequest<LoadsResponse>(`/v1/loads?${searchParams.toString()}`, {
    method: 'GET',
  })
}

export async function createLoad(payload: Load): Promise<CreateLoadResponse> {
  return apiRequest<CreateLoadResponse>('/v1/integrations/webhooks/loads', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

async function readErrorMessage(response: Response) {
  try {
    const data = (await response.json()) as { error?: string }
    if (data?.error) {
      return data.error
    }
  } catch {
    // noop
  }

  return `Request failed with status ${response.status}`
}
