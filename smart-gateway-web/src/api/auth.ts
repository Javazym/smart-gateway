import { post } from './http'
import type { LoginParams, LoginResult } from '@/types'

export const login = (params: LoginParams) =>
  post<LoginResult>('/gateway_auth/login', params)
