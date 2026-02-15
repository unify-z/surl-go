<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo-section">
        <h1>{{ configStore.siteName }}</h1>
      </div>

      <n-tabs v-model:value="activeTab" size="large" animated>
        <n-tab-pane name="login" tab="登录">
          <n-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            size="large"
          >
            <n-form-item path="username" label="用户名">
              <n-input
                v-model:value="loginForm.username"
                placeholder="请输入用户名"
                :input-props="{ autocomplete: 'username' }"
              >
                <template #prefix>
                  <n-icon :component="PersonOutline" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="password" label="密码">
              <n-input
                v-model:value="loginForm.password"
                type="password"
                show-password-on="click"
                placeholder="请输入密码"
                :input-props="{ autocomplete: 'current-password' }"
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <n-icon :component="LockClosedOutline" />
                </template>
              </n-input>
            </n-form-item>
            <n-button
              type="primary"
              block
              size="large"
              :loading="loginLoading"
              @click="handleLogin"
            >
              登录
            </n-button>
          </n-form>
        </n-tab-pane>

        <n-tab-pane v-if="configStore.allowRegistration" name="register" tab="注册">
          <n-form
            ref="registerFormRef"
            :model="registerForm"
            :rules="registerRules"
            size="large"
          >
            <n-form-item path="username" label="用户名">
              <n-input
                v-model:value="registerForm.username"
                placeholder="请输入用户名"
                :input-props="{ autocomplete: 'username' }"
              >
                <template #prefix>
                  <n-icon :component="PersonOutline" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="email" label="邮箱">
              <n-input
                v-model:value="registerForm.email"
                placeholder="请输入邮箱"
                :input-props="{ autocomplete: 'email' }"
              >
                <template #prefix>
                  <n-icon :component="MailOutline" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="emailCode" label="验证码">
              <n-input-group>
                <n-input
                  v-model:value="registerForm.emailCode"
                  placeholder="请输入邮箱验证码"
                  style="width: 60%"
                />
                <n-button
                  :disabled="countdown > 0"
                  :loading="sendingCode"
                  @click="sendVerifyCode"
                  style="width: 40%"
                >
                  {{ countdown > 0 ? `${countdown}秒后重试` : "发送验证码" }}
                </n-button>
              </n-input-group>
            </n-form-item>
            <n-form-item path="password" label="密码">
              <n-input
                v-model:value="registerForm.password"
                type="password"
                show-password-on="click"
                placeholder="请输入密码（至少6位）"
                :input-props="{ autocomplete: 'new-password' }"
              >
                <template #prefix>
                  <n-icon :component="LockClosedOutline" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="confirmPassword" label="确认密码">
              <n-input
                v-model:value="registerForm.confirmPassword"
                type="password"
                show-password-on="click"
                placeholder="请再次输入密码"
                :input-props="{ autocomplete: 'new-password' }"
                @keyup.enter="handleRegister"
              >
                <template #prefix>
                  <n-icon :component="LockClosedOutline" />
                </template>
              </n-input>
            </n-form-item>
            <n-button
              type="primary"
              block
              size="large"
              :loading="registerLoading"
              @click="handleRegister"
            >
              注册
            </n-button>
          </n-form>
        </n-tab-pane>
      </n-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import {
  useMessage,
  NTabs,
  NTabPane,
  NForm,
  NFormItem,
  NInput,
  NInputGroup,
  NButton,
  NIcon,
} from "naive-ui";
import {
  PersonOutline,
  LockClosedOutline,
  MailOutline,
} from "@vicons/ionicons5";
import md5 from "js-md5";
import api from "../api";
import { useUserStore } from "../stores/user";
import { useConfigStore } from "../stores/config";

const router = useRouter();
const message = useMessage();
const userStore = useUserStore();
const configStore = useConfigStore();

const activeTab = ref("login");
const loginLoading = ref(false);
const registerLoading = ref(false);
const sendingCode = ref(false);
const countdown = ref(0);

const loginForm = reactive({
  username: "",
  password: "",
});

const registerForm = reactive({
  username: "",
  email: "",
  emailCode: "",
  password: "",
  confirmPassword: "",
});

const loginRules = {
  username: { required: true, message: "请输入用户名", trigger: "blur" },
  password: { required: true, message: "请输入密码", trigger: "blur" },
};

const registerRules = {
  username: {
    required: true,
    message: "请输入用户名",
    trigger: "blur",
    min: 3,
    max: 20,
  },
  email: {
    required: true,
    message: "请输入邮箱",
    trigger: "blur",
    type: "email",
  },
  emailCode: { required: true, message: "请输入验证码", trigger: "blur" },
  password: {
    required: true,
    message: "请输入密码",
    trigger: "blur",
    min: 6,
  },
  confirmPassword: [
    { required: true, message: "请确认密码", trigger: "blur" },
    {
      validator: (rule, value) => value === registerForm.password,
      message: "两次密码不一致",
      trigger: "blur",
    },
  ],
};

const loginFormRef = ref(null);
const registerFormRef = ref(null);

async function handleLogin() {
  try {
    await loginFormRef.value?.validate();
    loginLoading.value = true;

    const res = await api.post("/user/login", {
      username: loginForm.username,
      password_md5: md5(loginForm.password),
    });
    if (res.code === 200) {
      userStore.setUser(res.data);
      userStore.setToken(res.data.token);
      message.success("登录成功");
      router.push("/dashboard");
    } else {
      message.error(res.message || "登录失败");
    }
  } catch (error) {
    if (error?.response?.data?.message) {
      message.error(error.response.data.message);
    }
  } finally {
    loginLoading.value = false;
  }
}

async function sendVerifyCode() {
  if (!registerForm.email) {
    message.warning("请输入邮箱");
    return;
  }

  try {
    sendingCode.value = true;
    const res = await api.post("/user/create_email_code", {
      email: registerForm.email,
    });

    if (res.code === 200) {
      message.success("验证码已发送");
      countdown.value = 60;
      const timer = setInterval(() => {
        countdown.value--;
        if (countdown.value <= 0) {
          clearInterval(timer);
        }
      }, 1000);
    } else {
      message.error(res.message || "发送失败");
    }
  } catch (error) {
    if (error?.response?.data?.message) {
      message.error(error.response.data.message);
    }
  } finally {
    sendingCode.value = false;
  }
}

async function handleRegister() {
  try {
    await registerFormRef.value?.validate();
    registerLoading.value = true;

    const res = await api.post("/user/register", {
      username: registerForm.username,
      email: registerForm.email,
      password_md5: md5(registerForm.password),
      email_verify_code: registerForm.emailCode,
    });

    if (res.status === 200) {
      userStore.setUser(res.data);
      userStore.setToken(res.data.token);
      message.success("注册成功");
      router.push("/dashboard");
    } else {
      message.error(res.message || "注册失败");
    }
  } catch (error) {
    if (error?.response?.data?.message) {
      message.error(error.response.data.message);
    }
  } finally {
    registerLoading.value = false;
  }
}

onMounted(() => {
  configStore.loadConfig();
});
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  background: white;
  padding: 40px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.logo-section {
  text-align: center;
  margin-bottom: 32px;
}

.logo-section h1 {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.logo-section p {
  color: #666;
  margin: 0;
  font-size: 14px;
}

:deep(.n-tabs-nav) {
  margin-bottom: 24px;
}

:deep(.n-form-item) {
  margin-bottom: 20px;
}
</style>
