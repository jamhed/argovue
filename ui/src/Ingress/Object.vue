<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">ingress/{{namespace}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body v-if="loaded">
        <b-tabs card no-key-nav v-model="tab" @input="onTab">
          <b-tab title="Ingress">
            <b-row>
              <b-col>
                <b-link target="_new" :href="`https://${object.spec.rules[0].host}`">{{object.spec.rules[0].host}}</b-link>
              </b-col>
              <b-col cols="auto">{{object.metadata.creationTimestamp}}</b-col>
            </b-row>
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

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
  },
  methods: {
    uri() {
      return `/k8s/ingress/${this.namespace}/${this.name}`
    }
  },
}
</script>
