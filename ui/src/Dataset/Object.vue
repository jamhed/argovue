<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">datasource/{{namespace}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body v-if="loaded">
        <b-tabs card no-key-nav v-model="tab" @input="onTab">
          <b-tab title="Datasource">
            <b-row>
              <b-col>{{object.spec.location}}</b-col>
              <b-col class="text-right">{{object.metadata.creationTimestamp}}</b-col>
            </b-row>
          </b-tab>
          <b-tab title="Sync">
            <sync :namespace="object.metadata.namespace" :name="object.metadata.name"></sync>
          </b-tab>
          <b-tab title="PVCs">
            <pvcs :namespace="object.metadata.namespace" :name="object.metadata.name"></pvcs>
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
import Sync from '@/Datasource/Sync'
import PVCs from '@/Datasource/PVCs'
import JsonEditor from '@/JsonEditor'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    sync: Sync,
    pvcs: PVCs
  },
  methods: {
    uri() {
      return `/k8s/datasource/${this.namespace}/${this.name}`
    }
  },
}
</script>
