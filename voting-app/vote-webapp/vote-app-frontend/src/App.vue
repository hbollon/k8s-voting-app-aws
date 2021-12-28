<template>
  <h1>Welcome to the vote webapp from the voting-app stack !</h1>
  <h2>Vote now:</h2>
  <ul>
    <li>Option 1: <button v-on:click="submitVote(0)">Vote</button> </li>
    <li>Option 2: <button v-on:click="submitVote(1)">Vote</button> </li>
  </ul>
</template>

<script lang="ts">
  import { Vue, Options } from 'vue-class-component';
  import axios from "axios";

  @Options({
    data() {
      return {
        vote: 0
      }
    },
    methods: {
      submitVote(choice: number) {
        this.vote = choice;
        axios.post("/api/vote", {
          "vote": this.vote
        })
          .then(response => {
            console.log(response)
            alert(`Your vote for the option ${this.vote} is submitted!`);
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
    }
  })
  export default class App extends Vue {}
</script>

<style lang="scss" scoped>

</style>
