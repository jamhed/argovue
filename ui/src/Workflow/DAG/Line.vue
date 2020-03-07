<template>
<div class="line" :key="line.i" :style="style()" :class="{ 'line-no-arrow': line.noArrow }">
</div>
</template>

<script>
export default {
  props: ["line"],
  methods: {
    style() {
      let line = this.line
      return `width: ${line.distance}px; left: ${line.xMid - line.distance / 2}px; top: ${line.yMid}px; transform: rotate(${line.angle}deg)`
    }
  },
}
</script>

<style scoped>
.line {
  position: absolute;
  border-top: 1px dashed #8FA4B1;
}

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