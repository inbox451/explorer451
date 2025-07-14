<script setup>
import Header from '~/components/navigation/Header.vue';
const year = new Date().getFullYear();
</script>

<template>
  <div>
    <Header />
    <main>
      <NuxtPage />
    </main>
    <footer>
      <p>&copy; {{ year }} My Application</p>
    </footer>
  </div>
</template>