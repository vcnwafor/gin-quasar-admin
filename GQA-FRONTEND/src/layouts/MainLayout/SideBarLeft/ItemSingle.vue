<template>
    <q-item clickable v-ripple exact @click="toPath(trueItem)" :inset-level="initLevel" :active="checkActive"
        :active-class="darkThemeSelect">
        <q-item-section avatar>
            <q-icon :name="trueItem.icon" />
        </q-item-section>
        <q-item-section>{{ $t(trueItem.title) }}</q-item-section>
    </q-item>
</template>

<script setup>
import { computed, toRefs } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import useDarkTheme from 'src/composables/useDarkTheme';

const { darkThemeSelect } = useDarkTheme()
const router = useRouter()
const route = useRoute()
const props = defineProps({
    trueItem: {
        default: function () {
            return {}
        },
        type: Object,
    },
    initLevel: {
        type: Number,
        default: 0,
    },
})
const { trueItem, initLevel } = toRefs(props)

const toPath = (item) => {
    if (item.is_link === 'yes') {
        window.open(item.path)
    } else {
        router.push(item.path)
    }
}

const checkActive = computed(() => {
    if (route.path === trueItem.value.path) {
        return true
    } else {
        return false
    }
})
</script>

<style lang="scss" scoped>
.item-active-class {
    color: $primary;
    background: $primary;
    // background: lighten($color: $primary, $amount: 30%);
    // background: scale-color($primary, $lightness: 5%);
}
</style>