<template>
<b-row class="hover">
  <b-col v-if="isPod()">
    <b-link :to="`/k8s/pod/${namespace}/${content.id}`">{{content.displayName}}</b-link>
  </b-col>
  <b-col v-else>
    {{content.displayName}}
  </b-col>
  <b-col md=auto class="text-right">{{content.type}}</b-col>
  <b-col md=auto class="text-right">{{phase(content.phase)}}</b-col>
  <b-col cols=2 class="text-right">{{duration}}s</b-col>
</b-row>
</template>

<script>
import moment from 'moment'
import util from '@/util'
import 'moment-duration-format'

export default {
  props: ['content', 'name', 'namespace'],
  data () {
    return {
      duration: 0
    }
  },
  created () {
    setInterval(() => this.calcDuration(), 1000)
    this.calcDuration()
  },
  methods: {
    isPod () {
      return this.content.type == 'Pod'
    },
    phase(phase) {
      return util.phase(phase)
    },
    calcDuration() {
      let start = moment(this.content.startedAt)
      let end = this.content.finishedAt? moment(this.content.finishedAt) : moment()
      let duration = moment.duration(end.unix() - start.unix(), 'seconds')
      this.duration = duration.format("h:mm:ss")
    },
  }
}
</script>