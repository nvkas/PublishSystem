<template>
  <div class="app-container">
    <div class="mb16 search-text">
      <span>角色名称：</span>
      <el-input
        v-model="name"
        clearable
        style="width: 240px;"/>
      <el-button type="primary" size="small" @click="handleSearch">查询</el-button>
    </div>
    <div class="mb16">
      <el-button type="primary" size="small" @click="handleDialog(1)">新增</el-button>
      <el-button :disabled="disabled" type="primary" size="small" @click="handleDialog(2)">编辑</el-button>
      <el-button :disabled="disabled" type="danger" size="small">删除</el-button>
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
          prop="Name"
          label="角色名"/>
        <el-table-column
          prop="Remarks"
          label="备注"/>
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
      <el-form ref="ruleForm" :model="ruleForm" :rules="rules" label-width="100px" class="demo-ruleForm">
        <el-form-item label="角色名" prop="role_name">
          <el-input v-model.trim="ruleForm.role_name"/>
        </el-form-item>
        <el-form-item label="备注" prop="msg">
          <el-input v-model.trim="ruleForm.msg"/>
        </el-form-item>
        <el-form-item label="菜单权限" prop="menu">
          <el-tree
            v-loading="loading2"
            ref="tree"
            :data="data2"
            :default-expand-all="false"
            :default-checked-keys="menuId"
            :props="defaultProps"
            show-checkbox
            node-key="id"
            highlight-current/>
        </el-form-item>
        <div class="text-center">
          <el-button :loading="loading1" type="primary" size="small" @click="handleSave">保存</el-button>
        </div>
      </el-form>
    </el-dialog>
  </div>
</template>
<script>
import { TableData, tree, addModifyRole, checkRoleName, roleMenu } from '@/api/roleManage'
export default {
  name: 'RoleManage',
  data() {
    const name = (rule, value, callback) => {
      if (this.title === '新增角色') {
        checkRoleName(this.ruleForm.role_name).then(res => {
          if (res.Status) {
            callback()
          } else {
            callback(new Error('用户名已存在'))
          }
        })
      } else {
        callback()
      }
    }
    const menu = (rule, value, callback) => {
      if (this.$refs.tree.getCheckedNodes().length > 0) {
        callback()
      } else {
        callback(new Error('请勾选菜单项'))
      }
    }
    return {
      userOptions: [],
      dialogSaveVisible: false,
      loading1: false,
      loading2: false,
      disabled: true,
      menuId: [],
      page: false,
      menu: [],
      defaultProps: {
        children: 'children',
        label: 'label'
      },
      loading: false,
      name: '',
      title: '',
      total: '',
      pageSize: 10,
      currentPage: 1,
      tableData: [],
      multipleSelection: [],
      ruleForm: {
        role_name: '',
        msg: ''
      },
      rules: {
        role_name: [
          { required: true, message: '请输入角色名', trigger: 'blur' },
          { validator: name, trigger: 'blur' }
        ],
        menu: [
          { required: true, validator: menu, trigger: 'blur' }
        ]
      },
      data2: []
    }
  },
  created() {
    this.initRoleTableData()
    this.initTree()
  },
  methods: {
    // 处理保存
    handleSave() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const id = this.$refs.tree.getCheckedNodes()
          const node = []
          for (let i = 0; i < id.length; i++) {
            node.push(id[i].id)
          }
          this.loading1 = true
          if (this.title === '新增角色') {
            addModifyRole(this.ruleForm, node, 0).then(res => {
              this.loading1 = false
              if (res.Status) {
                this.$message({
                  message: '新增成功',
                  type: 'success'
                })
                this.dialogSaveVisible = false
                this.initRoleTableData()
              } else {
                this.$message({
                  message: '新增失败，请联系系统管理员！',
                  type: 'warning'
                })
              }
            })
            this.$refs.ruleForm.resetFields()
          } else {
            addModifyRole(this.ruleForm, node, this.multipleSelection[0].ID).then(res => {
              this.loading1 = false
              if (res.Status) {
                this.$message({
                  message: '修改成功',
                  type: 'success'
                })
                this.dialogSaveVisible = false
                this.initRoleTableData()
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
    // 处理弹窗关闭
    handleClose(done) {
      this.$refs.ruleForm.resetFields()
      done()
    },
    // 查询
    handleSearch() {
      this.currentPage = 1
      this.initRoleTableData()
    },
    // 初始化表格数据
    initRoleTableData() {
      this.page = false
      this.loading = true
      TableData(this.name, this.pageSize, this.currentPage).then(res => {
        this.tableData = res.Data
        this.total = res.PageCount
        this.loading = false
        if (this.total > 0) {
          this.page = true
        }
      }).catch(err => {
        this.loading = false
        this.page = false
      })
    },
    // 初始化弹窗菜单权限tree
    initTree() {
      tree().then(res => {
        const menu = res.Data
        for (let i = 0; i < menu.length; i++) {
          const child = []
          this.data2.push({
            id: menu[i].ID,
            label: menu[i].MenuName,
            children: child
          })
          if (menu[i].Menu) {
            for (let j = 0; j < menu[i].Menu.length; j++) {
              const child1 = []
              child.push({
                id: menu[i].Menu[j].ID,
                label: menu[i].Menu[j].MenuName,
                children: child1
              })
              if (menu[i].Menu[j].Menu) {
                for (let k = 0; k < menu[i].Menu[j].Menu.length; k++) {
                  child1.push({
                    id: menu[i].Menu[j].Menu[k].ID,
                    label: menu[i].Menu[j].Menu[k].MenuName
                  })
                }
              }
            }
          }
        }
      }).catch(err => {

      })
    },
    // 处理弹窗
    handleDialog(val) {
      if (val === 1) {
        this.title = '新增角色'
        this.dialogSaveVisible = true
        this.ruleForm.role_name = ''
        this.ruleForm.msg = ''
        this.$refs.tree.setCheckedNodes([])
      } else {
        if (this.multipleSelection.length === 1) {
          this.title = '修改角色'
          this.dialogSaveVisible = true
          this.ruleForm.role_name = this.multipleSelection[0].Name
          this.ruleForm.msg = this.multipleSelection[0].Remarks
          this.loading2 = true
          roleMenu(this.multipleSelection[0].ID).then(res => {
            this.loading2 = false
            if (res.Data) {
              this.$refs.tree.setCheckedKeys(res.Data)
            } else {
              this.$refs.tree.setCheckedNodes([])
            }
          })
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
      this.initRoleTableData()
    },
    // 处理分页当前页
    handleCurrentChange(val) {
      this.currentPage = val
      this.initRoleTableData()
    }
  }
}
</script>
<style scoped>
  .app-container{
    padding: 20px 15px;
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
  .search-text span{
    line-height: 36px;
    font-size: 14px;
  }
</style>

