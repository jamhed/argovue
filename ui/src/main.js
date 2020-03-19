import Vue from 'vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import '@fortawesome/fontawesome-free/css/solid.css'

import App from '@/App'
import Version from '@/Version'
import Front from '@/Front'
import Objects from '@/Objects'
import Datasources from '@/Datasources'
import Pvcs from '@/Pvcs'
import Services from '@/Services'
import WorkflowObject from '@/Workflow/Object'
import CatalogueObject from '@/Catalogue/Object'
import CatalogueInstance from '@/Catalogue/Instance'
import PodObject from '@/Pod/Object'
import JobObject from '@/Job/Object'
import ServiceObject from '@/Service/Object'
import PvcObject from '@/Pvc/Object'
import IngressObject from '@/Ingress/Object'
import DatasourceObject from '@/Datasource/Object'
import Profile from '@/Profile'
import API from '@/API'

Vue.use(VueRouter)
Vue.use(BootstrapVue)

if (window.argovue && process.env.VUE_APP_API_BASE_URL) {
  window.console.log("Setting up environment")
  window.argovue.api_base_url = process.env.VUE_APP_API_BASE_URL
}

Vue.prototype.$api = new Vue(API)
Vue.prototype.$log = window.console.log.bind(console)

function routeProps({params}) {
  return params
}

const router = new VueRouter({
  routes: [
    { path: '/', component: Front },
    { path: '/version', component: Version },
    { path: '/profile', component: Profile },
    { path: '/watch/datasources', component: Datasources, props: routeProps },
    { path: '/watch/pvcs', component: Pvcs, props: routeProps },
    { path: '/watch/services', component: Services, props: routeProps },
    { path: '/watch/:kind', component: Objects, props: routeProps },
    { path: '/workflows/:namespace/:name', component: WorkflowObject, props: routeProps },
    { path: '/k8s/job/:namespace/:name', component: JobObject, props: routeProps },
    { path: '/k8s/pod/:namespace/:name', component: PodObject, props: routeProps },
    { path: '/k8s/persistentvolumeclaim/:namespace/:name', component: PvcObject, props: routeProps },
    { path: '/k8s/ingress/:namespace/:name', component: IngressObject, props: routeProps },
    { path: '/k8s/datasource/:namespace/:name', component: DatasourceObject, props: routeProps },
    { path: '/k8s/service/:namespace/:name', component: ServiceObject, props: routeProps },
    { path: '/k8s/workflow/:namespace/:name', component: WorkflowObject, props: routeProps },
    { path: '/catalogue/:namespace/:name', component: CatalogueObject, props: routeProps },
    { path: '/catalogue/:namespace/:name/instance/:instance', component: CatalogueInstance, props: routeProps },
  ]
})

router.beforeEach((to, from, next) => {
  if (router.app.$api.isAuth()) {
    next()
  } else {
    router.app.$api.verifyAuth()
    next()
  }
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')