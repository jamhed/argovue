<template>
  <b-container>
    <b-row class="hover" v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
      <b-col>
        <b-link :to="`/catalogue/${namespace}/${name}/instance/${obj.metadata.name}`">{{ obj.metadata.name }}</b-link>
      </b-col>
      <b-col md="auto" class="text-right">{{ owner(obj) }}</b-col>
      <b-col cols=2 class="text-right">{{ formatTs(obj) }}</b-col>
      <b-col md="auto">
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
  props: ['name', 'kind', 'namespace'],
  extends: SSE,
  compnents: {
  },
  data () {
    return {
    }
  },
  methods: {
    uri() {
      return `/catalogue/${this.namespace}/${this.name}/instances`
    },
    del(instance) {
      this._action(`/catalogue/${this.namespace}/${this.name}/instance/${instance}/action/delete`)
    },
  },
}
</script>