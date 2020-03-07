<template>
<div></div>
</template>

<script>
import JSONEditor from 'jsoneditor'
import "jsoneditor/dist/jsoneditor.min.css"

export default {
	props: ['content'],
	data () {
		return {
			editor: null
		}
	},
	methods: {
	},
	mounted () {
		this.editor = new JSONEditor(this.$el, { mode: "view", search: false, navigationBar: false, statusBar: false })
		this.editor.set(this.content)
		this.editor.expandAll()
	},
	watch: {
		content (v) {
			this.editor.set(v)
			this.editor.expandAll()
		}
	}
}
</script>

<style>
.jsoneditor-menu {
	display: none !important;
}
.jsoneditor {
	border: 0px !important;
}
.jsoneditor-tree-inner {
	padding-bottom: 0px !important;
}
</style>