<template>
<b-container>
  <b-row>
    <b-col class="mb-3">
      <b-button @click="create" size="sm" variant="primary">Create</b-button>
    </b-col>
  </b-row>
  <b-row class="hover" v-for="obj in orderedCache" :key="obj.metadata.uid">
    <b-col class="text-left">{{ obj.spec.value }}</b-col>
    <b-col cols=3>{{ obj.metadata.name }}</b-col>
    <b-col cols=2>{{ obj.metadata.creationTimestamp }}</b-col>
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
      kind: "token",
    }
  },
  methods: {
    create: async function () {
      await this.$api.post(`/k8s/service/${this.namespace}/${this.name}/tokens/create`)
    },
    del: async function (name) {
      await this.$api.post(`/k8s/service/${this.namespace}/${this.name}/token/${name}/delete`)
    },
    uri() {
      return `/k8s/service/${this.namespace}/${this.name}/tokens`
    },
  },
}
</script>

