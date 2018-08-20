<template>
  <body>
    <ul id="navbar">
      <li><router-link to="/" v-bind:class="{selected: $route.name == 'index'}">Home</router-link></li>
      <li><router-link to="/my_profile" v-bind:class="{selected: $route.name == 'profile'}">Profile</router-link></li>
      <li><router-link to="/help" v-bind:class="{selected: $route.name == 'help'}">How to Play</router-link></li>
      <li><router-link to="/about" v-bind:class="{selected: $route.name == 'about'}">About</router-link></li>
    </ul>
    <router-view></router-view>
  </body>
</template>
<script>
export default {
  
}
</script>
