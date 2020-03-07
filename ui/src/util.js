export default {
  removeRetryNodes (nodes) {
  var toReplace = {}
  var toRemove = {}
  Object.values(nodes).forEach( (node) => {
    if (node.type == 'Retry' && node.children && node.children.length > 0) {
      var last = node.children.length-1
      toRemove[node.id] = true
      node.children.forEach( (nodeId, i) => i != last? toRemove[nodeId] = true : '')
      toReplace[node.id] = node.children[last]
    }
  })
  var re = {}
  Object.values(nodes).forEach( (node) => {
    if (!toRemove[node.id]) {
      if (node.children) {
        var children = node.children.
          map( (nodeId) => toReplace[nodeId]? toReplace[nodeId] : nodeId )
        node.children = children
      }
      re[node.id] = node
    }
  })
  return re
  },
  deepCopy (thing) {
    return JSON.parse(JSON.stringify(thing))
  },
  phase (phase) {
    switch (phase) {
      case "Running":
        return "R"
      case "Pending":
        return "P"
      case "Succeeded":
        return "S"
      case "Failed":
        return "F"
      case "Unknown":
        return "U"
      case "Bound":
        return "B"
      case undefined:
        return "_"
      default:
        return phase
    }
  },
  owner(obj) {
    if (obj && obj.metadata) {
      if (obj.metadata.annotations && obj.metadata.annotations['oidc.argovue.io/owner']) {
        return obj.metadata.annotations['oidc.argovue.io/owner']
      } else if (obj.metadata.labels && obj.metadata.labels['oidc.argovue.io/group']) {
        return obj.metadata.labels['oidc.argovue.io/group']
      } else {
        return "_"
      }
    } else {
      return "*"
    }
  },
}