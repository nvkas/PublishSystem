// 用户管理
import request from '@/utils/request'

// 获取用户管理表格数据
export function userTableData(Name, LoginName, CurrentPage, PageSize) {
  const data = {
    'user': {
      Name, LoginName
    },
    'page': {
      CurrentPage, PageSize
    }
  }
  return request({
    url: '/users/select/user',
    method: 'post',
    data: data
  })
}

// 删除
export function deleteData(id) {
  return request({
    url: '/users/deleteid',
    method: 'post',
    data: { id }
  })
}

// 新增用户
export function addModifyUser(form, id) {
  return request({
    url: '/users/save/user',
    method: 'post',
    data: {
      'user': {
        'DepartmentName': form.department,
        'Email': form.email,
        'LoginName': form.loginName,
        'Name': form.realName,
        'Password': form.password,
        'Phone': form.phone,
        'Sex': parseInt(form.sex),
        'ProjectWorkPath': form.path,
        'UserTypeId': form.userType

      },
      'roleId': form.role,
      id
    }
  })
}

// 验证用户名是否存在
export function checkLoginName(loginName) {
  return request({
    url: '/users/select/login/name',
    method: 'post',
    data: {
      loginName
    }
  })
}

// 初始化用户角色
export function role() {
  return request({
    url: '/users/role',
    method: 'get'
  })
}
// 初始化用户类型
export function userType() {
  return request({
    url: '/users/find/user/type/all',
    method: 'get'
  })
}
