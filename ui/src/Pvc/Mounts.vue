<template>
<b-container fluid>
  <b-row>
    <b-col>
      <b-form @submit="onSubmit">
        <b-form-group label="Owner" label-for="owner">
          <b-form-select id="owner" v-model="mountOwner" :options="owners()"></b-form-select>
        </b-form-group>
        <b-form-checkbox id="canonical" v-model="canonical" name="canonical" value="true" unchecked-value="false">
          Canonical name
        </b-form-checkbox>
        <b-button class="mt-2" type="submit" size="sm" variant="primary">Deploy</b-button>
      </b-form>
    </b-col>
  </b-row>
  <b-row class="hover mt-3" v-for="obj in orderedCache" v-bind:key="obj.metadata.uid">
    <b-col>
      <b-link v-for="port in obj.spec.ports" :key="port.port" target="_blank" :href="proxy_uri(obj, port.port)">
        {{ obj.metadata.name }}:{{ port.port }}
      </b-link>
    </b-col>
    <b-col md=auto class="text-right">{{ owner(obj) }}</b-col>
    <b-col cols=2 class="text-right">{{ formatTs(obj) }}</b-col>
    <b-col md="auto">
      <b-dropdown variant="link" toggle-class="p-0">
        <b-dropdown-item-button @click="del(obj.metadata.name)">Delete</b-dropdown-item-button>
      </b-dropdown>
    </b-col>
  </b-row>
</b-container>
</template>

<script>
import SSE from '@/SSE/Objects.vue'

export default {
  extends: SSE,
  props: ['name', 'namespace'],
  data () {
    return {
      mountOwner: this.$api.effective_id(),
      canonical: "false",
    }
  },
  methods: {
    proxy_uri(obj, port) {
      return this.$api.uri(`/proxy/${obj.metadata.namespace}/${obj.metadata.name}/${port}`)
    },
    uri() {
      return `/k8s/pvc/${this.namespace}/${this.name}/mounts`
    },
    owners() {
      var owners = this.$api.effective_groups()
      owners.push(this.$api.effective_id())
      return owners
    },
    del: async function (name) {
      let re = await this.$api.post(`/k8s/pvc/${this.namespace}/${this.name}/mount/${name}/delete`)
      this.showReply(re)
    },
    onSubmit: async function(ev) {
      ev.preventDefault()
      let re = await this.$api.post2(`/k8s/pvc/${this.namespace}/${this.name}/mounts/create`, {
        owner: this.mountOwner,
        canonical: this.canonical
      })
      this.showReply(re)
    },
  },
}
</script>


