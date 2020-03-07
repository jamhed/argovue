<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{kind}}</h1>
    </div>
    <b-container>
      <b-row class="hover" v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
        <b-col>
          <b-link :to="`/${kind}/${obj.metadata.namespace}/${obj.metadata.name}`">{{obj.metadata.namespace}}/{{ obj.metadata.name }}</b-link>
        </b-col>
        <b-col md="auto" class="text-right">{{ owner(obj) }}</b-col>
        <b-col md="auto" v-if="obj.status" class="text-right">{{ phase(obj) }}</b-col>
        <b-col cols=2 class="text-right">{{ formatTs(obj) }}</b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import SSE from '@/SSE/Objects'

export default {
  props: ["kind"],
  extends: SSE,
  data() {
    return {
    }
  },
  methods: {
  },
};
</script>
