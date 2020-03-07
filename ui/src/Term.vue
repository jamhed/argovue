<template>
<b-container style="font-family: 'Lucida Console', Monaco, monospace; font-size: 12px">
  <b-row><b-col><b>{{visible_title}}</b></b-col></b-row>
  <b-row v-for="line in lines" :key="line"><b-col>{{line}}</b-col></b-row>
</b-container>
</template>

<script>
export default {
  props: ['lines', 'title'],
  data () {
    return {
      visible_title: ''
    }
  },
  methods: {
  },
  created () {
    if (this.lines.length > 0) {
      this.visible_title = this.title
    }
  },
  watch: {
    content () {
      if (this.lines.length > 0) {
        this.visible_title = this.title
      }
    }
  }
}
</script>