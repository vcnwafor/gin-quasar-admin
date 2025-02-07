<template>
    <div class="items-center column">
        <div class="justify-between row" style="width: 100%">
            <q-btn color="negative" @click="handleClear">
                {{ $t('Clear') + $t('All') }}
            </q-btn>
            <q-btn color="negative" @click="handleAll">
                {{ $t('Select') + $t('All') }}
            </q-btn>
            <q-btn color="primary" @click="handleRoleMenu" :disable="row.role_code === 'super-admin'">
                {{ $t('Save') }}
            </q-btn>
        </div>
        <q-card-section style="width: 100%; max-height: 70vh" class="scroll">
            <q-tree dense style="width: 100%" :nodes="menuTree" default-expand-all node-key="name" label-key="name"
                selected-color="primary" v-if="menuTree.length !== 0" tick-strategy="strict" v-model:ticked="ticked">
                <template v-slot:default-header="prop">
                    <div class="items-center row">
                        <q-icon :name="prop.node.icon || 'share'" size="28px" class="q-mr-sm" />
                        <div class="text-weight-bold">{{ $t(prop.node.title) }}</div>
                    </div>
                </template>
            </q-tree>
        </q-card-section>
        <q-inner-loading :showing="loading">
            <q-spinner-gears size="50px" color="primary" />
        </q-inner-loading>
    </div>

</template>

<script setup>
import useTableData from 'src/composables/useTableData'
import { useQuasar } from 'quasar'
import { postAction } from 'src/api/manage'
import { computed, onMounted, ref, toRefs } from 'vue'

const $q = useQuasar()
const url = {
    list: 'menu/get-menu-list',
    roleMenuList: 'role/get-role-menu-list',
    roleMenuEdit: 'role/edit-role-menu',
}

const props = defineProps({
    row: {
        type: Object,
        required: true,
    }
})
const { row } = toRefs(props)

const {
    pagination,
    loading,
    tableData,
    getTableData,
} = useTableData(url)

const menuTree = computed(() => {
    // 如果是super-admin角色，那么禁用角色管理编辑，保证可以访问此界面。
    if (row.value.role_code === 'super-admin' && tableData.value.length) {
        for (let m of tableData.value[0].children) {
            if (m.name === 'role') {
                m.disabled = true
            }
        }
    }
    return tableData.value
})
onMounted(() => {
    pagination.value.rowsPerPage = 99999
    getTableData()
    getRoleMenuList()
})
const ticked = ref([])
const getRoleMenuList = () => {
    // 每次获取前，清空ticked
    ticked.value = []
    postAction(url.roleMenuList, {
        role_code: row.value.role_code,
    }).then((res) => {
        if (res.code === 1) {
            res.data.records.forEach((item) => {
                ticked.value.push(item.sys_menu_name)
            })
        }
    })
}
const handleRoleMenu = () => {
    const roleMenu = []
    for (let i of ticked.value) {
        roleMenu.push({
            role_code: row.value.role_code,
            menu_name: i,
        })
    }
    postAction(url.roleMenuEdit, {
        role_code: row.value.role_code,
        role_menu: roleMenu,
    }).then((res) => {
        if (res.code === 1) {
            $q.notify({
                type: 'positive',
                message: res.message,
            })
            getRoleMenuList()
        }
    })
}
const handleClear = () => {
    ticked.value = []
}
const handleAll = () => {
    tableData.value.forEach((item) => {
        ticked.value.push(item.name)
    })
}
</script>
