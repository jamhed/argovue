<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">pvc/{{namespace}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body v-if="loaded">
        <b-tabs card no-key-nav v-model="tab" @input="onTab">
          <b-tab title="PVC">
            <b-row>
              <b-col cols="auto">{{object.metadata.creationTimestamp}}</b-col>
              <b-col v-if="object.status" cols="auto">{{object.status.phase}}</b-col>
              <b-col v-if="object.status.capacity" cols="auto">{{object.status.capacity.storage}}</b-col>
              <b-col v-if="object.status" cols="auto">{{object.status.accessModes}}</b-col>
            </b-row>
          </b-tab>
          <b-tab title="Mounts" lazy>
            <mounts :name="name" :namespace="namespace"></mounts>
          </b-tab>
          <b-tab title="Datasource">
            <datasource :namespace="namespace" :name="name"></datasource>
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
import Datasource from '@/Pvc/Datasource'
import Mounts from '@/Pvc/Mounts'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    datasource: Datasource,
    mounts: Mounts,
  },
  methods: {
    uri() {
      return `/k8s/pvc/${this.namespace}/${this.name}`
    }
  },
}
</script>
