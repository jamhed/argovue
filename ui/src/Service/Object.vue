<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{kind}}/{{namespace}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body v-if="loaded">
        <b-tabs card no-key-nav v-model="tab" @input="onTab" lazy>
          <b-tab title="Proxy">
            <a v-for="port in object.spec.ports" :key="port.port" target="_blank" :href="proxy_uri(port.port)">
              {{ name }}:{{ port.port }}
            </a>
          </b-tab>
          <b-tab title="Tokens">
            <tokens :name="name" :namespace="namespace"></tokens>
          </b-tab>
          <b-tab title="Ingresses">
            <ingresses :name="name" :namespace="namespace"></ingresses>
          </b-tab>
          <b-tab title="Service">
            <jsoneditor :content="object"></jsoneditor>
          </b-tab>
        </b-tabs>
      </b-card>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Object'
import Tokens from '@/Service/Tokens'
import Ingresses from '@/Service/Ingresses'
import JsonEditor from '@/JsonEditor'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    tokens: Tokens,
    ingresses: Ingresses,
  },
  data() {
    return {
      kind: "service",
    }
  },
  methods: {
    proxy_uri(port) {
      return this.$api.uri(`/proxy/${this.namespace}/${this.name}/${port}`)
    },
  },
}
</script>
