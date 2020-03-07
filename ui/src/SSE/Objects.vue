<script>
import moment from 'moment'
import util from '@/util.js'

function objKey(obj) {
  let key = obj.kind + obj.metadata.namespace + obj.metadata.name
  return key
}

export default {
  data() {
    return {
      cache: {},
      es: undefined
    };
  },
  created() {
    this.setupStream()
  },
  destroyed() {
    this.tearDown()
  },
  computed: {
    orderedCache: function() {
      return Object.values(this.cache).sort( (a, b) => objKey(a).localeCompare(objKey(b)) )
    }
  },
  watch: {
    kind () {
      this.tearDown()
      this.setupStream()
    }
  },
  methods: {
    uri() {
      return `/watch/${this.kind}`
    },
    tearDown() {
      if (this.es) {
        this.es.close()
      }
      this.cache = {}
      this.es = undefined
    },
    _action: async function(uri) {
      let re = await this.$api.post(uri)
      this.showReply(re)
    },
    showReply(re) {
      this.$bvToast.toast(`${re.data.action} ${re.data.status} ${re.data.message}`, {
        title: re.data.action,
        toaster: 'b-toaster-bottom-right',
        autoHideDelay: re.data.status == 'ok' ? 3000 : 6000,
        noCloseButton: true,
        variant: re.data.status == 'ok'? 'info' : 'danger'
      })
    },
    formatTs(obj) {
      return moment(obj.metadata.creationTimestamp).format("YYYY-MM-DD HH:mm:ss")
    },
    owner (obj) {
      return util.owner(obj)
    },
    phase (obj) {
      return util.phase(obj.status.phase)
    },
    setupStream() {
      this.es = this.$api.sse(this.uri(), (event) => {
        var msg = JSON.parse(event.data)
        var obj = msg.Content
        switch (msg.Action) {
          case "delete":
            this.$delete(this.cache, obj.metadata.uid)
            break
          case "add":
            this.$set(this.cache, obj.metadata.uid, obj)
            break
          case "update":
            this.$set(this.cache, obj.metadata.uid, obj)
            break
        }
      })
    }
  }
}
</script>