<template>
<b-container fluid>
  <b-row class="hover" v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
    <b-col>
      <b-link v-for="port in obj.spec.ports" :key="port.port" target="_blank" :href="proxy_uri(obj, port.port)">
        {{ obj.metadata.name }}:{{ port.port }}
      </b-link>
    </b-col>
    <b-col md=auto class="text-right">{{ owner(obj) }}</b-col>
    <b-col cols=2 class="text-right">{{ formatTs(obj) }}</b-col>
  </b-row>
</b-container>
</template>

<script>
import SSE from '@/SSE/Objects.vue'

export default {
  extends: SSE,
  props: ['name', 'namespace'],
  data () {
    return {
    }
  },
  methods: {
    proxy_uri(obj, port) {
      return this.$api.uri(`/proxy/${obj.metadata.namespace}/${obj.metadata.name}/${port}`)
    },
    uri() {
      return `/workflow/${this.namespace}/${this.name}/mounts`
    },
  },
}
</script>


