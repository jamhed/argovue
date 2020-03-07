<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{namespace}}/{{kind}}/{{name}}</h1>
    </div>
    <div>
      <jsoneditor :content="object"></jsoneditor>
    </div>
  </div>
</template>

<script>
import SSE from '@/SSE/Object.vue'
import JsonEditor from '@/JsonEditor'

export default {
  props: ["namespace", "kind", "name"],
  extends: SSE,
  components: {
    jsoneditor: JsonEditor
  },
  data() {
    return {
    };
  },
};
</script>
