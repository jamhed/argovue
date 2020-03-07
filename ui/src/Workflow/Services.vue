<template>
<b-container fluid>
  <b-row class="hover" v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
    <b-col>
      <b-link :to="`/${kind}/${obj.metadata.namespace}/${obj.metadata.labels['service.argovue.io/name']}/instance/${obj.metadata.name}`">
        {{obj.metadata.namespace}}/{{ obj.metadata.name }}
      </b-link>
    </b-col>
    <b-col md=auto class="text-right">{{ owner(obj) }}</b-col>
    <b-col cols=2 class="text-right">{{ formatTs(obj) }}</b-col>
    <b-col md=auto>
      <b-dropdown variant="link" toggle-class="p-0">
        <b-dropdown-item-button @click="del(obj.metadata.name)">Delete</b-dropdown-item-button>
      </b-dropdown>
    </b-col>
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
      kind: "catalogue",
    }
  },
  methods: {
    uri() {
      return `/workflow/${this.namespace}/${this.name}/services`
    },
    del(instance) {
      return this._action(`/workflow/${this.namespace}/${this.name}/service/${instance}/action/delete`)
    },
  },
}
</script>


