<template>
  <nav class="fe-tree" :class="{ 'fe-tree--collapsed': collapsed }" aria-label="Directory tree">
    <div class="fe-tree__title">
      <i class="pi pi-folder-open" v-if="!collapsed" />
      <span v-if="!collapsed">Files</span>
      <div v-if="!collapsed" style="flex:1" />
      <button class="fe-tree__collapse-btn" :title="collapsed ? 'Expand tree' : 'Collapse tree'" @click.stop="collapsed = !collapsed">
        <i class="pi" :class="collapsed ? 'pi-chevron-right' : 'pi-chevron-left'" />
      </button>
    </div>

    <template v-if="!collapsed">
      <!-- Root entry -->
      <div
        class="fe-tree__item"
        :class="{ 'fe-tree__item--active': currentPath === '' }"
        style="padding-left: 8px"
        @click="emit('navigate', '')"
      >
        <i class="pi pi-home fe-tree__icon" />
        <span class="fe-tree__name">/</span>
      </div>

      <!-- Flat list of directory tree nodes (directories + files) -->
      <div
        v-for="item in flatTree"
        :key="item.path"
        class="fe-tree__item"
        :class="{ 'fe-tree__item--active': item.isDir ? currentPath === item.path : false }"
        :style="{ paddingLeft: `${(item.depth + 1) * 14 + 8}px` }"
        @click="item.isDir ? emit('navigate', item.path) : emit('open-file', item.path, item.name)"
      >
        <!-- Expand toggle: dirs only -->
        <button
          v-if="item.isDir"
          class="fe-tree__toggle-btn"
          :aria-expanded="item.expanded"
          :aria-label="item.expanded ? `Collapse ${item.name}` : `Expand ${item.name}`"
          @click.stop="toggleExpand(item)"
        >
          <i
            class="pi fe-tree__toggle"
            :class="item.loading ? 'pi-spin pi-spinner' : item.expanded ? 'pi-chevron-down' : 'pi-chevron-right'"
            aria-hidden="true"
          />
        </button>
        <span v-else class="fe-tree__toggle" />
        <!-- Icon -->
        <i class="pi fe-tree__icon" :class="item.isDir ? 'pi-folder' : 'pi-file'" />
        <span class="fe-tree__name">{{ item.name }}</span>
      </div>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { FileEntry } from '../../api/files'

interface TreeItem {
  path: string
  name: string
  depth: number
  expanded: boolean
  loading: boolean
  isDir: boolean
}

const props = defineProps<{
  currentPath: string
  fetchFiles: (path: string) => Promise<{ entries: FileEntry[] }>
}>();

const emit = defineEmits<{
  (e: 'navigate', path: string): void
  (e: 'open-file', path: string, name: string): void
}>()

const flatTree = ref<TreeItem[]>([])
const collapsed = ref(false)
// Cache: parent path → loaded child tree items
const childCache = new Map<string, TreeItem[]>()

async function loadChildren(parentPath: string, depth: number): Promise<TreeItem[]> {
  if (childCache.has(parentPath)) return childCache.get(parentPath)!
  try {
    const { entries } = await props.fetchFiles(parentPath)
    const items: TreeItem[] = entries.map(e => ({
      path: parentPath ? `${parentPath}/${e.name}` : e.name,
      name: e.name,
      depth,
      expanded: false,
      loading: false,
      isDir: e.isDir,
    }))
    // Dirs first, then files — both sorted alphabetically
    items.sort((a, b) => {
      if (a.isDir !== b.isDir) return a.isDir ? -1 : 1
      return a.name.localeCompare(b.name)
    })
    childCache.set(parentPath, items)
    return items
  } catch {
    return []
  }
}

async function toggleExpand(item: TreeItem) {
  if (!item.isDir || item.loading) return
  if (item.expanded) {
    // Collapse: remove all descendants from flat list
    item.expanded = false
    const idx = flatTree.value.findIndex(i => i.path === item.path)
    if (idx !== -1) {
      let end = idx + 1
      while (end < flatTree.value.length && flatTree.value[end].depth > item.depth) end++
      flatTree.value.splice(idx + 1, end - idx - 1)
    }
  } else {
    // Expand: load children and insert after this item
    item.loading = true
    const children = await loadChildren(item.path, item.depth + 1)
    item.loading = false
    item.expanded = true
    const idx = flatTree.value.findIndex(i => i.path === item.path)
    if (idx !== -1) flatTree.value.splice(idx + 1, 0, ...children)
  }
}

onMounted(async () => {
  flatTree.value = await loadChildren('', 0)
})
</script>

<style scoped>
.fe-tree {
  width: 220px;
  min-width: 220px;
  overflow-y: auto;
  overflow-x: hidden;
  background: var(--bg, #0c0e18);
  border-right: 1px solid var(--border, #252a42);
  font-size: 0.8rem;
  font-family: Lato, sans-serif;
  user-select: none;
  color: var(--text, #dde3f8);
  transition: width 0.2s ease, min-width 0.2s ease;
  flex-shrink: 0;
}

.fe-tree--collapsed {
  width: 32px;
  min-width: 32px;
  overflow: hidden;
}

.fe-tree__title {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 12px 6px;
  font-weight: 600;
  color: var(--muted, #5a6490);
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  white-space: nowrap;
  overflow: hidden;
}

.fe-tree--collapsed .fe-tree__title {
  padding: 8px 0;
  justify-content: center;
}

.fe-tree__collapse-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--muted, #5a6490);
  border-radius: 3px;
  padding: 0;
  transition: color 0.12s, background 0.12s;
}
.fe-tree__collapse-btn:hover { color: var(--text, #dde3f8); background: var(--surface2, #1c2038); }
.fe-tree__collapse-btn .pi { font-size: 11px; }

.fe-tree__item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 8px;
  cursor: pointer;
  border-radius: 4px;
  color: var(--muted, #5a6490);
  transition: background 0.1s, color 0.1s;
  margin: 0 4px;
}
.fe-tree__item:hover { background: var(--surface2, #1c2038); color: var(--text, #dde3f8); }
.fe-tree__item--active {
  background: var(--sel-bg, rgba(79,142,247,0.1));
  color: var(--accent, #4f8ef7);
}

.fe-tree__toggle {
  width: 14px;
  font-size: 10px;
  flex-shrink: 0;
  color: var(--muted, #5a6490);
}

.fe-tree__toggle-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
  color: var(--muted, #5a6490);
}

.fe-tree__icon {
  font-size: 13px;
  color: var(--warn, #f59e0b);
  flex-shrink: 0;
}

.fe-tree__name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
}
</style>
