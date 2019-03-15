<template>
  <div class="app-container">
    <el-form ref="form" :model="form" label-width="80px">
      <el-row class="search-text">
        <el-col :span="5">
          <el-form-item label="登录账号:" prop="loginName">
            <el-input
              v-model="form.loginName"
              placeholder="登录账号"
              clearable
              class="wid100"/>
          </el-form-item>
        </el-col>
        <el-col :span="5">
          <el-form-item label="真实姓名:" prop="realName">
            <el-input
              v-model="form.realName"
              placeholder="真实姓名"
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
          prop="LoginName"
          label="登录账号"
          show-overflow-tooltip/>
        <el-table-column
          prop="Name"
          label="真实姓名"
          show-overflow-tooltip/>
        <el-table-column
          prop="roleName"
          label="用户角色"
          show-overflow-tooltip/>
        <el-table-column
          prop="Phone"
          label="手机号"
          width="120"/>
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
            <el-form-item label="所属部门" prop="department">
              <el-select v-model.trim="ruleForm.department" clearable placeholder="请选择" class="wid100">
                <el-option
                  v-for="item in departmentOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"/>
              </el-select>
            </el-form-item>
            <el-form-item label="用户角色" prop="role">
              <el-select v-model="ruleForm.role" clearable multiple placeholder="请选择" class="wid100">
                <el-option
                  v-for="item in roleOptions"
                  :key="item.ID"
                  :label="item.Name"
                  :value="item.ID"/>
              </el-select>
            </el-form-item>
            <el-form-item label="E-mail" prop="email">
              <el-input v-model.trim="ruleForm.email"/>
            </el-form-item>
            <el-form-item label="手机号码" prop="phone">
              <el-input v-model.trim="ruleForm.phone"/>
            </el-form-item>
            <el-form-item label="性别" prop="sex">
              <el-radio v-model="ruleForm.sex" label="1">男</el-radio>
              <el-radio v-model="ruleForm.sex" label="2">女</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户类型" prop="userType">
              <el-select v-model.trim="ruleForm.userType" clearable placeholder="请选择" class="wid100">
                <el-option
                  v-for="item in userTypeOption"
                  :key="item.ID"
                  :label="item.Name"
                  :value="item.ID"/>
              </el-select>
            </el-form-item>
            <el-form-item label="登录账号" prop="loginName">
              <el-input v-model.trim="ruleForm.loginName" :disabled="loginNameDisabled"/>
            </el-form-item>
            <el-form-item label="真实姓名" prop="realName">
              <el-input v-model.trim="ruleForm.realName"/>
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model.trim="ruleForm.password" type="password"/>
            </el-form-item>
            <el-form-item label="项目配置路径" prop="path">
              <el-input v-model.trim="ruleForm.path"/>
            </el-form-item>
          </el-col>
        </el-row>
        <div class="text-center">
          <el-button :loading="loading1" type="primary" size="small" @click="handleSave">保存</el-button>
        </div>
      </el-form>
    </el-dialog>
  </div>
</template>
<script>
import { validatePhone, validateEmail, validateName, validateUserName } from '@/utils/validate'
import { deleteData, addModifyUser, userTableData, checkLoginName, role, userType } from '@/api/userManage'
export default {
  name: 'UserManage',
  data() {
    const phone = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入您的手机号码'))
      } else if (!validatePhone(value)) {
        callback(new Error('请正确填写你的手机号码'))
      } else {
        callback()
      }
    }
    const LoginName = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入用户名'))
      } else if (!validateUserName(value)) {
        callback(new Error('首字母必须是英文，后面是英文和数字，长度6-20个字符'))
      } else if (this.title === '新增用户') {
        checkLoginName(value).then((res) => {
          if (res.Status) {
            callback()
          } else {
            callback(new Error('用户名已存在！'))
          }
        }).catch(() => {
          return false
        })
      } else {
        callback()
      }
    }
    const Name = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入本人真实姓名'))
      } else if (!validateName(value)) {
        callback(new Error('真实姓名必须是汉字'))
      } else {
        callback()
      }
    }
    const Email = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入邮箱号码'))
      } else if (!validateEmail(value)) {
        callback(new Error('邮箱格式错误'))
      } else {
        callback()
      }
    }
    return {
      disabled: true,
      loginNameDisabled: false,
      form: {
        type: '',
        unit: '',
        loginName: '',
        realName: ''
      },
      loading1: false,
      page: false,
      loading: false,
      title: '',
      dialogSaveVisible: false,
      total: '',
      currentPage: 1,
      pageSize: 10,
      departmentOptions: [
        {
          value: '前端',
          label: '前端'
        }, {
          value: '后台',
          label: '后台'
        }
      ],
      roleOptions: [],
      tableData: [],
      multipleSelection: [],
      userTypeOption: [],
      ruleForm: {
        phone: '',
        password: '',
        department: '',
        email: '',
        loginName: '',
        realName: '',
        doPassword: '',
        sex: '',
        role: [],
        userType: '',
        path: ''
      },
      rules: {
        realName: [
          { required: true, message: '请输入您的真实姓名', trigger: 'blur' },
          { validator: Name, trigger: 'blur' }
        ],
        loginName: [
          { required: true, message: '请输入您的登录账号', trigger: 'blur' },
          { validator: LoginName, trigger: 'blur' }
        ],
        phone: [
          { required: true, message: '请输入您的手机号码', trigger: 'blur' },
          { validator: phone, trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入您的密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度为6~20个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入您的email', trigger: 'blur' },
          { validator: Email, trigger: 'blur' }
        ],
        department: [
          { required: true, message: '请选择您的部门名称', trigger: 'blur' }
        ],
        role: [
          { required: true, message: '请选择您的用户角色', trigger: 'blur' }
        ],
        sex: [
          { required: true, message: '请选择性别', trigger: 'blur' }
        ],
        path: [
          { required: true, message: '请输入您的项目配置路径', trigger: 'blur' }
        ],
        userType: [
          { required: true, message: '请选择您的用户类型', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.initUserTableData()
    this.initRole()
    this.initUserType()
  },
  methods: {
    // 初始化用户类型
    initUserType() {
      userType().then(res => {
        this.userTypeOption = res.Data
      })
    },
    // 初始化用户角色
    initRole() {
      role().then(res => {
        this.roleOptions = res.Data
      })
    },
    // 处理弹窗关闭
    handleClose(done) {
      this.$refs.ruleForm.resetFields()
      done()
    },
    // 查询
    handleSearch() {
      this.currentPage = 1
      this.initUserTableData()
    },
    // /初始化表格数据
    initUserTableData() {
      this.page = false
      this.loading = true
      userTableData(this.form.realName, this.form.loginName, this.currentPage, this.pageSize).then(res => {
        this.tableData = res.Data
        for (let i = 0, len = this.tableData.length; i < len; i++) {
          const roles = []
          if (this.tableData[i].Role) {
            for (let j = 0, len = this.tableData[i].Role.length; j < len; j++) {
              roles.push(this.tableData[i].Role[j].Name)
            }
            this.tableData[i].roleName = roles.join(',')
          }
        }
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
    // 新增和修改用户
    handleSave() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          this.loading1 = true
          if (this.title === '新增用户') {
            console.log(1)
            addModifyUser(this.ruleForm, 0).then(res => {
              this.loading1 = false
              if (res.Status) {
                this.$message({
                  message: '新增成功',
                  type: 'success'
                })
                this.dialogSaveVisible = false
                this.initUserTableData()
              } else {
                this.$message({
                  message: '新增失败，请联系系统管理员！',
                  type: 'warning'
                })
              }
            })
            this.$refs.ruleForm.resetFields()
          } else {
            addModifyUser(this.ruleForm, this.multipleSelection[0].ID).then(res => {
              this.loading1 = false
              if (res.Status) {
                this.$message({
                  message: '修改成功',
                  type: 'success'
                })
                this.dialogSaveVisible = false
                this.initUserTableData()
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
          if (res.data.Status) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.initUserTableData()
          } else {
            this.$message({
              message: res.data.Msg,
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
        this.title = '新增用户'
        this.ruleForm.phone = ''
        this.ruleForm.password = ''
        this.ruleForm.department = ''
        this.ruleForm.email = ''
        this.ruleForm.loginName = ''
        this.ruleForm.realName = ''
        this.ruleForm.doPassword = ''
        this.ruleForm.sex = ''
        this.ruleForm.role = []
        this.ruleForm.path = ''
        this.ruleForm.userType = ''
      } else {
        if (this.multipleSelection.length === 1) {
          this.dialogSaveVisible = true
          this.title = '编辑用户'
          this.ruleForm.phone = this.multipleSelection[0].Phone
          this.ruleForm.password = this.multipleSelection[0].Password
          this.ruleForm.department = this.multipleSelection[0].DepartmentName
          this.ruleForm.email = this.multipleSelection[0].Email
          this.ruleForm.loginName = this.multipleSelection[0].LoginName
          this.ruleForm.realName = this.multipleSelection[0].Name
          this.ruleForm.doPassword = this.multipleSelection[0].Password
          this.ruleForm.sex = this.multipleSelection[0].Sex + ''
          this.ruleForm.path = this.multipleSelection[0].ProjectWorkPath
          this.ruleForm.userType = this.multipleSelection[0].UserTypeId
          this.ruleForm.role = []
          for (let i = 0, len = this.multipleSelection[0].Role.length; i < len; i++) {
            this.ruleForm.role.push(this.multipleSelection[0].Role[i].ID)
          }
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
      if (this.multipleSelection.length > 0) {
        this.disabled = false
      } else {
        this.disabled = true
      }
    },
    // 处理每页条数
    handleSizeChange(val) {
      this.pageSize = val
      this.initUserTableData()
    },
    // 处理分页当前页
    handleCurrentChange(val) {
      this.currentPage = val
      this.initUserTableData()
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

