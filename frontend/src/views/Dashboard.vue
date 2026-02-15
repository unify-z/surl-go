<template>
  <div class="dashboard">
    <div class="header">
      <div class="header-content">
        <div class="logo">
          <n-icon :component="LinkOutline" size="24" />
          <span class="logo-text">{{ configStore.siteName }}</span>
        </div>
        <div class="user-section">
          <n-dropdown v-if="userStore.isLoggedIn" :options="userOptions" @select="handleUserAction">
            <n-button text>
              <template #icon>
                <n-icon :component="PersonCircleOutline" size="24" />
              </template>
              {{ userStore.username }}
            </n-button>
          </n-dropdown>
          <n-button v-else type="primary" quaternary @click="router.push('/login')">
            登录 / 注册
          </n-button>
        </div>
      </div>
    </div>

    <div class="main-content">
      <div class="container">
        <div class="welcome-section" v-if="userStore.isLoggedIn">
          <h2>欢迎，{{ userStore.username }}</h2>
          <p>创建和管理你的短链接</p>
        </div>

        <div class="create-section" v-if="userStore.isLoggedIn || configStore.allowGuestToCreateUrl">
          <n-card title="创建短链接">
            <n-form
              ref="createFormRef"
              :model="createForm"
              :rules="createRules"
            >
              <n-form-item path="url" label="原始链接">
                <n-input-group>
                  <n-input
                    v-model:value="createForm.url"
                    placeholder="请输入要缩短的链接"
                    size="large"
                    @keyup.enter="handleCreate"
                  >
                    <template #prefix>
                      <n-icon :component="LinkOutline" />
                    </template>
                  </n-input>
                  <n-button
                    type="primary"
                    size="large"
                    :loading="creating"
                    @click="handleCreate"
                  >
                    生成短链接
                  </n-button>
                </n-input-group>
              </n-form-item>
            </n-form>
          </n-card>
        </div>

        <div class="links-section" v-if="userStore.isLoggedIn">
          <n-card title="我的短链接">
            <template #header-extra>
              <n-button text @click="loadLinks">
                <template #icon>
                  <n-icon :component="RefreshOutline" />
                </template>
              </n-button>
            </template>

            <n-spin :show="loading">
              <n-empty v-if="links.length === 0" description="还没有短链接" />
              <n-list v-else hoverable clickable>
                <n-list-item v-for="link in links" :key="link.ID">
                  <template #prefix>
                    <n-icon
                      :component="LinkOutline"
                      size="20"
                      color="#18a058"
                    />
                  </template>
                  <n-thing>
                    <template #header>
                      <n-space align="center">
                        <n-text strong>{{
                          getShortUrl(link.ShortCode)
                        }}</n-text>
                        <n-button
                          text
                          size="small"
                          @click="copyLink(getShortUrl(link.ShortCode))"
                        >
                          <template #icon>
                            <n-icon :component="CopyOutline" />
                          </template>
                        </n-button>
                      </n-space>
                    </template>
                    <template #description>
                      <n-ellipsis style="max-width: 600px">
                        {{ link.OriginalURL }}
                      </n-ellipsis>
                    </template>
                    <template #footer>
                      <n-text depth="3" style="font-size: 12px">
                        创建于 {{ formatDate(link.CreatedAt) }}
                      </n-text>
                    </template>
                  </n-thing>
                  <template #suffix>
                    <n-space>
                      <n-button size="small" @click="handleEdit(link)">
                        <template #icon>
                          <n-icon :component="CreateOutline" />
                        </template>
                        编辑
                      </n-button>
                      <n-button
                        size="small"
                        type="error"
                        @click="handleDelete(link)"
                      >
                        <template #icon>
                          <n-icon :component="TrashOutline" />
                        </template>
                        删除
                      </n-button>
                    </n-space>
                  </template>
                </n-list-item>
              </n-list>
            </n-spin>
          </n-card>
        </div>
      </div>
    </div>

    <n-modal v-model:show="showEditModal" preset="dialog" title="编辑短链接">
      <n-form :model="editForm">
        <n-form-item label="短链接代码">
          <n-input :value="editForm.short_code" disabled />
        </n-form-item>
        <n-form-item label="原始链接">
          <n-input
            v-model:value="editForm.original_url"
            placeholder="输入新的链接"
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-space>
          <n-button @click="showEditModal = false">取消</n-button>
          <n-button type="primary" :loading="updating" @click="confirmEdit"
            >确定</n-button
          >
        </n-space>
      </template>
    </n-modal>

    <n-modal v-model:show="showResultModal" preset="dialog" title="生成成功">
      <n-thing>
        <n-input-group style="margin-top: 10px">
          <n-input :value="resultUrl" readonly />
          <n-button type="primary" @click="copyLink(resultUrl)">复制</n-button>
        </n-input-group>
      </n-thing>
      <template #action>
        <n-button @click="showResultModal = false">关闭</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h, computed } from "vue";
import { useRouter } from "vue-router";
import {
  useMessage,
  useDialog,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NInputGroup,
  NButton,
  NIcon,
  NSpin,
  NEmpty,
  NList,
  NListItem,
  NThing,
  NSpace,
  NText,
  NEllipsis,
  NDropdown,
  NModal,
} from "naive-ui";
import {
  PersonCircleOutline,
  LinkOutline,
  CopyOutline,
  CreateOutline,
  TrashOutline,
  RefreshOutline,
  LogOutOutline,
  ShieldCheckmarkOutline,
} from "@vicons/ionicons5";
import api from "../api";
import { useUserStore } from "../stores/user";
import { useConfigStore } from "../stores/config";

const router = useRouter();
const message = useMessage();
const dialog = useDialog();
const userStore = useUserStore();
const configStore = useConfigStore();

const loading = ref(false);
const creating = ref(false);
const updating = ref(false);
const links = ref([]);
const showEditModal = ref(false);
const showResultModal = ref(false);
const resultUrl = ref("");

const createForm = reactive({
  url: "",
});

const editForm = reactive({
  short_code: "",
  original_url: "",
});

const createRules = {
  url: [
    { required: true, message: "请输入链接", trigger: "blur" },
    { type: "url", message: "请输入有效的链接", trigger: "blur" },
  ],
};

const createFormRef = ref(null);

const userOptions = computed(() => {
  const options = [];
  if (userStore.isAdmin) {
    options.push({
      label: "管理面板",
      key: "admin",
      icon: () => h(ShieldCheckmarkOutline),
    });
  }
  options.push({
    label: "退出登录",
    key: "logout",
    icon: () => h(LogOutOutline),
  });
  return options;
});

function handleUserAction(key) {
  if (key === "logout") {
    userStore.logout();
    router.push("/login");
  } else if (key === "admin") {
    router.push("/admin");
  }
}

function getShortUrl(code) {
  return `${window.location.origin}/s/${code}`;
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleString("zh-CN");
}

async function copyLink(url) {
  try {
    await navigator.clipboard.writeText(url);
    message.success("已复制到剪贴板");
  } catch {
    message.error("复制失败");
  }
}

async function loadLinks() {
  if (!userStore.isLoggedIn) return;
  try {
    loading.value = true;
    const res = await api.get("/surl/list");
    if (res.code === 200) {
      links.value = res.data || [];
    }
  } catch (error) {
    message.error("加载列表失败");
  } finally {
    loading.value = false;
  }
}

async function handleCreate() {
  try {
    await createFormRef.value?.validate();
    creating.value = true;

    const res = await api.post("/surl/create", {
      original_url: createForm.url,
    });

    if (res.code === 200) {
      const generatedUrl = getShortUrl(res.data.short_code);
      createForm.url = "";

      if (userStore.isLoggedIn) {
        message.success("创建成功");
        loadLinks();
      } else {
        resultUrl.value = generatedUrl;
        showResultModal.value = true;
      }
    } else {
      message.error(res.message || "创建失败");
    }
  } catch (error) {
    if (error?.response?.data?.message) {
      message.error(error.response.data.message);
    }
  } finally {
    creating.value = false;
  }
}

function handleEdit(link) {
  editForm.short_code = link.ShortCode;
  editForm.original_url = link.OriginalURL;
  showEditModal.value = true;
}

async function confirmEdit() {
  try {
    updating.value = true;
    const res = await api.post("/surl/update", {
      short_code: editForm.short_code,
      original_url: editForm.original_url,
    });

    if (res.code === 200) {
      message.success("更新成功");
      showEditModal.value = false;
      loadLinks();
    } else {
      message.error(res.message || "更新失败");
    }
  } catch (error) {
    if (error?.response?.data?.message) {
      message.error(error.response.data.message);
    }
  } finally {
    updating.value = false;
  }
}

function handleDelete(link) {
  dialog.warning({
    title: "确认删除",
    content: "确定要删除这个短链接吗？",
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        const res = await api.post("/surl/delete", {
          short_code: link.ShortCode,
        });

        if (res.code === 200) {
          message.success("删除成功");
          loadLinks();
        } else {
          message.error(res.message || "删除失败");
        }
      } catch (error) {
        if (error?.response?.data?.message) {
          message.error(error.response.data.message);
        }
      }
    },
  });
}

onMounted(() => {
  if (userStore.isLoggedIn) {
    loadLinks();
  }
});
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background: #f0f2f5;
}

.header {
  background: white;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.main-content {
  padding: 24px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-section {
  background: white;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.welcome-section h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  color: #333;
}

.welcome-section p {
  margin: 0;
  color: #666;
}

.create-section {
  margin-bottom: 24px;
}

:deep(.n-card) {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

:deep(.n-card__content) {
  padding: 24px;
}

:deep(.n-list-item) {
  padding: 16px;
  margin-bottom: 8px;
  transition: all 0.3s;
}

:deep(.n-list-item:hover) {
  background: #fafafa;
}
</style>