<template>
  <div class="admin">
    <div class="header">
      <div class="header-content">
        <div class="logo">
          <n-icon :component="ShieldCheckmarkOutline" size="24" />
          <span class="logo-text">管理面板</span>
        </div>
        <div class="user-section">
          <n-button text @click="router.push('/dashboard')">
            <template #icon>
              <n-icon :component="ArrowBackOutline" />
            </template>
            返回主页
          </n-button>
          <n-dropdown :options="userOptions" @select="handleUserAction">
            <n-button text>
              <template #icon>
                <n-icon :component="PersonCircleOutline" size="24" />
              </template>
              {{ userStore.username }}
            </n-button>
          </n-dropdown>
        </div>
      </div>
    </div>

    <div class="main-content">
      <div class="container">
        <n-tabs v-model:value="activeTab" type="line" size="large" animated>
          <n-tab-pane name="users" tab="用户管理">
            <n-card>
              <template #header>
                <n-space justify="space-between">
                  <span>用户列表</span>
                  <n-button text @click="loadUsers">
                    <template #icon>
                      <n-icon :component="RefreshOutline" />
                    </template>
                  </n-button>
                </n-space>
              </template>

              <n-spin :show="usersLoading">
                <n-data-table
                  :columns="userColumns"
                  :data="users"
                  :pagination="false"
                />
              </n-spin>
            </n-card>
          </n-tab-pane>

          <n-tab-pane name="links" tab="短链接管理">
            <n-card>
              <template #header>
                <n-space justify="space-between">
                  <span>短链接列表</span>
                  <n-space>
                    <n-pagination
                      v-model:page="currentPage"
                      :page-count="totalPages"
                      @update:page="loadLinks"
                    />
                    <n-button text @click="loadLinks">
                      <template #icon>
                        <n-icon :component="RefreshOutline" />
                      </template>
                    </n-button>
                  </n-space>
                </n-space>
              </template>

              <n-spin :show="linksLoading">
                <n-data-table
                  :columns="linkColumns"
                  :data="links"
                  :pagination="false"
                />
              </n-spin>
            </n-card>
          </n-tab-pane>
        </n-tabs>
      </div>
    </div>

    <n-modal
      v-model:show="showEditLinkModal"
      preset="dialog"
      title="编辑短链接"
    >
      <n-form :model="editLinkForm">
        <n-form-item label="短链接代码">
          <n-input :value="editLinkForm.short_code" disabled />
        </n-form-item>
        <n-form-item label="原始链接">
          <n-input v-model:value="editLinkForm.original_url" />
        </n-form-item>
        <n-form-item label="用户ID">
          <n-input v-model:value="editLinkForm.user_id" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-space>
          <n-button @click="showEditLinkModal = false">取消</n-button>
          <n-button type="primary" :loading="updating" @click="confirmEditLink"
            >确定</n-button
          >
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h } from "vue";
import { useRouter } from "vue-router";
import {
  useMessage,
  useDialog,
  NButton,
  NSpace,
  NTag,
  NCard,
  NTabs,
  NTabPane,
  NSpin,
  NDataTable,
  NPagination,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NIcon,
  NDropdown,
} from "naive-ui";
import {
  PersonCircleOutline,
  LogOutOutline,
  RefreshOutline,
  ArrowBackOutline,
  TrashOutline,
  CreateOutline,
  BanOutline,
  CheckmarkCircleOutline,
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

const activeTab = ref("users");
const usersLoading = ref(false);
const linksLoading = ref(false);
const updating = ref(false);
const users = ref([]);
const links = ref([]);
const currentPage = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const showEditLinkModal = ref(false);

const editLinkForm = reactive({
  id: null,
  short_code: "",
  original_url: "",
  user_id: "",
});

const userOptions = [
  {
    label: "退出登录",
    key: "logout",
    icon: () => h(LogOutOutline),
  },
];

function handleUserAction(key) {
  if (key === "logout") {
    userStore.logout();
    router.push("/login");
  }
}

const userColumns = [
  {
    title: "ID",
    key: "ID",
    width: 80,
  },
  {
    title: "用户名",
    key: "Username",
  },
  {
    title: "邮箱",
    key: "Email",
  },
  {
    title: "状态",
    key: "IsBanned",
    render: (row) => {
      return h(
        NTag,
        { type: row.IsBanned ? "error" : "success" },
        { default: () => (row.IsBanned ? "已封禁" : "正常") },
      );
    },
  },
  {
    title: "角色",
    key: "IsAdmin",
    render: (row) => {
      return h(
        NTag,
        { type: row.IsAdmin ? "warning" : "default" },
        { default: () => (row.IsAdmin ? "管理员" : "普通用户") },
      );
    },
  },
  {
    title: "创建时间",
    key: "CreatedAt",
    render: (row) => new Date(row.CreatedAt).toLocaleString("zh-CN"),
  },
  {
    title: "操作",
    key: "actions",
    render: (row) => {
      return h(
        NSpace,
        {},
        {
          default: () => [
            !row.IsBanned &&
              h(
                NButton,
                {
                  size: "small",
                  type: "warning",
                  onClick: () => handleBanUser(row),
                },
                {
                  icon: () => h(BanOutline),
                  default: () => "封禁",
                },
              ),
            row.IsBanned &&
              h(
                NButton,
                {
                  size: "small",
                  type: "success",
                  onClick: () => handleUnbanUser(row),
                },
                {
                  icon: () => h(CheckmarkCircleOutline),
                  default: () => "解封",
                },
              ),
            h(
              NButton,
              {
                size: "small",
                type: "error",
                onClick: () => handleDeleteUser(row),
              },
              {
                icon: () => h(TrashOutline),
                default: () => "删除",
              },
            ),
          ],
        },
      );
    },
  },
];

const linkColumns = [
  {
    title: "ID",
    key: "ID",
    width: 80,
  },
  {
    title: "短链接代码",
    key: "ShortCode",
  },
  {
    title: "原始链接",
    key: "OriginalURL",
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: "用户ID",
    key: "UserID",
    width: 100,
  },
  {
    title: "创建时间",
    key: "CreatedAt",
    render: (row) => new Date(row.CreatedAt).toLocaleString("zh-CN"),
  },
  {
    title: "操作",
    key: "actions",
    render: (row) => {
      return h(
        NSpace,
        {},
        {
          default: () => [
            h(
              NButton,
              {
                size: "small",
                onClick: () => handleEditLink(row),
              },
              {
                icon: () => h(CreateOutline),
                default: () => "编辑",
              },
            ),
            h(
              NButton,
              {
                size: "small",
                type: "error",
                onClick: () => handleDeleteLink(row),
              },
              {
                icon: () => h(TrashOutline),
                default: () => "删除",
              },
            ),
          ],
        },
      );
    },
  },
];

async function loadUsers() {
  try {
    usersLoading.value = true;
    const res = await api.get("/admin/user/list");
    if (res.code === 200) {
      users.value = res.data || [];
    }
  } catch (error) {
    message.error("加载用户列表失败");
  } finally {
    usersLoading.value = false;
  }
}

async function loadLinks() {
  try {
    linksLoading.value = true;
    const res = await api.get("/admin/surl/list", {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
      },
    });
    if (res.code === 200) {
      links.value = res.data.surls || [];
      totalPages.value = Math.ceil(
        (res.data.total_count) / pageSize.value,
      );
    }
  } catch (error) {
    message.error("加载短链接列表失败");
  } finally {
    linksLoading.value = false;
  }
}

function handleBanUser(user) {
  dialog.warning({
    title: "确认封禁",
    content: `确定要封禁用户 ${user.Username} 吗？`,
    positiveText: "封禁",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        const res = await api.post("/admin/user/ban", {
          user_id: user.ID,
        });
        if (res.code === 200) {
          message.success("封禁成功");
          loadUsers();
        } else {
          message.error(res.message || "封禁失败");
        }
      } catch (error) {
        message.error("封禁失败");
      }
    },
  });
}

function handleUnbanUser(user) {
  dialog.info({
    title: "确认解封",
    content: `确定要解封用户 ${user.Username} 吗？`,
    positiveText: "解封",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        const res = await api.post("/admin/user/unban", {
          user_id: user.ID,
        });
        if (res.code === 200) {
          message.success("解封成功");
          loadUsers();
        } else {
          message.error(res.message || "解封失败");
        }
      } catch (error) {
        message.error("解封失败");
      }
    },
  });
}

function handleDeleteUser(user) {
  dialog.error({
    title: "确认删除",
    content: `确定要删除用户 ${user.Username} 吗？此操作不可恢复！`,
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        const res = await api.post("/admin/user/delete", {
          id: user.ID,
        });
        if (res.code === 200) {
          message.success("删除成功");
          loadUsers();
        } else {
          message.error(res.message || "删除失败");
        }
      } catch (error) {
        message.error("删除失败");
      }
    },
  });
}

function handleEditLink(link) {
  editLinkForm.id = link.ID;
  editLinkForm.short_code = link.ShortCode;
  editLinkForm.original_url = link.OriginalURL;
  editLinkForm.user_id = String(link.UserID);
  showEditLinkModal.value = true;
}

async function confirmEditLink() {
  try {
    updating.value = true;
    const res = await api.post("/admin/surl/update", {
      short_code: editLinkForm.short_code,
      original_url: editLinkForm.original_url,
      user_id: editLinkForm.user_id,
    });

    if (res.code === 200) {
      message.success("更新成功");
      showEditLinkModal.value = false;
      loadLinks();
    } else {
      message.error(res.message || "更新失败");
    }
  } catch (error) {
    message.error("更新失败");
  } finally {
    updating.value = false;
  }
}

function handleDeleteLink(link) {
  dialog.error({
    title: "确认删除",
    content: `确定要删除短链接 ${link.ShortCode} 吗？`,
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        const res = await api.post("/admin/surl/delete", {
          id: link.ID,
        });
        if (res.code === 200) {
          message.success("删除成功");
          loadLinks();
        } else {
          message.error(res.message || "删除失败");
        }
      } catch (error) {
        message.error("删除失败");
      }
    },
  });
}

onMounted(() => {
  configStore.loadConfig();
  loadUsers();
  loadLinks();
});
</script>

<style scoped>
.admin {
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
  max-width: 1400px;
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

.user-section {
  display: flex;
  gap: 16px;
  align-items: center;
}

.main-content {
  padding: 24px;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
}

:deep(.n-card) {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

:deep(.n-tabs) {
  background: white;
  padding: 16px;
}

:deep(.n-data-table) {
  font-size: 14px;
}
</style>
