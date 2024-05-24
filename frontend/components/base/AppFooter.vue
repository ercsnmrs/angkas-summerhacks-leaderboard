<template>
  <footer class="text-gray-600 body-font">
    <div class="sm:container sm:mx-auto px-5 py-10 mx-auto flex items-center sm:flex-row flex-col">

      <span class="inline-flex sm:ml-auto sm:mt-0 mt-4 justify-center sm:justify-start">

      </span>
    </div>
  </footer>
</template>

<script>
import { mapState } from "vuex";
import feather from "feather-icons";
export default {
  data: () => {
    return {
      // Todo
    };
  },
  computed: {
    ...mapState(["copyrightDate", "socialProfiles"]),
  },
  mounted() {
    feather.replace();
  },
  updated() {
    feather.replace();
  },
};
</script>

<style scoped></style>
