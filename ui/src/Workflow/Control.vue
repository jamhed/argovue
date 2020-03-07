<template>
<b-button-toolbar>
  <b-button-group size="sm" class="mr-1">
    <b-button :disabled="cantRetry()" @click="retry()">Retry</b-button>
    <b-button @click="resubmit()">Resubmit</b-button>
    <b-button @click="del()">Delete</b-button>
  </b-button-group>
  <b-button-group size="sm" class="mr-1">
    <b-button :disabled="cantSuspend()" @click="suspend()">Suspend</b-button>
    <b-button :disabled="cantResume()" @click="resume()">Resume</b-button>
    <b-button :disabled="cantTerminate()" @click="terminate()">Terminate</b-button>
  </b-button-group>
  <b-button-group size="sm" class="mr-1">
    <b-button :disabled="cantMount()" @click="mount()">Mount</b-button>
  </b-button-group>
</b-button-toolbar>
</template>

<script>
export default {
  props: ['object', 'name', 'namespace'],
  data () {
    return {
      nodes: []
    }
  },
  methods: {
    status (status) {
      return this.object && this.object.status && this.object.status.phase == status
    },
    _action: async function(uri) {
      let re = await this.$api.post(uri)
      this.$bvToast.toast(`${re.data.action} ${re.data.status} ${re.data.message}`, {
        title: re.data.action,
        toaster: 'b-toaster-bottom-right',
        autoHideDelay: re.data.status == 'ok' ? 3000 : 6000,
        noCloseButton: true,
        variant: re.data.status == 'ok'? 'info' : 'danger'
    })},
    action: async function(action) {
      this.$bvModal.msgBoxConfirm(`Please confirm workflow action ${action}:`, { buttonSize: "sm" }).then(value => {
          if (value) {
            this._action(`/workflow/${this.namespace}/${this.name}/action/${action}`)
          }
        })
    },
    cantRetry () {
      return ! (this.status('Failed') || this.status('Error'))
    },
    cantSuspend () {
      return ! (this.status('Running') || this.isSuspended())
    },
    cantResume () {
      return ! this.isSuspended()
    },
    cantTerminate () {
      return ! this.status('Running')
    },
    isSuspended () {
      return this.object && this.object.spec && this.object.spec.suspend
    },
    cantMount () {
      return false
    },
    retry () {
      this.action('retry')
    },
    del () {
      this.action('delete')
    },
    resubmit () {
      this.action('resubmit')
    },
    suspend () {
      this.action('suspend')
    },
    resume () {
      this.action('resume')
    },
    terminate () {
      this.action('terminate')
    },
    mount () {
      this.action('mount')
    }
  },
}
</script>
