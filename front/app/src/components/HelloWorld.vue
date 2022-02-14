<template>
  <Draggable
    :flatData="flatData"
    idKey="id"
    parentIdKey="pid"
    childrenLazyLoading
    :childrenLoader="onClickUnfold"
  >
    <template v-slot="{ node, tree }">
      <b @click="tree.toggleFold(node)">{{ node.$folded ? "+" : "-" }}</b>
      <span>{{ node.text }}</span>
    </template>
  </Draggable>
</template>
<script>
import "@he-tree/vue3/dist/he-tree-vue3.css";
import { Draggable } from "@he-tree/vue3";
import axios from "axios";

export default {
  components: { Draggable },
  data() {
    return {
      flatData: [],
    };
  },
  mounted() {
    console.log("asdf");
    axios.get("/tree.json").then((response) => {
      console.log(response);
      this.flatData = response.data;
    });
  },
  methods: {
    onClickUnfold: async (node) => {
      if (node.left + 1 == node.right) {
        return [];
      } else {
        return [{ text: "child" }];
      }
    },
  },
};
</script>