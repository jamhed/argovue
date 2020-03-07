<template>
<div class="w-100" style="height: 400px"></div>
</template>

<script>
import Vis from 'vis-network'
import util from '@/util.js'

function color(node) {
  if (node.type == "DAG") {
    return "#8EFFAB"
  }
  switch(node.phase) {
    case "Failed":
      return "#FB7E81"
    default:
      return "#97C2FC"
  }
}

function shape(node) {
  switch(node.type) {
    default:
      return "box"
  }
}

export default {
  props: ['content', 'namespace', 'name'],
  data () {
    return {
      nodes: undefined,
      edges: undefined,
      options: {},
      network: undefined,
    }
  },
  mounted () {
    if (this.content.status) {
      this.update()
    }
  },
  methods: {
    update () {
      let wfNodes = util.removeRetryNodes(util.deepCopy(this.content.status.nodes))
      this.nodes = new Vis.DataSet([])
      this.edges = new Vis.DataSet([])
      Object.values(wfNodes).forEach( (node) => {
        this.nodes.add([{ id: node.id, label: node.displayName, shape: shape(node), color: color(node), type: node.type }])
      })
      Object.values(wfNodes).forEach( (node) => {
        (node.children || []).forEach( (child) => {
          this.edges.add([{ from: node.id, to: child, arrows: "to" }])
        })
      })
      this.network = new Vis.Network(this.$el, { nodes: this.nodes, edges: this.edges }, this.options)
      this.network.on("doubleClick", (ev) => {
        let node = this.nodes.get(ev.nodes[0])
        if (node && node.type == "Pod") {
          this.$router.push(`/k8s/pod/${this.namespace}/${ev.nodes[0]}`)
        }
      })
    }
  },
  watch: {
    content () {
      this.update()
    }
  }
}
</script>