<template>
<b-container>
  <b-row>
    <b-col class="mb-3">
      <b-button @click="create" size="sm" variant="primary">Create</b-button>
    </b-col>
  </b-row>
  <b-row class="hover" v-for="obj in orderedCache" :key="obj.metadata.uid">
    <b-col class="text-left">
      <b-link :to="`/k8s/ingress/${obj.metadata.namespace}/${obj.metadata.name}`">{{ obj.metadata.name }}</b-link>
    </b-col>
    <b-col class="text-left">
      <b-link target="_new" :href="`https://${obj.spec.rules[0].host}`">link</b-link>
    </b-col>
    <b-col cols=2 sm="3">{{ obj.metadata.creationTimestamp }}</b-col>
    <b-col cols="auto" class="text-right">
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
      kind: "ingress",
    }
  },
  methods: {
    create: async function () {
      let re = await this.$api.post(`/k8s/service/${this.namespace}/${this.name}/ingresses/create`)
      this.showReply(re)
    },
    del: async function (name) {
      let re = await this.$api.post(`/k8s/service/${this.namespace}/${this.name}/ingress/${name}/delete`)
      this.showReply(re)
    },
    uri() {
      return `/k8s/service/${this.namespace}/${this.name}/ingresses`
    },
  },
}
</script>

