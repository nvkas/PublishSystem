import request from '@/utils/request'

export function loginByUsername(LoginName, password) {
  const data = {
    LoginName,
    password
  }
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function getMenu(id) {
  return request({
    url: '/users/select/menu',
    method: 'post',
    data: { id }
  })
}
