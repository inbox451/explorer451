export default defineNuxtPlugin(({ $pinia, $router, $route }) => {
  $pinia.use(({ store }) => {
    store.$router = $router
    store.$route = $route
  })
})
