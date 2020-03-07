<template>
<div>
  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Version</h1>
  </div>
  <b-container>
    <b-row>
      <b-col cols=3>App</b-col>
      <b-col md=auto>{{version.version}}</b-col>
      <b-col md=auto>{{version.builddate}}</b-col>
      <b-col md=auto>{{version.commit}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=3>Kubernetes</b-col>
      <b-col md=auto>{{version.kubernetes.gitVersion}}</b-col>
      <b-col md=auto>{{version.kubernetes.buildDate}}</b-col>
    </b-row>
  </b-container>
  <b-container class="mt-2">
    <b-row><b-col><b>Group maps:</b></b-col></b-row>
    <b-row v-for="(oidc, k8s) in groups" v-bind:key="oidc">
      <b-col cols="2">{{oidc}}</b-col>
      <b-col cols="auto">{{k8s}}</b-col>
    </b-row>
  </b-container>
</div>
</template>

<script>
export default {
  data () {
    return {
      groups: {},
      version: {
        kubernetes: {}
      }
    }
  },
  created: async function() {
    var re = await this.$api.get("/version")
    this.version = re.data.version
    this.groups = re.data.groups
  }
}
</script>
