<template>
  <div
    class="bg-secondary-light dark:bg-black min-h-screen flex flex-col"
  >
    <noscript>
      <strong>We're sorry but vuejs-typeeffect doesn't work properly without JavaScript enabled. Please enable it to continue.</strong>
    </noscript>
    <AppHeader />
    <transition name="fade" mode="out-in">
      <Nuxt />
    </transition>
    <div
      class="
        flex
        space-x-2
        mr-8
        mb-6
        right-0
        bottom-0
        z-50
        fixed
        items-center
        sm:space-x-4
      "
    >
      <BackToTop />
    </div>
  </div>
</template>

<script>
import feather from "feather-icons";
import AppHeader from "../components/base/AppHeader.vue";
import BackToTop from "../components/common/BackToTop.vue";
export default {
  data: () => {
    return {
      // Todo
    };
  },
  mounted() {
    feather.replace();
  },
  components: { AppHeader, BackToTop },
};
</script>

<style>
.vue-back-to-top {
  @apply p-2 sm:p-4 bg-indigo-500 hover:bg-indigo-600 text-white;
  border-radius: 50%;
  font-size: 22px;
  line-height: 22px;
}

.fade-enter-active {
  animation: coming 0.4s;
  animation-delay: 0.2s;
  opacity: 0;
}

.fade-leave-active {
  animation: going 0.4s;
}

@keyframes going {
  from {
    transform: translateX(0);
  }

  to {
    transform: translateX(-10px);
    opacity: 0;
  }
}

@keyframes coming {
  from {
    transform: translateX(-10px);
    opacity: 0;
  }

  to {
    transform: translateX(0px);
    opacity: 1;
  }
}
</style>
