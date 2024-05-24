<template>
  <div>
    <Leaderboard class="container mx-auto" />
  </div>
</template>

<script>
import feather from "feather-icons";
import Leaderboard from "../components/manila/MNLLeaderboard.vue";

export default {
  components: { Leaderboard },
  name: "Manila",
  scrollToTop: true,
  mounted() {
    feather.replace();
  },
  updated() {
    feather.replace();
  },
};
</script>

<style scoped></style>
