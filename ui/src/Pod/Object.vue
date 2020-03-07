<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{kind}}/{{namespace}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body v-if="loaded">
        <b-tabs card no-key-nav v-model="tab" @input="onTab">
          <b-tab title="Pod">
            <b-row class="mb-4">
              <b-col cols=1>Phase:</b-col>
              <b-col cols=11>{{object.status.phase}}</b-col>
            </b-row>
            <b-row v-for="(c,index) in object.status.conditions" :key="index">
              <b-col cols=2>{{c.lastTransitionTime}}</b-col>
              <b-col cols=2>{{c.type}}</b-col>
              <b-col cols=2>{{c.reason}}</b-col>
              <b-col md="auto">{{c.message}}</b-col>
            </b-row>
          </b-tab>
          <b-tab title="Object">
            <jsoneditor :content="object"></jsoneditor>
          </b-tab>
          <b-tab v-for="container in containers" :key="container.name" :title="container.name" lazy>
            <logs :name="name" :namespace="namespace" :container="container.name"></logs>
          </b-tab>
        </b-tabs>
      </b-card>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Object'
import JsonEditor from '@/JsonEditor'
import Logs from '@/Pod/Logs'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    logs: Logs,
  },
  data() {
    return {
      loaded: false,
      kind: "pod",
      containers: [],
    }
  },
  watch: {
    object (c) {
      this.containers = c.spec.containers
    }
  }
}
</script>
