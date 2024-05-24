<template>
<section>
  <TypedEffect/>
  <div class="container mb-10">
    <div class="mx-auto max-w-xl text-center my-5 lg:mx-0 lg:flex-auto lg:text-left dark:text-white">
      <h1 class="text-3xl font-bold tracking-tight sm:text-4xl text-justify">We are team Pakopya ni Edgar </h1>
      <h3 class="text-2xl font-semibold tracking-tight sm:text-2xl text-left">Mix of product and technology</h3>

      <h3 class="text-1xl my-5 text-blue-700 dark:text-violet-200 font-semibold tracking-tight sm:text-2xl text-left">The Idea</h3>

      <p class="my-5 leading-relaxed text-justify">
        We have a real-time leaderboard to show who the top bikers are for the week based on the bikers’ RFM (Recency, Frequency, Monetary) score.
      </p>

      <p class="my-5 leading-relaxed text-justify">
        Top bikers of the week will enjoy certain rewards on the following week such as lower commission rates and free vouchers for their friends and family
      </p>

      <h3 class="text-1xl my-5 text-blue-700 dark:text-violet-200 font-semibold tracking-tight sm:text-2xl text-left">Want to know your ranking? Use the quick search</h3>
      <div class="md:w-1/1 sm:w-1/1 mt-4 w-full p-4">
        <div class="flex  w-full sm:flex-row flex-col mx-auto px-8 sm:space-x-4 sm:space-y-0 space-y-4 sm:px-0 items-end">
          <div class="relative flex-grow w-full">
            <input
              v-model="userID"
              type="text" id="addr" name="addr" class="w-full bg-gray-100 dark:bg-violet-200 bg-opacity-50 rounded border border-gray-300 focus:border-indigo-500 focus:bg-transparent focus:ring-2 focus:ring-indigo-200 text-base outline-none text-gray-700 py-1 px-3 leading-8 transition-colors duration-200 ease-in-out"
              required
            >
          </div>
          <button
            @click="postTransaction()"
            class="text-white bg-blue-500 hover:bg-blue-700 dark:bg-violet-600 dark:hover:bg-violet-700 border-0 py-2 px-8 focus:outline-none  rounded text-lg"
          >
            Send
        </button>
        </div>
        <p v-if="is_processed" class="my-5 leading-relaxed text-justify">
          {{ placeText }}
        </p>

      </div>

    </div>
  </div>
</section>

</template>

<script>
import feather from "feather-icons";
import TypedEffect from '../common/TypedEffect.vue'

export default {
  components: { TypedEffect },
  data () {
    return {
      is_processed : false,
      placeText : "",
      userID : "",
    }
  },

  mounted() {
    feather.replace();
  },
  updated() {
    feather.replace();
  },
  methods: {
    postTransaction() {
      this.is_processed = true;

      let confirm_text = 'Are you sure to send?';
      let confirm_result = confirm(confirm_text);
      if (confirm_result !== true) {
          alert('Canceled');
          return
      }
      this.$axios.get(`http://127.0.0.1:8000/leaderboard/ranking/${this.userID}`)
        .then((response) => {
          const serviceZone = response.data.service_zone;
          const average = response.data.rating.average;
          const placement = 1; // Replace with the actual placement if available in the response
          this.placeText = `You are ${serviceZone}'s ${placement}${this.getOrdinalSuffix(placement)} with an Average of ${average}`;

        })
        .catch((e) => {
          console.error(e);
          this.placeText = "You do not have record on our system. Please try again.";

        });
    },
    getOrdinalSuffix(n) {
      const s = ["th", "st", "nd", "rd"];
      const v = n % 100;
      return s[(v - 20) % 10] || s[v] || s[0];
    }
  }
};
</script>

<style lang="scss" scoped></style>
