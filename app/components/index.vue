<template>
  <div id="main-column">
    <!-- Right side is first because it'll look proper in the end -->
    <div id="side-right">
      <div id="user-box">
        <div id="user-box-login" v-if="!user">
          <h3 id="user-box-title">Log in</h3>
          <form @submit.prevent="loginSendForm()">
            <div class="user-box-form-field">Username <input type="text" v-bind:class="{invalid: loginErrorFields.includes('username')}" v-model="loginUsername"></div>
            <div class="user-box-form-field">Password <input type="password" v-bind:class="{invalid: loginErrorFields.includes('password')}" v-model="loginPassword"></div>
            <span class="user-box-form-error red" v-show="loginError">{{ loginError }}</span>
            <button class="form-submit" v-show="!loginFormSending">Submit</button>
            <button class="form-submit disabled" disabled="disabled" v-show="loginFormSending">Sending...</button>
          </form>
          <div class="user-box-form-footer">If you don't have an account, you can <router-link to="/signup">sign up here</router-link>.</div>
        </div>
        <div id="user-box-info" v-else>
          <img id="user-box-avatar" v-show="user.avatar" v-bind:src="user.avatar">
          <div class="user-box-welcome">Hi, {{ user.nickname }}!</div>
          <a href="/logout" @click.prevent="logOut()">Log out</a>
        </div>
      </div>
    </div>

    <span id="welcome-text">Welcome to <strong>Una</strong>!</span>
    <div id="about">
      <span id="about-title"><strong>Una</strong>: an Uno card game simulation</span>
      <p id="about-text"><strong>Una</strong> is an online multiplayer simulation of the Uno card game.
      Just create or join a lobby and start playing. You can even set a password to play privately.
      Create an account to get an avatar, name, points, and more.</p>
    </div>
    lobbies
    <p v-if="lobbiesLoading">lobbies are loading</p>
    <ul style="text-align:center;" v-else>
      <li v-for="lobby in lobbies">{{ lobby }}</li>
    </ul>
    <h3>create a lobby</h3>
    <form @submit.prevent="lobbySendForm()">
      name <input v-model="lobbyName">
      <br><div class="red" v-show="lobbyError">{{ lobbyError }}</div>
      <button>go</button>
    </form>
  </div>
</template>
<script>
export default {
  computed: {
    user() {
      return window.user;
    }
  },
  data() {
    return {
      loginUsername: '',
      loginPassword: '',
      loginError: '',
			loginErrorFields: [],
			loginFormSending: false,
      lobbiesLoading: true,
      lobbies: null,
      lobbyName: '',
      lobbyError: ''
    }
  },
  created() {
    var _this = this;
    this.$http.get('/app/lobbies', {
      loginUsername: this.loginUsername,
      loginPassword: this.loginPassword,
    }).then(response => {
      _this.lobbiesLoading = false;
      _this.lobbies = response.body;
    });
  },
  methods: {
    loginSendForm() {
      this.loginErrorFields = [];
      if(!this.loginUsername) {
        this.loginErrorFields.push('username');
        this.loginError = 'You have to enter a username.';
        return;
      }
      if(!this.loginPassword) {
        this.loginErrorFields.push('password');
        this.loginError = 'You have to enter your password.';
        return;
      }
			this.$set(this, 'loginFormSending', true);
			this.loginError = '';
      this.$http.post('/app/login', {
        username: this.loginUsername,
				password: this.loginPassword,
			}).then(response => {
				location.href = '/';
			}, response => {
				this.$set(this, 'loginFormSending', false);
				if(response.body.error) {
					this.loginError = response.body.error;
				} else {
					if(response.status) {
						this.loginError = 'Error: ' + response.status + ' ' + response.statusText;
					} else {
						this.loginError = 'The server seems to be down right now, try again in a moment.';
					}
				}
			});
    },
    logOut() {
      this.$http.post('/app/logout')
      .then(response => {
				location.href = '/';
			});
    },
    lobbySendForm() {
      this.lobbyError = '';
      var _this = this;
      this.$http.post('/app/lobby_create', {
        name: this.lobbyName,
			}).then(response => {
        _this.lobbyName = '';
        _this.lobbies.unshift(response.body);
			}, response => {
				if(response.body.error) {
					this.lobbyError = response.body.error;
				} else {
					if(response.status) {
						this.lobbyError = 'Error: ' + response.status + ' ' + response.statusText;
					} else {
						this.lobbyError = 'The server seems to be down right now, try again in a moment.';
					}
				}
			});
    }
  }
}
</script>
