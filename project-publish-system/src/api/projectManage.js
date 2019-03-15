// 项目管理
import request from '@/utils/request'

// 获取项目管理表格数据
export function TableData(form, CurrentPage, PageSize, userId) {
  const data = {
    'project': {
      Name: form.projectName,
      WarehouseName: form.storeName
    },
    'page': {
      CurrentPage, PageSize
    },
    projectTypeId: form.projectType + '',
    userId
  }
  return request({
    url: '/users/find/project/all',
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

// 新增编辑项目
export function addModifyProjecy(form, id, UserId) {
  return request({
    url: '/users/save/project',
    method: 'post',
    data: {
      id,
      'Name': form.projectName,
      'ProjectTypeId': form.projectType,
      'ServerAddress': form.serverAddress,
      'GitAddress': form.gitAddress,
      'OnlineAccessAddress': form.onlineAddress,
      'ConfAddr': form.path,
      'WarehouseName': form.gitName,
      UserId
    }
  })
}

// 仓库名是否存在
export function checkGitName(warehouseName) {
  return request({
    url: '/users/find/warehouse/name',
    method: 'post',
    data: {
      warehouseName
    }
  })
}

// 初始化项目类型
export function projectType() {
  return request({
    url: '/users/find/project/type/all',
    method: 'get'
  })
}

// 保存环境配置
export function config(form, ProjectId, start1, start2, start3, start4, start5, start6, start7, start8, id1, id2, id3, id4, id5, id6, id7, id8) {
  return request({
    url: '/users/save/setting',
    method: 'post',
    data: [{
      'id': id1,
      ProjectId,
      'Keys': 'openUsername',
      'Values': form.gitLoginName,
      'Enable': start1 + ''
    },
    {
      'id': id2,
      ProjectId,
      'Keys': 'openPassword',
      'Values': form.gitLoginPsw,
      'Enable': start2 + ''
    },
    {
      'id': id3,
      ProjectId,
      'Keys': 'deployKey',
      'Values': form.publishKey,
      'Enable': start3 + ''
    },
    {
      'id': id4,
      ProjectId,
      'Keys': 'pullPath',
      'Values': form.pullPath,
      'Enable': start4 + ''
    },
    {
      'id': id5,
      ProjectId,
      'Keys': 'runPath',
      'Values': form.runPath,
      'Enable': start5 + ''
    },
    {
      'id': id6,
      ProjectId,
      'Keys': 'beforeCommand',
      'Values': form.publishBeforeCmd,
      'Enable': start6 + ''
    },
    {
      'id': id7,
      ProjectId,
      'Keys': 'publishCommand',
      'Values': form.publishCmd,
      'Enable': start7 + ''
    },
    {
      'id': id8,
      ProjectId,
      'Keys': 'afterCommand',
      'Values': form.publishAfterCmd,
      'Enable': start8 + ''
    }
    ]
  })
}

// 环境配置查询
export function configData(projectId) {
  return request({
    url: '/users/find/project/setting',
    method: 'post',
    data: {
      projectId
    }
  })
}
