<template>
  <div class="container mx-auto">
    <AppBanner />
  </div>
</template>

<script>
import AppBanner from "../components/base/AppBanner.vue";
import Button from "../components/common/Button.vue";
import Industry from "../components/home/Industry.vue";
import Stack from "../components/home/Stack.vue";
export default {
  scrollToTop: true,
  components: { AppBanner, Button, Industry, Stack },
};
</script>

<style scoped></style>
