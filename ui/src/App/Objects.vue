<template>
<ul class="nav flex-column">
  <b-nav-item v-for="obj in objects" :key="obj" :to="`/watch/${obj}`">{{obj}}</b-nav-item>
  <b-nav-item :to="`/version`">version</b-nav-item>
</ul>
</template>

<script>
export default {
  data() {
    return {
      objects: []
    };
  },
  mounted: async function () {
    let re = await this.$api.get("/objects")
    this.objects = re.data.sort()
  },
}
</script>