import Vue from 'vue';
import VueRouter from 'vue-router';
import VueResource from 'vue-resource';
Vue.use(VueRouter);
Vue.use(VueResource);

// Pass window.csrf as _csrf to all POST requests
Vue.http.interceptors.push((request, next) => {
	if(request.method == 'POST') {
		if(request.body == undefined) {
			request.body = {};
		}
		request.body['gorilla.csrf.Token'] = window.csrf;
	}
	next();
});
Vue.http.options.emulateJSON = true;

import Layout from './components/layout.vue';

import Index from './components/index.vue';
import Login from './components/login.vue';
import Signup from './components/signup.vue';

const router = new VueRouter({
  mode: 'history',
  routes: [
    {path: '/', component: Index, name: 'index'},
    {path: '/login', component: Login, name: 'login'},
    {path: '/signup', component: Signup, name: 'signup'}
  ]
});

const app = new Vue({
  router,
  render: h => h(Layout)
}).$mount('body');
