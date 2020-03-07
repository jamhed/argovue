<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">job/{{namespace}}/{{name}}</h1>
    </div>
    <div>
      <b-card no-body v-if="loaded">
        <b-tabs card no-key-nav v-model="tab" @input="onTab">
          <b-tab title="Job">
            <b-row>
              <b-col cols="2">{{object.metadata.name}}</b-col>
              <b-col class="text-right">{{object.status.startTime}}</b-col>
              <b-col cols="auto" class="text-right">{{calcDuration()}}s</b-col>
            </b-row>
          </b-tab>
          <b-tab title="Pods">
            <pod :namespace="namespace" :name="name"></pod>
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
import Pod from '@/Job/Pod'
import moment from 'moment'

export default {
  props: ["namespace", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor,
    pod: Pod,
  },
  methods: {
    uri() {
      return `/k8s/job/${this.namespace}/${this.name}`
    },
    calcDuration() {
      let start = moment(this.object.status.startTime)
      let end = this.object.status.completionTime? moment(this.object.status.completionTime) : moment()
      let duration = moment.duration(end.unix() - start.unix(), 'seconds')
      return duration.format("h:mm:ss")
    },
  },
}
</script>
