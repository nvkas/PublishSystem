// 个人信息
import request from '@/utils/request'

// 初始化个人信息
export function person(id) {
  return request({
    url: 'users/user/info',
    method: 'post',
    data: { id }
  })
}
// 修改个人信息
export function modify(ID, form) {
  return request({
    url: 'users/update/user',
    method: 'post',
    data: {
      ID,
      Name: form.name,
      Phone: form.phone,
      DepartmentName: form.company,
      Email: form.email,
      Birthday: form.birth,
      Address: form.address,
      Sex: parseInt(form.sex)
    }
  })
}
// 修改密码
export function modifyPsw(ID, form) {
  return request({
    url: 'users/update/pwd',
    method: 'post',
    data: {
      ID,
      OldPassword: form.oldPsw,
      Password: form.newPsw
    }
  })
}

