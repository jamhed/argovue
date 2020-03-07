<script>
import axios from "axios"

export default {
  data() {
    return {
      auth: "undefined",
      username: undefined,
      profile: {},
      baseURL: '',
      $axios: undefined,
    };
  },
  created () {
    if (window.argovue && window.argovue.api_base_url) {
      this.baseURL = window.argovue.api_base_url
    } else {
      this.baseURL = ''
    }
    this.$axios = axios.create({ baseURL: this.baseURL, withCredentials: true})
  },
  methods: {
    uri(uri) {
      return this.baseURL + uri
    },
    redirect(url) {
      window.location.href = this.baseURL+url
    },
    get(url) {
      return this.$axios.get(url)
    },
    post(url) {
      return this.$axios.post(url)
    },
    post2(url, data) {
      return this.$axios.request({
        url: url,
        method:'post',
        withCredentials: true,
        data: data,
      })
    },
    sse(url, onMessage) {
      let es = new EventSource(this.baseURL+url, { withCredentials: true })
      es.onerror = (err) => this.$log("SSE", err)
      es.onmessage = onMessage
      return es
    },
    isAuth() {
      return this.auth == "true"
    },
    isNot() {
      return this.auth == "false"
    },
    verifyAuth: async function() {
      let ev = await this.$axios.get("/profile")
      if (ev.data && ev.data.name) {
        this.auth = "true"
        this.username = ev.data.name
        this.profile = ev.data
      } else {
        this.auth = "false"
        this.username = undefined
      }
    },
    login() {
      this.redirect('/auth')
    },
    logout() {
      this.redirect('/logout')
    },
    copy(ar) {
      return JSON.parse(JSON.stringify(ar || []))
    },
    groups () {
      return this.copy(this.profile.groups).sort()
    },
    effective_groups () {
      return this.copy(this.profile.effective_groups).sort()
    },
    effective_id () {
      return this.profile.id
    },
  }
};
</script>