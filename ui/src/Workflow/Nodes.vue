<template>
<b-container fluid>
  <node v-for="node in nodes" v-bind:key="node.name" :content="node" :name="content.metadata.name" :namespace="content.metadata.namespace"></node>
</b-container>
</template>

<script>
import Node from '@/Workflow/Node'

function sortNodes(hash, nodes) {
  if (!nodes || nodes.length == 0) {
    return []
  }
  return nodes.map(
    (name) => {
      let node = hash[name];
      if (!node) {
        return [];
      }
      delete hash[name];
      return [node].concat(sortNodes(hash, node.children))
    }).flat()
}

export default {
  props: ['content', 'namespace'],
  components: {
    node: Node
  },
  data () {
    return {
      nodes: []
    }
  },
  methods: {
  },
  watch: {
    content (c) {
      let nodes = JSON.parse(JSON.stringify(c.status.nodes)) // deep copy
      let names = Object.values(nodes).map( node => node.id )
      this.$set(this, "nodes", sortNodes(nodes, [names.shift()]))
    }
  }
}
</script>


