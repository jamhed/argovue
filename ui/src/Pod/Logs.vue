<template>
<term :lines="logs"></term>
</template>

<script>
import Term from '@/Term'

export default {
  props: ["namespace", "name", "container"],
  components: {
    term: Term
  },
  data() {
    return {
      logs: [],
      es: undefined,
    }
  },
  created: async function() {
    this.es = this.$api.sse(`/k8s/pod/${this.namespace}/${this.name}/container/${this.container}/logs`,
      (event) => {
        this.logs.push(event.data)
      }
    )
  },
  destroyed () {
    if (this.es) {
      this.es.close()
    }
    this.es = undefined
  },
}
</script>