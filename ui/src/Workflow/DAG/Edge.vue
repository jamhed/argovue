<template>
<div :key="this.key()" class="edge">
  <edgeline v-for="line in lines" v-bind:key="line.id" :line="line"></edgeline>
</div>
</template>

<script>
import Line from '@/Workflow/DAG/Line'

export default {
  props: ["edge"],
  components: {
    edgeline: Line,
  },
  data () {
    return {
      lines: []
    }
  },
  created () {
    this.lines = this.edge.lines.map((line, i) => { return {
      i: i,
      distance: Math.sqrt(Math.pow(line.x1 - line.x2, 2) + Math.pow(line.y1 - line.y2, 2)),
      xMid: (line.x1 + line.x2) / 2,
      yMid: (line.y1 + line.y2) / 2,
      angle: (Math.atan2(line.y1 - line.y2, line.x1 - line.x2) * 180) / Math.PI
    }})
  },
  methods: {
    key() {
      return `${this.edge.from}-${this.edge.to}`
    },
  },
}
</script>

<style scoped>
.edge .line:last-child:not(.line-no-arrow):after {
  content: '\25BA';
  position: absolute;
  color: #A3A3A3;
  font-size: 12px;
  top: -9px;
  left: -1px;
  transform: rotate(180deg);
}

.line, .node {
  transition: left 0.2s, top 0.2s, width 0.2s, height 0.2s;
}
</style>