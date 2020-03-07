<template>
<div class="node" :style="this.style()" @click="click()">
  <div
    :class="`fas node-status node-status-${node.phase.toLocaleLowerCase()}`"
    :style="`lineHeight: ${node.height}px`">
  </div>
  <div class="node-title" :style="`lineHeight: ${node.height}px`">{{node.displayName}}</div>
</div>
</template>

<script>
export default {
  props: ['node', 'name', 'namespace'],
  data () {
    return {
    }
  },
  created () {
  },
  methods: {
    click() {
      if (this.node.type == "Pod") {
        this.$router.push(`/k8s/pod/${this.namespace}/${this.node.id}`)
      }
    },
    left() {
      return this.node.x - this.node.width / 2
    },
    top() {
      return this.node.y - this.node.height / 2
    },
    width() {
      return this.node.width
    },
    height() {
      return this.node.height
    },
    style() {
      return `left: ${this.left()}px; top: ${this.top()}px; width: ${this.width()}px; height: ${this.height()}px`
    }
  }
}
</script>

<style scoped>
.node {
  position: absolute;
  padding-left: 3.5em;
  box-shadow: 1px 1px 1px #CCD6DD;
  background-color: white;
  border-radius: 4px;
  border: 1px solid #d7d7d8;
  cursor: pointer;
}

.node.active {
  border-color: #00A2B3;
}

.node.virtual {
  background-color: transparent;
  box-shadow: none;
  border: none;
  padding-left: 0;
}

.node.virtual:after {
  content: '';
  position: absolute;
  display: block;
  border-radius: 10px;
  width: 20px;
  height: 20px;
  left: -10px;
  top: -10px;
  border: 1px dashed #8FA4B1;
}

.node.virtual .node-status, .node.virtual .node-title {
  display: none;
}

.node.virtual.active:after {
  border-color: #00A2B3;
}

.node.active .node-status {
  border: 1px solid #00A2B3;
  border-right: none;
}

.node-status {
  position: absolute;
  left: -1px;
  bottom: -1px;
  top: -1px;
  width: 3em;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
  text-align: center;
  color: white;
}

.node-status-error, .node-status-failed {
  background-color: #E96D76;
}

.node-status-error::after, .node-status-failed::after {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: inline-block;
  font-style: normal;
  font-variant: normal;
  font-weight: normal;
  line-height: 1;
  content: "\f057";
}

.node-status-pending {
  background-color: #f4c030;
}

.node-status-pending::after {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: inline-block;
  font-style: normal;
  font-variant: normal;
  font-weight: normal;
  line-height: 1;
  content: "\f017";
  font-size: 1em;
  animation-name: spin;
  animation-duration: 10000ms;
  animation-iteration-count: infinite;
  animation-timing-function: linear;
}

.node-status-running {
  background-color: #0DADEA;
}

.node-status-running::after {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: inline-block;
  font-style: normal;
  font-variant: normal;
  font-weight: normal;
  line-height: 1;
  content: "\f1ce";
  animation-name: spin;
  animation-duration: 4000ms;
  animation-iteration-count: infinite;
  animation-timing-function: linear;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
    }
  }

.node-status-succeeded {
  background-color: #18BE94;
}

.node-status-succeeded::after {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: inline-block;
  font-style: normal;
  font-variant: normal;
  font-weight: normal;
  line-height: 1;
  content: "\f00c";
}

.node-status-skipped {
  background-color: #CCD6DD;
}

.node-title {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.line, .node {
  transition: left 0.2s, top 0.2s, width 0.2s, height 0.2s;
}

</style>