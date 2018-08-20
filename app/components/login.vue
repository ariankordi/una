<template>
  <div id="main-column">
    <h1 id="page-title">Log in</h1>
    <form id="form-box" @submit.prevent="submitForm()">
      <div class="form-field">Username <input class="form-input" type="text" v-bind:class="{invalid: errorFields.includes('username')}" v-model="username"></div>
      <div class="form-field">Password <input class="form-input" type="password" v-bind:class="{invalid: errorFields.includes('password')}" v-model="password"></div>
      <p class="red" v-show="error">{{ error }}</p>
      <button class="form-submit" v-show="!formSending">Submit</button>
      <button class="form-submit disabled" disabled="disabled" v-show="formSending">Sending...</button>
    </form>
    <div class="page-footer-link">If you don't have an account, you can <router-link to="/signup">sign up here</router-link>.</div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      username: '',
      password: '',
      error: '',
			errorFields: [],
			formSending: false
    }
  },
  methods: {
    submitForm() {
      this.errorFields = [];
      if(!this.username) {
        this.errorFields.push('username');
        this.error = 'You have to enter a username.';
        return;
      }
      if(!this.password) {
        this.errorFields.push('password');
        this.error = 'You have to enter your password.';
        return;
      }
			this.$set(this, 'formSending', true);
			this.error = '';
      this.$http.post('/app/login', {
        username: this.username,
				password: this.password,
			}).then(response => {
				location.href = '/';
			}, response => {
				this.$set(this, 'formSending', false);
				if(response.body.error) {
					this.error = response.body.error;
				} else {
					if(response.status) {
						this.error = 'Error: ' + response.status + ' ' + response.statusText;
					} else {
						this.error = 'The server seems to be down right now, try again in a moment.';
					}
				}
			});
    }
  }
}
</script>
