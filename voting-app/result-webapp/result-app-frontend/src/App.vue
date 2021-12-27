<template>
  <h1>Welcome to the result webapp from the voting-app stack !</h1>
  <h2>Vote results:</h2>
  <ul>
    <li>Option 1: {{ votes.firstOption }}</li>
    <li>Option 2: {{ votes.secondOption }}</li>
  </ul>
</template>

<script lang="ts">
  import { Vue, Options } from 'vue-class-component';
  import axios from "axios";

  @Options({
    data() {
      return {
        votes: {
          firstOption: 0,
          secondOption: 0
        }
      }
    },
    methods: {
      getVotes() {
        axios.get("/api/votes")
          .then(response => {
            console.log(response)
            this.votes.firstOption = response.data[0];
            this.votes.secondOption = response.data[1];
          })
          .catch(function (err) {
            if (err.response) {
              console.log("Server Error:", err)
            } else if (err.request) {
              console.log("Network Error:", err)
            } else {
              console.log("Client Error:", err)
            }
          });
      }
    },
    mounted() {
      this.getVotes();
    }  
  })
  export default class App extends Vue {}
</script>

<style lang="scss" scoped>

</style>
