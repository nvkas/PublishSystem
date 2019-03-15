<template>
  <div class="navbar">
    <hamburger :toggle-click="toggleSideBar" :is-active="sidebar.opened" class="hamburger-container"/>

    <breadcrumb class="breadcrumb-container"/>

    <div class="right-menu">
      <template v-if="device!=='mobile'">
        <lang-select v-show="false" class="international right-menu-item"/>
      </template>
      <el-dropdown class="avatar-container right-menu-item hover-effect" trigger="click">
        <div class="avatar-wrapper">
          <img src="../../../images/tou.gif" class="user-avatar">
          <span class="top-role">{{ username }}</span>
          <i class="el-icon-caret-bottom xia"/>
        </div>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item>
            <span style="display:block;" @click="dialogEditVisible = true">个人信息</span>
          </el-dropdown-item>
          <el-dropdown-item divided>
            <span style="display:block;" @click="logout">{{ $t('navbar.logOut') }}</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
    <el-dialog :visible.sync="dialogEditVisible" title="个人设置" >
      <p class="person-title text-center mb16">个人资料中心</p>
      <el-tabs v-model="activeName2" type="card">
        <el-tab-pane label="基本信息" name="first">
          <el-form ref="ruleForm" :model="ruleForm" :rules="rules" label-width="100px" class="demo-ruleForm" label-position="right">
            <el-row>
              <el-col :span="12">
                <el-form-item label="登录账号：">
                  {{ loginName }}
                </el-form-item>
                <el-form-item label="真实姓名：" prop="name">
                  <el-input v-model="ruleForm.name"/>
                </el-form-item>
                <el-form-item label="性别：" prop="sex">
                  <el-radio v-model="ruleForm.sex" label="1">男</el-radio>
                  <el-radio v-model="ruleForm.sex" label="2">女</el-radio>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="手机号码：" prop="phone">
                  <el-input v-model="ruleForm.phone"/>
                </el-form-item>
                <el-form-item label="邮箱：" prop="email">
                  <el-input v-model="ruleForm.email"/>
                </el-form-item>
              </el-col>
            </el-row>
            <div class="text-center">
              <el-button type="primary" size="small">保存</el-button>
            </div>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="修改密码" name="second">
          <el-form ref="ruleForm1" :model="ruleForm1" :rules="rules1" label-width="100px" class="demo-ruleForm" label-position="right">
            <el-form-item prop="oldPsw" label="旧密码" class="relative">
              <el-input
                :type="passwordType"
                v-model="ruleForm1.oldPsw"
                placeholder="请输入旧密码"
                name="password"/>
              <span class="show-pwd" @click="showPwd"><svg-icon icon-class="eye" /></span>
            </el-form-item>
            <el-form-item prop="newPsw" label="新密码">
              <el-input
                :type="passwordType1"
                v-model="ruleForm1.newPsw"
                placeholder="请输入新密码"
                name="password"/>
              <span class="show-pwd" @click="showPwd1"><svg-icon icon-class="eye" /></span>
            </el-form-item>
            <el-form-item prop="doPsw" label="确认密码">
              <el-input
                :type="passwordType2"
                v-model="ruleForm1.doPsw"
                placeholder="请确认密码"
                name="password"/>
              <span class="show-pwd" @click="showPwd2"><svg-icon icon-class="eye" /></span>
            </el-form-item>
            <div class="text-center">
              <el-button type="primary" size="small">保存</el-button>
            </div>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import ErrorLog from '@/components/ErrorLog'
import Screenfull from '@/components/Screenfull'
import SizeSelect from '@/components/SizeSelect'
import LangSelect from '@/components/LangSelect'
import ThemePicker from '@/components/ThemePicker'
import Search from '@/components/HeaderSearch'
import user from './../../../store/modules/user'
import { person, modify, modifyPsw } from '@/api/personMsg'
import { validatePhone, validateEmail } from '@/utils/validate'

export default {
  components: {
    Breadcrumb,
    Hamburger,
    ErrorLog,
    Screenfull,
    SizeSelect,
    LangSelect,
    ThemePicker,
    Search
  },
  data() {
    const phone = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入您的手机号码！'))
      } else if (!validatePhone(value)) {
        callback(new Error('请正确填写你的手机号码!'))
      } else {
        callback()
      }
    }
    const email = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入您的信箱！'))
      } else if (!validateEmail(value)) {
        callback(new Error('请正确输入您的信箱！格式：xx@xx.com'))
      } else {
        callback()
      }
    }
    const do_password = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请确认密码'))
      } else if (value !== this.ruleForm1.newPsw) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    return {
      loginName: '',
      passwordType: 'password',
      passwordType1: 'password',
      passwordType2: 'password',
      dialogEditVisible: false,
      username: user.state.name,
      activeName2: 'first',
      ruleForm: {
        name: '',
        birth: '',
        company: '',
        address: '',
        sex: '',
        phone: '',
        email: ''
      },
      rules: {
        name: [
          { required: true, message: '请输入您的真实姓名', trigger: 'blur' }
        ],
        sex: [
          { required: true, message: '请选择您的性别', trigger: 'blur' }
        ],
        phone: [
          { required: true, message: '请输入您的手机号码', trigger: 'blur' },
          { validator: phone, trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入您的邮箱', trigger: 'blur' },
          { validator: email, trigger: 'blur' }
        ]
      },
      ruleForm1: {
        oldPsw: '',
        newPsw: '',
        doPsw: ''
      },
      rules1: {
        oldPsw: [
          { required: true, message: '请输入旧密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度为6~20个字符', trigger: 'blur' }
        ],
        newPsw: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度为6~20个字符', trigger: 'blur' }
        ],
        doPsw: [
          { required: true, message: '请确认密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度为6~20个字符', trigger: 'blur' },
          { validator: do_password, trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    window.addEventListener('beforeunload', () => {
      sessionStorage.setItem('name', this.username)
    })
    if (sessionStorage.getItem('name')) {
      this.username = sessionStorage.getItem('name')
    }
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'name',
      'avatar',
      'device'
    ])
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('toggleSideBar')
    },
    // 小眼睛
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
    },
    showPwd1() {
      if (this.passwordType1 === 'password') {
        this.passwordType1 = ''
      } else {
        this.passwordType1 = 'password'
      }
    },
    showPwd2() {
      if (this.passwordType2 === 'password') {
        this.passwordType2 = ''
      } else {
        this.passwordType2 = 'password'
      }
    },
    // 退出登录
    logout() {
      this.$store.dispatch('FedLogOut').then(() => {
        this.$router.push({ path: '/login' })
      })
    }
  }
}
</script>
<style>
  .navbar .right-menu .avatar-container .avatar-wrapper .el-icon-caret-bottom[data-v-797e31be].xia{
    top: 14px;
  }
  .el-dialog__header {
    border-bottom: 1px solid #ddd;
  }
</style>
<style rel="stylesheet/scss" lang="scss" scoped>
  $dark_gray:#889aa4;
.navbar {
  height: 50px;
  overflow: hidden;
  .person-title {
    line-height: 20px;
    font-size: 16px;
    margin: 0;
  }
  .show-pwd {
    position: absolute;
    right: 10px;
    top: 1px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
  .top-role{
    line-height: 40px;
    height: 40px;
    vertical-align: top;
    font-size: 16px;
  }
  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;

    &:hover {
      background: rgba(0, 0, 0, .025)
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .avatar-container {
      margin-right: 30px;

      .avatar-wrapper {
        margin-top: 5px;
        position: relative;

        .user-avatar {
          cursor: pointer;
          width: 40px;
          height: 40px;
          border-radius: 10px;
        }

        .el-icon-caret-bottom {
          cursor: pointer;
          position: absolute;
          right: -20px;
          top: 25px;
          font-size: 12px;
        }
      }
    }
  }
}
</style>
