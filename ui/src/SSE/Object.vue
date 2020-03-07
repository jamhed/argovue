<script>
export default {
  data() {
    return {
      loaded: false,
      tab: 0,
      object: {},
      es: undefined
    };
  },
  created() {
    this.tab = parseInt(this.get('tab') || 0)
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
    uri() {
      return `/k8s/${this.kind}/${this.namespace}/${this.name}`
    },
    parent() {
      return `/k8s/${this.kind}`
    },
    onTab (id) {
      this.set('tab', id)
    },
    set(lkey, value) {
      let key = `${this.uri()}/${lkey}`
      localStorage.setItem(key, value)
    },
    get(lkey) {
      let key = `${this.uri()}/${lkey}`
      return localStorage.getItem(key)
    },
    setupStream() {
      this.es = this.$api.sse(this.uri(), (event) => {
        let msg = JSON.parse(event.data);
        let obj = msg.Content;
        switch (msg.Action) {
          case "delete":
            this.$router.replace(this.parent())
            break
          case "add":
            this.$set(this, "object", obj)
            this.loaded = true
            break
          case "update":
            this.$set(this, "object", obj)
            this.loaded = true
            break
        }
      })
    }
  }
};
</script>