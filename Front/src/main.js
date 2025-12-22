import { createApp } from "vue";
import App from "./App.vue";
import router from "./router/index.js";
import "./styles/base.css";
import "./styles/themes.css";
import "./styles/animations.css";

createApp(App).use(router).mount("#app");
