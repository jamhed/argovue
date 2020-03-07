<template>
<b-container>
  <b-row class="hover" v-for="obj in orderedCache" :key="obj.metadata.uid">
    <b-col cols=2>
      <b-link :to="`/k8s/pod/${obj.metadata.namespace}/${obj.metadata.name}`">{{ obj.metadata.name }}</b-link>
    </b-col>
    <b-col>{{ obj.spec.location }}</b-col>
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
    uri() {
      return `/k8s/job/${this.namespace}/${this.name}/pods`
    },
  },
}
</script>
