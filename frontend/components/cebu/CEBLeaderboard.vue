<template>
  <section>
    <div class="container mb-10">
      <div>
        <div class="mx-auto max-w-xl text-center my-2 lg:mx-0 lg:flex-auto lg:text-left">
          <h1 class="text-3xl font-bold dark:text-white tracking-tight sm:text-4xl text-justify">CEBU</h1>
          <p class="text-2xl font-semibold text-blue-700 dark:text-violet-200 font-medium tracking-tight sm:text-2xl text-justify">
            Checkout the leaderboard in CEBU Service Zone
          </p>
        </div>

        <p class="my-5 leading-relaxed text-justify text-blue-700 dark:text-violet-200">
          Disclaimer: Ang inyong score ay nakabase sa mga sumusunod na kategorya: <br />
          1. Completed Rides <br />
          2. Net Earnings <br />
          3. Tagal mula noong huling booking <br />
          4. Dalas ng booking <br />
        </p>

        <div v-if="drivers.length > 0">
          <h2 class="text-xl font-bold dark:text-white tracking-tight sm:text-2xl text-justify">Leaderboard</h2>
          <ul>
            <li v-for="(driver, index) in drivers" :key="driver.driver_id" class="my-2 p-2 border-b border-gray-300">
              <p class="text-lg font-semibold">Rank {{ index + 1 }}: {{ driver.service_zone }}'s Driver</p>
              <p class="text-sm">Average Rating: {{ driver.rating.average }}</p>
              <p class="text-sm">Completed Trips: {{ driver.number_of_completed_trips }}</p>
              <p class="text-sm">Net Income: {{ driver.net_income }}</p>
              <p class="text-sm">Last Completed Trip Date: {{ formatDate(driver.last_completed_trip_date) }}</p>
            </li>
          </ul>
        </div>
        <div v-else>
          <p>No data available.</p>
        </div>

        <p class="my-5 leading-relaxed text-justify text-blue-700 dark:text-violet-200">
          Rewards <br />
          Ang Top 10 ay makakuha ng mga sumusunod na premyo: <br />
          1. 1000 PHP <br />
          2. 15% Commission Fee <br />
          3. Voucher for 3 rides <br />
        </p>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      drivers: []
    };
  },
  mounted() {
    this.fetchLeaderboard();
  },
  methods: {
    fetchLeaderboard() {
      this.$axios.get('http://localhost:8000/leaderboard/CEB')
        .then((response) => {
          this.drivers = response.data.drivers;
        })
        .catch((error) => {
          console.error('Error fetching leaderboard:', error);
        });
    },
    formatDate(dateString) {
      const options = { year: 'numeric', month: 'long', day: 'numeric' };
      return new Date(dateString).toLocaleDateString(undefined, options);
    }
  }
};
</script>

<style scoped>
.container {
  padding: 20px;
}
</style>
