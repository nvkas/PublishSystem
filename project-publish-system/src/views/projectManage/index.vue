<template>
  <div class="app-container">
    <el-form ref="form" :model="form" label-width="80px">
      <el-row class="search-text">
        <el-col :span="5">
          <el-form-item label="项目名:" prop="projectName">
            <el-input
              v-model="form.projectName"
              placeholder="项目名"
              clearable
              class="wid100"/>
          </el-form-item>
        </el-col>
        <el-col :span="5">
          <el-form-item label="项目类型:" prop="projectType">
            <el-select v-model="form.projectType" clearable placeholder="请选择" class="wid100">
              <el-option
                v-for="item in projectOptions"
                :key="item.ID"
                :label="item.Name"
                :value="item.ID"/>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="5">
          <el-form-item label="仓库名:" prop="storeName">
            <el-input
              v-model="form.storeName"
              placeholder="仓库名"
              clearable
              class="wid100"/>
          </el-form-item>
        </el-col>
        <el-col :span="2" class="text-r">
          <el-button type="primary" size="small" @click="handleSearch">查询</el-button>
        </el-col>
      </el-row>
    </el-form>
    <div class="mb16 over">
      <el-button type="primary" size="small" @click="handleDialog(1)">新增</el-button>
      <el-button :disabled="disabled" type="primary" size="small" @click="handleDialog(2)">编辑</el-button>
      <el-button :disabled="disabled" type="danger" size="small" @click="handleDeleteData">删除</el-button>
      <el-button :disabled="disabled" type="primary" size="small" class="fr" @click="handleConfigDialog">环境配置</el-button>
    </div>
    <div class="demo-block">
      <el-table
        v-loading="loading"
        ref="multipleTable"
        :data="tableData"
        border
        tooltip-effect="dark"
        style="width: 100%"
        @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"/>
        <el-table-column
          prop="ID"
          label="ID"
          width="80"/>
        <el-table-column
          prop="Name"
          label="项目名"
          show-overflow-tooltip/>
        <el-table-column
          prop="ProjectType.Name"
          label="项目类型"
          width="120"/>
        <el-table-column
          prop="WarehouseName"
          label="项目仓库名"
          show-overflow-tooltip/>
        <el-table-column
          prop="GitAddress"
          label="项目Git服务器地址"
          show-overflow-tooltip/>
        <el-table-column
          prop="ServerAddress"
          label="项目服务器部署地址"
          show-overflow-tooltip/>
        <el-table-column
          prop="Remarks"
          label="备注"
          show-overflow-tooltip/>
      </el-table>
    </div>
    <div class="pad-tb20">
      <el-pagination
        v-if="page"
        :current-page="currentPage"
        :page-sizes="[10,20]"
        :page-size="pageSize"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"/>
    </div>
    <el-dialog :title="title" :visible.sync="dialogSaveVisible" :before-close="handleClose">
      <el-form ref="ruleForm" :model="ruleForm" :rules="rules" label-width="120px" class="demo-ruleForm">
        <el-row>
          <el-col :span="12">
            <el-form-item label="项目名" prop="projectName">
              <el-input v-model.trim="ruleForm.projectName"/>
            </el-form-item>
            <el-form-item label="项目类型" prop="projectType">
              <el-select v-model="ruleForm.projectType" clearable placeholder="请选择" class="wid100">
                <el-option
                  v-for="item in projectOptions"
                  :key="item.ID"
                  :label="item.Name"
                  :value="item.ID"/>
              </el-select>
            </el-form-item>
            <el-form-item label="服务器项目地址" prop="serverAddress">
              <el-input v-model.trim="ruleForm.serverAddress"/>
            </el-form-item>
            <el-form-item label="线上访问地址" prop="onlineAddress">
              <el-input v-model.trim="ruleForm.onlineAddress"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="git仓库名" prop="gitName">
              <el-input v-model.trim="ruleForm.gitName"/>
            </el-form-item>
            <el-form-item label="git仓库地址" prop="gitAddress">
              <el-input v-model.trim="ruleForm.gitAddress"/>
            </el-form-item>
            <el-form-item label="配置文件路径" prop="path">
              <el-input v-model.trim="ruleForm.path"/>
            </el-form-item>
          </el-col>
        </el-row>
        <div class="text-center">
          <el-button v-loading="loading1" type="primary" size="small" @click="handleSave">保存</el-button>
        </div>
      </el-form>
    </el-dialog>
    <el-dialog :visible.sync="dialogConfigVisible" :before-close="handleConfigClose" title="环境配置">
      <el-form v-loading="loading3" ref="configForm" :model="configForm" :rules="configRules" label-width="100px">
        <el-row>
          <el-form-item label="git账号" prop="gitLoginName">
            <el-input v-model.trim="configForm.gitLoginName" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start1"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
          <el-form-item label="git密码" prop="gitLoginPsw">
            <el-input v-model.trim="configForm.gitLoginPsw" type="password" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start2"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
          <el-form-item label="部署秘钥" prop="publishKey">
            <el-input v-model.trim="configForm.publishKey" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start3"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
          <el-form-item label="pull路径" prop="pullPath">
            <el-input v-model.trim="configForm.pullPath" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start4"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
          <el-form-item label="运行路径" prop="runPath">
            <el-input v-model.trim="configForm.runPath" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start5"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
          <el-form-item label="发布前命令" prop="publishBeforeCmd">
            <el-input v-model.trim="configForm.publishBeforeCmd" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start6"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
          <el-form-item label="发布命令" prop="publishCmd">
            <el-input v-model.trim="configForm.publishCmd" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start7"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
          <el-form-item label="发布后命令" prop="publishAfterCmd">
            <el-input v-model.trim="configForm.publishAfterCmd" class="wid60"/>
            <span class="ml24">启用：</span>
            <el-switch
              v-model="start8"
              active-color="#13ce66"
              inactive-color="#ff4949"/>
          </el-form-item>
        </el-row>
        <div class="text-center">
          <el-button v-loading="loading2" type="primary" size="small" @click="handleConfigSave">保存</el-button>
        </div>
      </el-form>
    </el-dialog>
  </div>
</template>
<script>
import { projectType, TableData, addModifyProjecy, checkGitName, configData, config } from '@/api/projectManage'
export default {
  name: 'ProjectManage',
  data() {
    const gitName = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入您的项目名'))
      } else if (this.title === '新增项目') {
        checkGitName(value).then((res) => {
          if (res.Status) {
            callback()
          } else {
            callback(new Error('项目名已存在！'))
          }
        }).catch(() => {
          return false
        })
      } else {
        callback()
      }
    }
    return {
      configId1: 0,
      configId2: 0,
      configId3: 0,
      configId4: 0,
      configId5: 0,
      configId6: 0,
      configId7: 0,
      configId8: 0,
      start1: true,
      start2: true,
      start3: true,
      start4: true,
      start5: true,
      start6: true,
      start7: true,
      start8: true,
      dialogConfigVisible: false,
      projectOptions: [],
      disabled: true,
      form: {
        projectName: '',
        projectType: '',
        storeName: ''
      },
      loading1: false,
      loading2: false,
      loading3: false,
      page: false,
      loading: false,
      title: '',
      dialogSaveVisible: false,
      total: '',
      currentPage: 1,
      pageSize: 10,
      tableData: [],
      multipleSelection: [],
      ruleForm: {
        projectName: '',
        projectType: '',
        serverAddress: '',
        gitName: '',
        gitAddress: '',
        path: '',
        onlineAddress: ''
      },
      rules: {
        projectName: [
          { required: true, message: '请输入您的项目名', trigger: 'blur' }
        ],
        projectType: [
          { required: true, message: '请选择您的项目类型', trigger: 'blur' }
        ],
        gitName: [
          { required: true, validator: gitName, trigger: 'blur' }
        ]
      },
      configForm: {
        gitLoginName: '',
        gitLoginPsw: '',
        publishKey: '',
        pullPath: '',
        runPath: '',
        publishBeforeCmd: '',
        publishCmd: '',
        publishAfterCmd: ''
      },
      configRules: {
        pullPath: [
          { required: true, message: '请输入您的项目pull路径', trigger: 'blur' }
        ],
        runPath: [
          { required: true, message: '请输入您的项目运行路径', trigger: 'blur' }
        ],
        publishBeforeCmd: [
          { required: true, message: '请输入发布前命令', trigger: 'blur' }
        ],
        publishCmd: [
          { required: true, message: '请输入发布命令', trigger: 'blur' }
        ],
        publishAfterCmd: [
          { required: true, message: '请输入发布后命令', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.initTableData()
    this.initProjectType()
  },
  methods: {
    // 处理环境配置保存
    handleConfigSave() {
      this.$refs.configForm.validate((valid) => {
        if (valid) {
          this.loading2 = true
          config(this.configForm, this.multipleSelection[0].ID, this.start1, this.start2, this.start3, this.start4, this.start5, this.start6, this.start7, this.start8, this.configId1, this.configId2, this.configId3, this.configId4, this.configId5, this.configId6, this.configId7, this.configId8).then(res => {
            this.loading2 = false
            this.$refs.configForm.resetFields()
            if (res.Status) {
              this.$message({
                message: '配置成功',
                type: 'success'
              })
              this.dialogConfigVisible = false
            } else {
              this.$message({
                message: '配置失败，请联系系统管理员！',
                type: 'warning'
              })
            }
          }).catch(err => {
            this.$refs.configForm.resetFields()
            console.log(err)
            this.loading2 = false
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    // 处理环境配置弹窗
    handleConfigDialog() {
      if (this.multipleSelection.length === 1) {
        this.dialogConfigVisible = true
        this.loading3 = true
        configData(this.multipleSelection[0].ID).then(res => {
          this.loading3 = false
          if (res.Status) {
            if (res.Data.length > 0) {
              this.configId1 = res.Data[0].ID
              this.configForm.gitLoginName = res.Data[0].Values
              if (res.Data[0].Enable === 'true') {
                this.start1 = true
              } else {
                this.start1 = false
              }

              this.configId2 = res.Data[1].ID
              this.configForm.gitLoginPsw = res.Data[1].Values
              this.start2 = res.Data[1].Enable
              if (res.Data[1].Enable === 'true') {
                this.start2 = true
              } else {
                this.start2 = false
              }

              this.configId3 = res.Data[2].ID
              this.configForm.publishKey = res.Data[2].Values
              this.start3 = res.Data[2].Enable
              if (res.Data[2].Enable === 'true') {
                this.start3 = true
              } else {
                this.start3 = false
              }

              this.configId4 = res.Data[3].ID
              this.configForm.pullPath = res.Data[3].Values
              if (res.Data[3].Enable === 'true') {
                this.start4 = true
              } else {
                this.start4 = false
              }

              this.configId5 = res.Data[4].ID
              this.configForm.runPath = res.Data[4].Values
              if (res.Data[4].Enable === 'true') {
                this.start5 = true
              } else {
                this.start5 = false
              }

              this.configId6 = res.Data[5].ID
              this.configForm.publishBeforeCmd = res.Data[5].Values
              if (res.Data[5].Enable === 'true') {
                this.start6 = true
              } else {
                this.start6 = false
              }

              this.configId7 = res.Data[6].ID
              this.configForm.publishCmd = res.Data[6].Values
              if (res.Data[6].Enable === 'true') {
                this.start7 = true
              } else {
                this.start7 = false
              }

              this.configId8 = res.Data[7].ID
              this.configForm.publishAfterCmd = res.Data[7].Values
              if (res.Data[7].Enable === 'true') {
                this.start8 = true
              } else {
                this.start8 = false
              }
            } else {
              console.log(1)
              this.start1 = true
              this.start2 = true
              this.start3 = true
              this.start4 = true
              this.start5 = true
              this.start6 = true
              this.start7 = true
              this.start8 = true
              this.configId1 = 0
              this.configId2 = 0
              this.configId3 = 0
              this.configId4 = 0
              this.configId5 = 0
              this.configId6 = 0
              this.configId7 = 0
              this.configId8 = 0
            }
          }
        }).catch(err => {
          console.log(err)
          this.loading3 = false
        })
      } else {
        this.$message({
          message: '只能选择一条要配置的数据',
          type: 'warning'
        })
      }
    },
    // 初始化用户角色
    initProjectType() {
      projectType().then(res => {
        this.projectOptions = res.Data
      })
    },
    // 处理保存弹窗关闭
    handleClose(done) {
      this.$refs.ruleForm.resetFields()
      done()
    },
    // 处理环境配置弹窗关闭
    handleConfigClose(done) {
      this.$refs.configForm.resetFields()
      done()
    },
    // 查询
    handleSearch() {
      this.currentPage = 1
      this.initTableData()
    },
    // /初始化表格数据
    initTableData() {
      this.page = false
      this.loading = true
      TableData(this.form, this.currentPage, this.pageSize, parseInt(sessionStorage.getItem('userId'))).then(res => {
        this.tableData = res.Data
        this.total = res.PageCount
        this.loading = false
        if (this.total > 0) {
          this.page = true
        }
      }).catch(err => {
        this.page = false
        this.loading = false
      })
    },
    // 新增和修改项目
    handleSave() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          this.loading1 = true
          if (this.title === '新增项目') {
            addModifyProjecy(this.ruleForm, 0, parseInt(sessionStorage.getItem('userId'))).then(res => {
              this.loading1 = false
              if (res.Status) {
                this.$message({
                  message: '新增成功',
                  type: 'success'
                })
                this.dialogSaveVisible = false
                this.initTableData()
              } else {
                this.$message({
                  message: '新增失败，请联系系统管理员！',
                  type: 'warning'
                })
              }
            })
            this.$refs.ruleForm.resetFields()
          } else {
            addModifyProjecy(this.ruleForm, this.multipleSelection[0].ID, parseInt(sessionStorage.getItem('userId'))).then(res => {
              this.loading1 = false
              if (res.Status) {
                this.$message({
                  message: '修改成功',
                  type: 'success'
                })
                this.dialogSaveVisible = false
                this.initTableData()
              } else {
                this.$message({
                  message: '修改失败，请联系系统管理员！',
                  type: 'warning'
                })
              }
            })
            this.$refs.ruleForm.resetFields()
          }
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    // 删除
    handleDeleteData() {
      const id = []
      for (let i = 0; i < this.multipleSelection.length; i++) {
        id.push(this.multipleSelection[i].ID)
      }
      this.$confirm('您是否确认删除此行（多行）数据！', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteData(id).then(res => {
          if (res.Status) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.initTableData()
          } else {
            this.$message({
              message: res.Msg,
              type: 'warning'
            })
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    // 处理弹窗
    handleDialog(val) {
      if (val === 1) {
        this.dialogSaveVisible = true
        this.title = '新增项目'
        this.ruleForm.projectName = ''
        this.ruleForm.projectType = ''
        this.ruleForm.serverAddress = ''
        this.ruleForm.gitAddress = ''
        this.ruleForm.path = ''
        this.ruleForm.onlineAddress = ''
        this.ruleForm.gitName = ''
      } else {
        if (this.multipleSelection.length === 1) {
          this.dialogSaveVisible = true
          this.title = '编辑项目'
          this.ruleForm.projectName = this.multipleSelection[0].Name
          this.ruleForm.projectType = this.multipleSelection[0].ProjectTypeId
          this.ruleForm.serverAddress = this.multipleSelection[0].ServerAddress
          this.ruleForm.gitAddress = this.multipleSelection[0].GitAddress
          this.ruleForm.path = this.multipleSelection[0].ConfAddr
          this.ruleForm.onlineAddress = this.multipleSelection[0].OnlineAccessAddress
          this.ruleForm.gitName = this.multipleSelection[0].WarehouseName
        } else {
          this.$message({
            message: '只能选择一条要编辑的数据',
            type: 'warning'
          })
        }
      }
    },
    // 处理勾选
    handleSelectionChange(val) {
      this.multipleSelection = val
      console.log(this.multipleSelection)
      if (this.multipleSelection.length > 0) {
        this.disabled = false
      } else {
        this.disabled = true
      }
    },
    // 处理每页条数
    handleSizeChange(val) {
      this.pageSize = val
      this.initTableData()
    },
    // 处理分页当前页
    handleCurrentChange(val) {
      this.currentPage = val
      this.initTableData()
    }
  }
}
</script>
<style scoped>
  .search-text span{
    line-height: 36px;
    font-size: 14px;
  }
  .zoom .el-button{
    padding:4px 0;
    width: 38px;
    margin: 0;
    line-height: 16px;
    font-size: 12px;
  }
  .el-table .cell.el-tooltip{
    color:#409EFF;
  }
</style>

