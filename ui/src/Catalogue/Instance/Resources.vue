<template>
  <b-container>
    <b-row class="hover" v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
      <b-col>
        <b-link :to="`/k8s/${obj.kind.toLowerCase()}/${obj.metadata.namespace}/${obj.metadata.name}`">
          {{ obj.kind.toLowerCase() }}/{{ obj.metadata.namespace }}/{{ obj.metadata.name }}
        </b-link>
      </b-col>
      <b-col md="auto" class="text-right">{{ owner(obj) }}</b-col>
      <b-col md="auto" class="text-right" v-if="obj.status">{{ phase(obj) }}</b-col>
      <b-col cols=2 class="text-right">{{ formatTs(obj) }}</b-col>
    </b-row>
  </b-container>
</template>

<script>
import SSE from '@/SSE/Objects.vue'

export default {
  props: ['name', 'namespace', 'instance'],
  extends: SSE,
  data () {
    return {
    }
  },
  methods: {
    uri() {
      return `/catalogue/${this.namespace}/${this.name}/instance/${this.instance}/resources`
    },
  },
}
</script>