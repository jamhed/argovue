<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}/{{name}}</h1>
    </div>
    <div>
      <control :object="object" :name="name" :namespace="namespace" style="margin-bottom: 20px"></control>
      <b-card no-body>
        <b-tabs card no-key-nav v-model="tab" @input="onTab">
          <b-tab title="Nodes">
            <nodes :content="object"></nodes>
          </b-tab>
          <b-tab title="Mounts" lazy>
            <mounts :name="name" :namespace="namespace"></mounts>
          </b-tab>
          <b-tab title="Services" lazy>
            <services :name="name" :namespace="namespace"></services>
          </b-tab>
          <b-tab title="Graph" lazy no-key-nav>
            <graph :content="object" :name="name" :namespace="namespace"></graph>
          </b-tab>
          <b-tab title="DAG" lazy no-key-nav>
            <dag :content="object" :name="name" :namespace="namespace"></dag>
          </b-tab>
          <b-tab title="Workflow">
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
import Nodes from '@/Workflow/Nodes'
import Services from '@/Workflow/Services'
import Mounts from '@/Workflow/Mounts'
import Control from '@/Workflow/Control'
import Graph from '@/Workflow/Graph'
import DAG from '@/Workflow/DAG'

export default {
  props: ['namespace', 'name'],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    nodes: Nodes,
    control: Control,
    services: Services,
    mounts: Mounts,
    graph: Graph,
    dag: DAG,
  },
  data() {
    return {
      kind: 'workflows'
    }
  },
  methods: {
    uri() {
      return `/workflow/${this.namespace}/${this.name}`
    },
    parent() {
      return `/watch/workflows`
    },
  }
}
</script>
