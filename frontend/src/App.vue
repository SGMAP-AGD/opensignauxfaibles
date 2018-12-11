<template>
  <v-app id="app">
    <v-content>
      <NavigationDrawer v-if="login && drawer"/>
      <Login v-if="!login" :state="login"/>
      <router-view v-if="login"/>
    </v-content>
  </v-app>
</template>

<script>
import Login from "@/components/Login";
import NavigationDrawer from "@/components/NavigationDrawer";

export default {
  methods: {
    handleResize() {
      this.height = Math.max(
        document.documentElement.clientHeight,
        window.innerHeight || 0
      )
    },
    handleScrolling(event) {
      
      this.scrollTop = event.pageY || window.scrollY
    }
  },
  components: { Login, NavigationDrawer },
  computed: {
    height: {
      get() {
        return this.$store.state.height;
      },
      set(height) {
        this.$store.dispatch("setHeight", height);
      }
    },
    scrollTop: {
      get() {
        return this.$store.state.scrolltop;
      },
      set(scrollTop) {
        this.$store.dispatch("setScrollTop", scrollTop);
      }
    },
    uploads() {
      return this.$store.getters.getUploads;
    },
    login() {
      return this.$store.state.token != null;
    },
    batches() {
      return JSON.stringify(this.$store.state.batches, null, 2);
    },
    dbstatus() {
      return this.$store.state.dbstatus;
    },
    messages() {
      return this.$store.getters.messages;
    },
    drawer: {
      get() {
        return this.$store.state.appDrawer;
      },
      set(val) {
        this.$store.dispatch("setDrawer", val);
      }
    }
  },
  data() {
    return {
      fixed: false,
      dialog: false,
      activityFilter: ["info", "warning", "critical"],
      menu: {
        title: "Accueil",
        color: "light-green darken-4"
      },
      items: [
        {
          title: "Accueil",
          action: "/",
          icon: "fa-home",
          color: "green darken-3"
        },
        {
          title: "Détection",
          action: "/Browse",
          icon: "fa-search",
          color: "indigo darken-4"
        },
        {
          title: "Données",
          action: "/data",
          icon: "fa-database",
          color: "red darken-4"
        },
        {
          title: "Administration",
          action: "/admin",
          icon: "fa-users",
          color: "gray darken-4"
        }
      ],
      avatar: {
        debug: { icon: "fa-cogs", color: "blue" },
        info: { icon: "fa-info", color: "green" },
        warning: { icon: "fa-exclamation-triangle", color: "yellow" },
        critical: { icon: "fa-sad-cry", color: "red" }
      }
    };
  },
  mounted() {
    document.onscroll = this.handleScrolling
    window.addEventListener("resize", this.handleResize)
    this.height = Math.max(
      document.documentElement.clientHeight,
      window.innerHeight || 0
    );
  },
  name: "App"
}
</script>

<style>
@import url("https://fonts.googleapis.com/css?family=Quicksand");
@import url("https://fonts.googleapis.com/css?family=Signika");
body {
  font-family: "Quicksand", sans-serif;
}
.task {
  background: linear-gradient(0deg, #1f1a24, #1d1c24 10%, #1d1c24 90%, #353c47);
}
.tasktext {
  color: white;
}
.tasktitle {
  color: white;
  font-size: 20px;
  font-weight: 500;
}
.toolbar {
  color: "black";
  font-family: "Quicksand", sans-serif;
  text-shadow: 0px 0px 2px rgb(0, 0, 0), 0px 0px 1px rgb(255, 255, 255);
  background: linear-gradient(0deg, #272629, #ffffff00 7%, transparent),
    radial-gradient(
      ellipse at top,
      #a49bbd,
      #716a81 35%,
      #3a3b4b 75%,
      #21213b 100%,
      transparent
    );
}
.toolbar-widget {
  color: "black";
  font-family: "Quicksand", sans-serif;
  text-shadow: 0px 0px 2px rgb(146, 146, 146), 0px 0px 1px rgb(255, 255, 255);
  background: linear-gradient(0deg, #b6b6b6, #6e6e6e00 7%, transparent),
    radial-gradient(
      ellipse at top,
      #ffffff,
      #e4e4e4 35%,
      #c9c9cc 80%,
      #afafcc 100%,
      transparent
    );
}
#app {
  background: radial-gradient(
    circle at center,
    rgb(255, 255, 255),
    rgb(233, 232, 232) 75%,
    rgb(187, 187, 187) 100%
  );
}
.rightDrawer {
  position: fixed;
  right: 0px;
}
.span {
  max-height: 10px;
}
</style>
