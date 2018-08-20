<template>
  <div id="main-column">
    <h1 id="page-title">Sign up</h1>
    <form id="form-box" @submit.prevent="submitForm()">
      <span class="form-field">Username <span class="form-note">(This is what you will use to sign in.)</span> <input class="form-input" type="text" v-bind:class="{invalid: errorFields.includes('username')}" v-model="username"></span>
      <span class="form-field">Nickname <span class="form-note">(This is what you will be seen as.)</span> <input class="form-input" type="text" v-bind:class="{invalid: errorFields.includes('nickname')}" v-model="nickname"></span>
      <span class="form-field">Password <input class="form-input" type="password" v-bind:class="{invalid: errorFields.includes('password')}" v-model="password"></span>
      <span class="form-field">Confirm Password <input class="form-input" type="password" v-bind:class="{invalid: errorFields.includes('passwordAgain')}" v-model="passwordAgain"></span>
      <p class="red" v-show="error">{{ error }}</p>
      <button class="form-submit" v-show="!formSending">Submit</button>
      <button class="form-submit disabled" disabled="disabled" v-show="formSending">Sending...</button>
    </form>
    <div class="page-footer-link">Or, you can <router-link to="/login">log in</router-link>.</div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      username: '',
      nickname: '',
      password: '',
      passwordAgain: '',
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
      if(!this.nickname) {
        this.errorFields.push('nickname');
        this.error = 'You have to enter a nickname.';
        return;
      }
      if(!this.password) {
        this.errorFields.push('password');
        this.error = 'You have to enter a password.';
        return;
      }
      if(this.password != this.passwordAgain) {
        this.errorFields.push('password');
        this.errorFields.push('passwordAgain');
        this.error = 'Your passwords don\'t match.';
        return;
      }
      if(!/^[A-Za-z0-9][^/]{1,32}$/.test(this.username)) {
        this.errorFields.push('username');
        this.error = 'Your username contains a forward slash (/), has no letters, or is too long or too short.';
        return;
      }
      if(this.nickname.length > 64) {
        this.errorFields.push('nickname');
        this.error = 'Your nickname is too long.';
        return;
      }
			this.$set(this, 'formSending', true);
			this.error = '';
      this.$http.post('/app/signup', {
        username: this.username,
				nickname: this.nickname,
				password: this.password,
				password_again: this.passwordAgain
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
