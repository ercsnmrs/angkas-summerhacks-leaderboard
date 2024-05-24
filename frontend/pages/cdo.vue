<template>
  <div class="container mx-auto">
    <Leaderboard />
  </div>
</template>

<script>
import feather from "feather-icons";
import Button from "../components/common/Button.vue";
import Leaderboard from "../components/cdo/CDOLeaderboard.vue";

export default {
  scrollToTop: true,
  components: { Button, Leaderboard },
  mounted() {
    feather.replace();
  },
  updated() {
    feather.replace();
  },
};
</script>

<style scoped></style>
