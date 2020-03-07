<template><div></div></template>

<script>
export default {
  data() {
    return {
      es: undefined,
      teardown: false
    };
  },
  created() {
    this.setupStream()
  },
  destroyed() {
    this.tearDown()
  },
  watch: {
    name () {
      this.tearDown()
      this.setupStream()
    }
  },
  methods: {
    tearDown() {
      if (this.es) {
        this.es.close()
      }
      this.object = {}
      this.es = undefined
    },
    setupStream() {
      this.es = this.$api.sse("/events", (ev) => {
        this.teardown = false
        return ev
      })
      this.es.onerror = () => {
        this.teardown = true
        setTimeout(() => {
          if (this.teardown) {
            this.$log("disconnect on connection error", this)
            this.$api.logout()
          }
        }, 30000)
      }
    }
  }
};
</script>