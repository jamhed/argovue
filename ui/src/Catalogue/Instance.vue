<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{instance}}</h1>
    </div>
    <div>
      <b-card no-body v-if="loaded">
        <b-tabs card lazy v-model="tab" @input="onTab">
          <b-tab title="Instance">
            <b-row class="mb-4">
              <b-col cols=1>Status:</b-col>
              <b-col cols=11>{{object.status.releaseStatus}}</b-col>
            </b-row>
            <b-row v-for="(c,index) in object.status.conditions" :key="index">
              <b-col cols=2>{{c.lastTransitionTime}}</b-col>
              <b-col cols=2>{{c.type}}</b-col>
              <b-col cols=2>{{c.reason}}</b-col>
              <b-col md="auto">{{c.message}}</b-col>
            </b-row>
          </b-tab>
          <b-tab title="Resources">
            <resources :name="name" :namespace="namespace" :instance="instance"></resources>
          </b-tab>
          <b-tab title="Object">
            <jsoneditor :content="object"></jsoneditor>
          </b-tab>
        </b-tabs>
      </b-card>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Object'
import JsonEditor from '@/JsonEditor'
import Resources from '@/Catalogue/Instance/Resources'

export default {
  props: ['namespace', 'name', 'instance'],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    resources: Resources,
  },
  methods: {
    uri() {
      return `/catalogue/${this.namespace}/${this.name}/instance/${this.instance}`
    },
  },
  data() {
    return {
      kind: "services",
      object: {
        spec: {
          ports: [],
        }
      }
    }
  },
}
</script>
