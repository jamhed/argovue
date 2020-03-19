<template>
<b-container>
  <b-row>
    <b-col class="mb-3">
      <b-button @click="create" size="sm" variant="primary">Sync</b-button>
    </b-col>
  </b-row>
  <b-row class="hover" v-for="obj in orderedCache" :key="obj.metadata.uid">
    <b-col cols=2>
      <b-link :to="`/k8s/job/${obj.metadata.namespace}/${obj.metadata.name}`">{{ obj.metadata.name }}</b-link>
    </b-col>
    <b-col class="text-right">{{ obj.metadata.creationTimestamp }}</b-col>
    <b-col md="auto">
      <b-dropdown variant="link" toggle-class="p-0">
        <b-dropdown-item-button @click="del(obj.metadata.name)">Delete</b-dropdown-item-button>
      </b-dropdown>
    </b-col>
  </b-row>
</b-container>
</template>

<script>
import SSE from '@/SSE/Objects'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  data() {
    return {
    }
  },
  methods: {
    create: async function () {
      let re = await this.$api.post(`/k8s/datasource/${this.namespace}/${this.name}/syncs/create`)
      this.showReply(re)
    },
    del: async function (name) {
      let re = await this.$api.post(`/k8s/datasource/${this.namespace}/${this.name}/sync/${name}/delete`)
      this.showReply(re)
    },
    uri() {
      return `/k8s/datasource/${this.namespace}/${this.name}/syncs`
    },
  },
}
</script>
