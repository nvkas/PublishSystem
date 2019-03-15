// 角色管理
import request from '@/utils/request'

// 获取用户管理表格数据
export function TableData(Name, PageSize, CurrentPage) {
  const data = {
    'role': {
      Name
    },
    'page': {
      CurrentPage, PageSize
    }
  }
  return request({
    url: '/users/find/role/all',
    method: 'post',
    data: data
  })
}

// 初始化tree
export function tree() {
  return request({
    url: '/users/menu',
    method: 'get'
  })
}

// 新增角色
export function addModifyRole(form, menuId, id) {
  return request({
    url: '/users/save/role',
    method: 'post',
    data: {
      'role': {
        'name': form.role_name,
        'remarks': form.msg
      },
      menuId,
      id
    }
  })
}

// 角色名查重
export function checkRoleName(roleName) {
  return request({
    url: '/users/find/role/name',
    method: 'POST',
    data: {
      roleName
    }
  })
}
// 根据角色ID查询对应的菜单权限iD
export function roleMenu(id) {
  return request({
    url: '/users/find/menu/id',
    method: 'POST',
    data: {
      id
    }
  })
}
