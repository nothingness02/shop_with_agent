<template>
  <component :is="layout">
    <RouterView />
  </component>
</template>

<script setup>
import { computed } from "vue";
import { RouterView, useRoute } from "vue-router";
import UserLayout from "./layouts/UserLayout.vue";
import MerchantLayout from "./layouts/MerchantLayout.vue";
import NeutralLayout from "./layouts/NeutralLayout.vue";

const route = useRoute();

const layout = computed(() => {
  if (route.meta.layout === "merchant") {
    return MerchantLayout;
  }
  if (route.meta.layout === "user") {
    return UserLayout;
  }
  return NeutralLayout;
});
</script>
