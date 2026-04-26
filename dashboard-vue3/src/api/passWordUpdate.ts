import request from '@/utils/http'

export function passWordUpdate(params: {
  oldPassword: string
  newPassword: string
}) {
  return request.put<{ success: boolean }>({
    url: '/user/passwordUpdate',
    data: params
  })
}