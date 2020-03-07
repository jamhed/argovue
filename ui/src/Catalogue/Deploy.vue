<template>
<b-form @submit="onSubmit">
  <b-form-group label="Owner" label-for="owner">
    <b-form-select id="owner" v-model="owner" :options="owners()"></b-form-select>
  </b-form-group>
  <b-form-checkbox id="canonical" v-model="canonical" name="canonical" value="true" unchecked-value="false">
    Canonical name
  </b-form-checkbox>
  <b-form-group class="mt-4" v-for="input in input()" :key="input.name" :label="input.caption" :label-for="input.name">
    <b-form-input :id="input.name" v-model="data[input.name]" type="text" placeholder=""></b-form-input>
  </b-form-group>
  <b-button class="mt-4" type="submit" size="sm" variant="primary">Deploy</b-button>
</b-form>
</template>

<script>
function serialize(data) {
  return Object.keys(data).map( (k) => ({ name: k, value: data[k] }))
}

export default {
  props: ['object', 'name', 'namespace'],
  data () {
    return {
      owner: this.$api.effective_id(),
      canonical: "false",
      data: {}
    }
  },
  methods: {
    owners() {
      var owners = this.$api.effective_groups()
      owners.push(this.$api.effective_id())
      return owners
    },
    input () {
      let iv = this.object && this.object.spec && this.object.spec.input? this.object.spec.input : []
      iv.forEach(i => this.data[i.name] = i.default? i.default : '')
      return iv
    },
    onSubmit: async function(ev) {
      ev.preventDefault()
      let re = await this.$api.post2(`/catalogue/${this.namespace}/${this.name}/deploy`, {
        owner: this.owner,
        input: serialize(this.data),
        canonical: this.canonical
      })
      if (re.data.status == 'ok') {
        this.data = {}
      }
      this.$bvToast.toast(`${re.data.action} ${re.data.status} ${re.data.message}`, {
        title: re.data.action,
        toaster: 'b-toaster-bottom-right',
        autoHideDelay: 3000,
        noCloseButton: true,
        variant: re.data.status == 'ok'? 'info' : 'error'
      })
    },
  },
}
</script>
