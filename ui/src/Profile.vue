<template>
  <b-container>
    <b-row class="pt-3 pb-2 mb-3 border-bottom">
      <b-col cols=10>
        <h1 class="h2">profile</h1>
      </b-col>
      <b-col class="text-right">
        <b-button size="sm" variant="primary" @click="$api.logout()">Logout</b-button>
      </b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Subject</b-col>
      <b-col>{{$api.profile.sub}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Name</b-col>
      <b-col>{{$api.profile.name}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Preferred username</b-col>
      <b-col>{{$api.profile.preferred_username}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Email</b-col>
      <b-col>{{$api.profile.email}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Time Zone</b-col>
      <b-col>{{$api.profile.zoneinfo}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Effective Id</b-col>
      <b-col>{{$api.profile.effective_id}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Effective Groups</b-col>
      <b-col>
        <div v-for="group in effective_groups()" :key="group">{{group}}</div>
      </b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Groups</b-col>
      <b-col>
        <div v-for="group in groups()" :key="group">{{group}}</div>
      </b-col>
    </b-row>
    <b-row class="mt-2 pt-2 border-top"><b-col>&nbsp;</b-col></b-row>
  </b-container>
</template>

<script>
export default {
  data() {
    return {
    };
  },
  methods: {
    copy(ar) {
      return JSON.parse(JSON.stringify(ar || []))
    },
    groups () {
      return this.copy(this.$api.profile.groups).sort()
    },
    effective_groups () {
      return this.copy(this.$api.profile.effective_groups).sort()
    },
  }
};
</script>
